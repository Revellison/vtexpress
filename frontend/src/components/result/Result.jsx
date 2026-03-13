import styles from './Result.module.css';
import AIAnswer from '../ai_answer/AIAnswer';

const Result = ({ scan, onSummaryRequest, isSummarizing }) => {
    if (!scan) {
        return null;
    }

    const summary = scan.payload?.vt_summary;
    const engines = scan.payload?.engines_verdict?.engines || [];
    const domains = scan.payload?.ip_domain?.domain || [];
    const ips = scan.payload?.ip_domain?.ip || [];

    return (
        <div className={styles.result}>
            <h2>Результат анализа</h2>
            <div className={styles.meta}>
                <span>Файл: {scan.fileName}</span>
                <span>SHA256: {scan.fileSha256}</span>
                <span>Время: {scan.scannedAt}</span>
            </div>

            <div className={styles.summaryGrid}>
                <div>Вердикт: {summary?.verdict || 'unknown'}</div>
                <div>Детектов: {summary?.malicious || 0}</div>
                <div>Подозрительных: {summary?.suspicious || 0}</div>
                <div>Недетектов: {summary?.undetected || 0}</div>
                <div>Всего движков: {summary?.totalEngines || engines.length}</div>
            </div>

            <div className={styles.section}>
                <h3>Engines</h3>
                <div className={styles.list}>
                    {engines.map((engine) => (
                        <div className={styles.item} key={engine.engineName}>
                            <span>{engine.engineName}</span>
                            <span>{engine.category}</span>
                            <span>{engine.result || '-'}</span>
                        </div>
                    ))}
                </div>
            </div>

            <div className={styles.section}>
                <h3>IP/Domain</h3>
                <div className={styles.list}>
                    {ips.map((ip) => (
                        <div className={styles.item} key={ip.address}>
                            <span>{ip.address}</span>
                            <span>{ip.provider || '-'}</span>
                            <span>{ip.region || '-'}</span>
                        </div>
                    ))}
                    {domains.map((domain) => (
                        <div className={styles.item} key={domain.domain}>
                            <span>{domain.domain}</span>
                            <span>{domain.resolved || '-'}</span>
                            <span>{domain.provider || '-'}</span>
                        </div>
                    ))}
                </div>
            </div>

            <button className={styles.summaryBtn} onClick={onSummaryRequest} disabled={isSummarizing}>
                {isSummarizing ? 'Генерация...' : 'Summary with AI'}
            </button>

            <AIAnswer summary={scan.aiSummary} />
        </div>
    );
};

export default Result;
