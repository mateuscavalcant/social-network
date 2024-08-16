import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import './App.css';
import Profile from './pages/Profile';
import Home from './pages/Home';
import LoginForm from './pages/Login';
import SignupForm from './pages/Signup';
import Messages from './pages/Messages';
import ChatMessages from './pages/ChatMessages';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          {/* Redirect from root ("/") to "/home" */}
          <Route path="/" element={<Navigate to="/home" />} />
          <Route path="/login" element={<LoginForm />} />
          <Route path="/signup" element={<SignupForm />} />
          <Route path="/home" element={<Home />} />
          <Route path="/:username" element={<Profile />} />
          <Route path="/chat/:username" element={<ChatMessages />} />
          <Route path="/chats" element={<Messages />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
