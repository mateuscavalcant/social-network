// src/components/Sidebar.js

import React from 'react';
import { handleLogout, handleProfile } from '../services/authService';

const Sidebar = ({ chatPartner }) => {
    return (
        <div className="bar-btn-container">
            <div className="vertical-bar">
                <button id="home-btn">
                    <img src="/images/home.png"
                    alt="Home" 
                    onClick={() => window.location.replace('home')}
                    style={{ cursor: 'pointer' }}
                    />
                </button>
                <button id="profile-btn">
                    <img
                        src="/images/profile.png"
                        alt="Profile"
                        onClick={() => handleProfile(chatPartner.username)} 
                        style={{ cursor: 'pointer' }}
                    />
                </button>
                <button id="search-btn">
                    <img src="/images/search.png" alt="Search" />
                </button>
                <button id='envelope-btn'>
                    <img
                        src="/images/envelope.png"
                        alt="Messages"
                        onClick={() => window.location.replace('chats')}
                        style={{ cursor: 'pointer' }}
                    />
                </button>
                <button id="configure-btn">
                    <img src="/images/config.png" alt="Configure" />
                </button>
                <button id='logout-btn'>
                    <img
                        src="/images/logout.png"
                        alt="Logout"
                        onClick={() => handleLogout()}
                        style={{ cursor: 'pointer' }}
                    />
                </button>
            </div>
        </div>
    );
};

export default Sidebar;
