import { useEffect, useState } from 'react';
import styles from './Home.module.css';
import FileArrowIcon from '../../components/icons/filearrow';
import Result from '../../components/result/Result';
import { OnFileDrop, OnFileDropOff } from '../../../wailsjs/runtime/runtime';
import { analyzeFileByPath, pickFile, summarizePayload } from '../../lib/backend';

const Home = ({ latestScan, onScanComplete }) => {
    const [isLoading, setIsLoading] = useState(false);
    const [isSummarizing, setIsSummarizing] = useState(false);
    const [status, setStatus] = useState('');

    const runScan = async (filePath) => {
        if (!filePath) {
            return;
        }
        setIsLoading(true);
        setStatus('Запущен анализ файла...');
        try {
            const scanResult = await analyzeFileByPath({ path: filePath, runAiSummary: false });
            onScanComplete(scanResult);
            setStatus('Анализ завершён');
        } catch (error) {
            setStatus(error?.message || 'Ошибка во время анализа');
        } finally {
            setIsLoading(false);
        }
    };

    const handleSelectFile = async () => {
        try {
            const filePath = await pickFile();
            await runScan(filePath);
        } catch (error) {
            setStatus(error?.message || 'Не удалось выбрать файл');
        }
    };

    const handleSummary = async () => {
        if (!latestScan?.payload) {
            return;
        }
        setIsSummarizing(true);
        try {
            const summary = await summarizePayload({ payload: latestScan.payload });
            onScanComplete({ ...latestScan, aiSummary: summary });
            setStatus('AI summary успешно сформирован');
        } catch (error) {
            setStatus(error?.message || 'Ошибка при AI суммаризации');
        } finally {
            setIsSummarizing(false);
        }
    };

    useEffect(() => {
        OnFileDrop((_, __, paths) => {
            if (Array.isArray(paths) && paths.length > 0) {
                runScan(paths[0]);
            }
        }, false);

        return () => {
            OnFileDropOff();
        };
    }, []);

    return (
        <div className={styles.home}>
            <div className={styles.container}>
                <div className={styles.ddrop}>
                    <div className={styles.ddrop_label}>
                        <FileArrowIcon width={48} height={48}/>
                        <span>Перетащите файл</span>
                    </div>
                    <div className={styles.or}>
                        <span>или</span>
                    </div>
                    <div className={styles.select}>
                        <button onClick={handleSelectFile} disabled={isLoading}>{isLoading ? 'Анализ...' : 'Выбрать файл'}</button>
                        <span>Поддерживаются: .exe, .apk, .zip и другие</span>
                        <span className={styles.status}>{status}</span>
                    </div>
                </div>
            </div>
            <Result scan={latestScan} onSummaryRequest={handleSummary} isSummarizing={isSummarizing} />
        </div>
    );
}

export default Home;