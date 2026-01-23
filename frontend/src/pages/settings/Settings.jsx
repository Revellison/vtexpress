import React from 'react';
import styles from './Settings.module.css';
import CopyIcon from '../../components/icons/copyicon';
import EyeToggle from '../../components/icons/eyeToggle';
const Settings = () => {
    return (
        <div className={styles.settings}>
            <div className={styles.apiSection}>
                <div className={styles.vtApiContainer}>
                    <div className={styles.vtApiContainerTop}>
                        <span>VirusTotal api key</span>
                    </div>
                    <div className={styles.vtApiContainerBottom}>
                        <input className={styles.vtapi} type="text" placeholder='xxxxx-xxxxx-xxxx'/>
                        {/*<button className={styles.vtApiSave}><EyeToggle /></button> */}
                    </div>
                </div>
                <div className={styles.aiApiContainer}>
                    <div className={styles.aiApiContainerTop}>
                        <span>Gemini api key</span>
                    </div>
                    <div className={styles.aiApiContainerBottom}>  
                        <input className={styles.aiapi} type="text" placeholder='xxxxx-xxxxx-xxxx'/>
                        {/*<button className={styles.aiApiSave}><EyeToggle /></button>*/}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Settings;