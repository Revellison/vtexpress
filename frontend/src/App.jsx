import {useState} from 'react';
import SideBar from './components/sidebar/SideBar';
import './App.css';
import { Routes, Route } from 'react-router-dom';
import Home from './pages/home/Home';
import Settings from './pages/settings/Settings';
import History from './pages/history/History';
const App = () => {
    const [latestScan, setLatestScan] = useState(null);

    return (
        <div id="App">
            <SideBar />
            <main>
                <Routes>
                    <Route path="/" element={<Home latestScan={latestScan} onScanComplete={setLatestScan} />} />
                    <Route path="/history" element={<History latestScan={latestScan} />} />
                    <Route path="/settings" element={<Settings />} />
                    <Route path="*" element={<div>Not found</div>} />
                </Routes>
            </main>
        </div>
    );
}

export default App
