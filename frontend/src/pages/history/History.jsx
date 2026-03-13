import React, { useEffect, useState } from 'react';
import styles from './History.module.css';
import { getHistory } from '../../lib/backend';
import AnalysCard from '../../components/analys_card/AnalysCard';

const History = ({ latestScan }) => {
    const [historyItems, setHistoryItems] = useState([]);
    const [error, setError] = useState('');

    useEffect(() => {
        let isMounted = true;
        getHistory()
            .then((items) => {
                if (isMounted) {
                    setHistoryItems(Array.isArray(items) ? items : []);
                }
            })
            .catch((fetchError) => {
                if (isMounted) {
                    setError(fetchError?.message || 'Не удалось загрузить историю');
                }
            });

        return () => {
            isMounted = false;
        };
    }, [latestScan]);

    return (
        <div className={styles.history}>
            <div className={styles.header}>
                <h1 className={styles.title}>History</h1>
                <p className={styles.subtitle}>История всех сканов из локальной базы SQLite</p>
            </div>

            <div className={styles.listPlaceholder}>
                {error ? <div className={styles.cardText}>{error}</div> : null}
                {historyItems.length === 0 ? (
                    <div className={styles.card}>
                        <div className={styles.cardTitle}>Нет записей</div>
                        <div className={styles.cardText}>Здесь будет отображаться история ваших действий.</div>
                    </div>
                ) : (
                    <div className={styles.list}>
                        {historyItems.map((item) => (
                            <AnalysCard key={item.id} item={item} />
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default History;