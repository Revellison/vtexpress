import React from 'react';
import styles from './History.module.css';

const History = () => {
    return (
        <div className={styles.history}>
            <div className={styles.header}>
                <h1 className={styles.title}>History</h1>
                <p className={styles.subtitle}>Заглушка страницы истории действий</p>
            </div>

            <div className={styles.listPlaceholder}>
                <div className={styles.card}>
                    <div className={styles.cardTitle}>Нет записей</div>
                    <div className={styles.cardText}>Здесь будет отображаться история ваших действий.</div>
                </div>
            </div>
        </div>
    );
};

export default History;