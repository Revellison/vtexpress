import { NavLink } from 'react-router-dom';
import styles from './SideBar.module.css';
import DashboardIcon from './icons/dashboard';
import HistoryIcon from './icons/history';
import SettingsIcon from './icons/settings';
const Sidebar = () => (
    <div className={styles.sidebar}>
        <div className={styles.nav}>
            <nav className={styles.navtop}>
                <NavLink to="/" className={({ isActive }) => isActive ? `${styles.item} ${styles.active}` : styles.item}><DashboardIcon/><span>Главная</span></NavLink>
                <NavLink to="/history" className={({ isActive }) => isActive ? `${styles.item} ${styles.active}` : styles.item}><HistoryIcon/><span>История</span></NavLink>
            </nav>
            <nav className={styles.navbottom}>
                <NavLink to="/settings" className={({ isActive }) => isActive ? `${styles.item} ${styles.active}` : styles.item}><SettingsIcon/><span>Настройки</span></NavLink>
            </nav>
        </div>
    </div>
);

export default Sidebar;