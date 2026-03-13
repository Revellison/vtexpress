package modules

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"vtgui/internal/models"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS settings (
            key TEXT PRIMARY KEY,
            value TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'string',
			updated_at TEXT NOT NULL
        );

		CREATE TABLE IF NOT EXISTS scans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at TEXT NOT NULL,
			file_name TEXT NOT NULL,
			file_sha256 TEXT NOT NULL,
			payload_json TEXT NOT NULL,
			vt_raw_json TEXT NOT NULL,
			ai_summary TEXT NOT NULL DEFAULT ''
		);
    `)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Close()
}

func (s *Storage) SetSetting(key, value, valueType string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return errors.New("setting key is empty")
	}
	if valueType == "" {
		valueType = "string"
	}
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(`
		INSERT INTO settings(key, value, type, updated_at) VALUES(?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			type = excluded.type,
			updated_at = excluded.updated_at
	`, key, value, valueType, now)
	return err
}

func (s *Storage) GetSetting(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var value string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (s *Storage) SaveScan(result models.ScanResult) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	payloadJSON, err := json.Marshal(result.Payload)
	if err != nil {
		return 0, fmt.Errorf("marshal payload: %w", err)
	}

	vtRaw := result.RawVT
	if len(vtRaw) == 0 {
		vtRaw = json.RawMessage("{}")
	}

	res, err := s.db.Exec(`
		INSERT INTO scans(created_at, file_name, file_sha256, payload_json, vt_raw_json, ai_summary)
		VALUES(?, ?, ?, ?, ?, ?)
	`, result.ScannedAt, result.FileName, result.FileSHA256, string(payloadJSON), string(vtRaw), result.AISummary)
	if err != nil {
		return 0, fmt.Errorf("insert scan: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return id, nil
}

func (s *Storage) ListScans(limit int) ([]models.HistoryItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 {
		limit = 200
	}

	rows, err := s.db.Query(`
		SELECT id, created_at, file_name, file_sha256, payload_json, vt_raw_json, ai_summary
		FROM scans
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("query scans: %w", err)
	}
	defer rows.Close()

	items := make([]models.HistoryItem, 0, limit)
	for rows.Next() {
		var (
			item        models.HistoryItem
			payloadJSON string
			vtRawJSON   string
		)

		err := rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.FileName,
			&item.FileSHA256,
			&payloadJSON,
			&vtRawJSON,
			&item.AISummary,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		if payloadJSON != "" {
			if err := json.Unmarshal([]byte(payloadJSON), &item.Payload); err != nil {
				return nil, fmt.Errorf("unmarshal payload: %w", err)
			}
		}

		if vtRawJSON == "" {
			vtRawJSON = "{}"
		}
		item.RawVT = json.RawMessage(vtRawJSON)

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return items, nil
}
