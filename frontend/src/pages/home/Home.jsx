import styles from './Home.module.css';
import FileArrowIcon from '../../components/icons/filearrow';
const Home = () => {
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
                        <button>Выбрать файл</button>
                        <span>Поддерживаются: .exe, .apk, .zip и другие</span>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Home;