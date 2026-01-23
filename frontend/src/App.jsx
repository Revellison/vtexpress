import {useState} from 'react';
import SideBar from './components/sidebar/SideBar';
import './App.css';
import { Routes, Route } from 'react-router-dom';
import Home from './pages/home/Home';
import Settings from './pages/settings/Settings';
import History from './pages/history/History';
const App = () => {
    return (
        <div id="App">
            <SideBar />
            <main>
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/history" element={<History />} />
                    <Route path="/settings" element={<Settings />} />
                    <Route path="*" element={<div>Not found</div>} />
                </Routes>
            </main>
        </div>
    );
}

export default App
