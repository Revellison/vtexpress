import {useState} from 'react';
import SideBar from './components/sidebar/SideBar';
import './App.css';
import { Routes, Route } from 'react-router-dom';

const App = () => {
    return (
        <div id="App">
            <SideBar />
            <main>
                <Routes>
                    <Route path="/" element={<div>Home</div>} />
                    <Route path="/history" element={<div>History</div>} />
                    <Route path="/settings" element={<div>Settings</div>} />
                    <Route path="*" element={<div>Not found</div>} />
                </Routes>
            </main>
        </div>
    );
}

export default App
