import React, { useEffect, useMemo, useState } from 'react';
import styles from './Settings.module.css';
import { getSettings, saveSettings } from '../../lib/backend';

const Settings = () => {
    const [vtApiKey, setVtApiKey] = useState('');
    const [aiApiKey, setAiApiKey] = useState('');
    const [aiProvider, setAiProvider] = useState('gemini');
    const [aiModel, setAiModel] = useState('gemini-2.0-flash');
    const [autoAiSummary, setAutoAiSummary] = useState(true);
    const [status, setStatus] = useState('');
    const [isSaving, setIsSaving] = useState(false);

    const canValidate = useMemo(() => {
        return vtApiKey.trim().length > 0 && aiApiKey.trim().length > 0 && aiModel.trim().length > 0;
    }, [vtApiKey, aiApiKey, aiModel]);

    useEffect(() => {
        let isMounted = true;
        getSettings()
            .then((settings) => {
                if (!isMounted) {
                    return;
                }
                if (settings?.aiProvider) {
                    setAiProvider(settings.aiProvider);
                }
                if (settings?.aiModel) {
                    setAiModel(settings.aiModel);
                }
                if (typeof settings?.autoAiSummary === 'boolean') {
                    setAutoAiSummary(settings.autoAiSummary);
                }
            })
            .catch((error) => {
                if (isMounted) {
                    setStatus(error?.message || 'Не удалось загрузить настройки');
                }
            });

        return () => {
            isMounted = false;
        };
    }, []);

    useEffect(() => {
        if (!canValidate) {
            return undefined;
        }

        const timer = setTimeout(async () => {
            setIsSaving(true);
            setStatus('Проверка API ключей...');
            try {
                const response = await saveSettings({
                    vtApiKey,
                    aiApiKey,
                    aiProvider,
                    aiModel,
                    autoAiSummary,
                });
                if (response?.vtValid && response?.aiValid) {
                    setStatus('API ключи проверены и сохранены в защищённое хранилище ОС');
                } else {
                    setStatus(response?.message || 'API ключи не прошли проверку');
                }
            } catch (error) {
                setStatus(error?.message || 'Ошибка валидации API ключей');
            } finally {
                setIsSaving(false);
            }
        }, 2000);

        return () => clearTimeout(timer);
    }, [canValidate, vtApiKey, aiApiKey, aiProvider, aiModel, autoAiSummary]);

    return (
        <div className={styles.settings}>
            <div className={styles.apiSection}>
                <div className={styles.vtApiContainer}>
                    <div className={styles.vtApiContainerTop}>
                        <span>VirusTotal api key</span>
                    </div>
                    <div className={styles.vtApiContainerBottom}>
                        <input
                            className={styles.vtapi}
                            type="password"
                            value={vtApiKey}
                            onChange={(event) => setVtApiKey(event.target.value)}
                            placeholder='Введите VirusTotal API ключ'
                        />
                    </div>
                </div>
                <div className={styles.aiApiContainer}>
                    <div className={styles.aiApiContainerTop}>
                        <span>AI api key</span>
                    </div>
                    <div className={styles.aiApiContainerBottom}>
                        <input
                            className={styles.aiapi}
                            type="password"
                            value={aiApiKey}
                            onChange={(event) => setAiApiKey(event.target.value)}
                            placeholder='Введите AI API ключ'
                        />
                    </div>
                </div>
                <div className={styles.aiApiContainer}>
                    <div className={styles.aiApiContainerTop}>
                        <span>AI provider</span>
                    </div>
                    <div className={styles.aiApiContainerBottom}>
                        <select className={styles.select} value={aiProvider} onChange={(event) => setAiProvider(event.target.value)}>
                            <option value="gemini">Gemini</option>
                            <option value="openrouter">OpenRouter</option>
                        </select>
                    </div>
                </div>
                <div className={styles.aiApiContainer}>
                    <div className={styles.aiApiContainerTop}>
                        <span>AI model</span>
                    </div>
                    <div className={styles.aiApiContainerBottom}>
                        <input
                            className={styles.aiapi}
                            type="text"
                            value={aiModel}
                            onChange={(event) => setAiModel(event.target.value)}
                            placeholder='Например: gemini-2.0-flash или openai/gpt-4o-mini'
                        />
                    </div>
                </div>
                <div className={styles.aiApiContainer}>
                    <div className={styles.aiApiContainerTop}>
                        <span>Auto summary with AI</span>
                    </div>
                    <div className={styles.aiApiContainerBottom}>
                        <label className={styles.checkboxRow}>
                            <input
                                type="checkbox"
                                checked={autoAiSummary}
                                onChange={(event) => setAutoAiSummary(event.target.checked)}
                            />
                            <span>Включить автосуммаризацию после скана</span>
                        </label>
                    </div>
                </div>
                <div className={styles.status}>{isSaving ? 'Сохранение...' : status}</div>
            </div>
        </div>
    );
};

export default Settings;