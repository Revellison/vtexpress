package modules

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	keyringService = "vtexpress"
	vtAPIKeyName   = "vt_api_key"
	aiAPIKeyName   = "ai_api_key"
)

type KeyringStore struct{}

func NewKeyringStore() *KeyringStore {
	return &KeyringStore{}
}

func (k *KeyringStore) SetVTAPIKey(value string) error {
	if value == "" {
		return errors.New("vt api key is empty")
	}
	if err := keyring.Set(keyringService, vtAPIKeyName, value); err != nil {
		return fmt.Errorf("set vt key: %w", err)
	}
	return nil
}

func (k *KeyringStore) SetAIAPIKey(value string) error {
	if value == "" {
		return errors.New("ai api key is empty")
	}
	if err := keyring.Set(keyringService, aiAPIKeyName, value); err != nil {
		return fmt.Errorf("set ai key: %w", err)
	}
	return nil
}

func (k *KeyringStore) GetVTAPIKey() (string, error) {
	value, err := keyring.Get(keyringService, vtAPIKeyName)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get vt key: %w", err)
	}
	return value, nil
}

func (k *KeyringStore) GetAIAPIKey() (string, error) {
	value, err := keyring.Get(keyringService, aiAPIKeyName)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get ai key: %w", err)
	}
	return value, nil
}
