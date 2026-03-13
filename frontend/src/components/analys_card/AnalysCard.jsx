import styles from './AnalysCard.module.css';

const AnalysCard = ({ item }) => {
    return (
        <article className={styles.card}>
            <div className={styles.row}><strong>{item.fileName}</strong></div>
            <div className={styles.row}>Дата: {item.createdAt}</div>
            <div className={styles.row}>SHA256: {item.fileSha256}</div>
            <div className={styles.row}>Вердикт: {item.payload?.vt_summary?.verdict || 'unknown'}</div>
            <div className={styles.row}>Детекты: {item.payload?.vt_summary?.malicious || 0}</div>
            {item.aiSummary ? <div className={styles.row}>AI: {item.aiSummary.slice(0, 240)}...</div> : null}
        </article>
    );
};

export default AnalysCard;
