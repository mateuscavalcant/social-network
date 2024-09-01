// src/components/VerticalNavBar.js
import React from 'react';
import {handleProfile } from './utils';
import { useNavigate } from 'react-router-dom';

const VerticalNavBar = ({ userInfos }) => {
    const navigate = useNavigate()
    return (
        <div className="vertical-bar">
            <button id="home-btn">
                <img src="/images/home.png"
                    alt="Home"
                    onClick={() => window.location.replace(`/`)}
                    style={{ cursor: 'pointer' }}
                />
            </button>
            <button id="profile-btn">
                <img
                    src="/images/profile.png"
                    alt="Profile"
                    onClick={() => handleProfile(userInfos.username)}
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
                    onClick={() => navigate('/chats')}
                    style={{ cursor: 'pointer' }}
                />
            </button>
            <button id="configure-btn">
                <img 
                src="/images/config.png" 
                alt="Configure" 
                onClick={() => navigate('/settings/account')}
                style={{ cursor: 'pointer' }}/>
            </button>
        </div>
    );
};

export default VerticalNavBar;
