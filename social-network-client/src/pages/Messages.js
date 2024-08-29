import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import '../styles/chats.css';
import { usePosts } from '../hooks/usePosts';
import { useNavigate } from 'react-router-dom';
import VerticalNavBar from '../components/VerticalNavBar';
import { useHomeChatMessages } from '../hooks/useHomeMessages';



const Messages = () => {
    const navigate = useNavigate()
    const { chats } = useHomeChatMessages();
    const [currentUsername, setCurrentUsername] = useState({ username: '' });
    const [token, setToken] = useState('');
    const { userInfos } = usePosts(navigate);

    const handleMessage = (username) => {
        const token = localStorage.getItem('token');
        axios.post(`http://localhost:8080/chat/${username}`, {}, {
            headers: { Authorization: `Bearer ${token}` }
        })
            .then(() => navigate(`/chat/${username}`))
            .catch(error => {
                console.error("Failed to start chat:", error.response ? error.response.data : error.message);
            });
    };



    return (
        <div className="home-page">
            <div className="bar-btn-container">
                <VerticalNavBar userInfos={userInfos} />
            </div>
            <div className="home-container">
                <div className="home-page-header">
                    <header>
                        <p>Messages</p>
                    </header>
                </div>
                <div id="chats-container">
                    {chats.map((post) => (
                        <div
                            className="post"
                            key={post.postID}
                            onClick={() => handleMessage(post.createdby)}
                            style={{ cursor: 'pointer' }}
                        >
                            <header>
                                {post.iconbase64 ? (
                                    <img
                                        src={`data:image/jpeg;base64,${post.iconbase64}`}
                                        alt="Profile"
                                        className="profile-icon"
                                        onClick={() => handleMessage(post.createdby)}
                                        style={{ cursor: 'pointer' }}
                                    />
                                ) : (
                                    <img
                                        src="default-profile-icon.png"
                                        alt="Profile"
                                        className="profile-icon"
                                        onClick={() => handleMessage(post.createdby)}
                                        style={{ cursor: 'pointer' }}
                                    />
                                )}
                                <div className="post-title">
                                    <div className="user-name" onClick={() => handleMessage(post.createdby)} style={{ cursor: 'pointer' }}>
                                        <p>{post.createdby}</p>
                                    </div>
                                </div>
                            </header>
                            <main>
                                <div className="post-content">
                                    <p>{post.content}</p>
                                </div>
                            </main>
                        </div>
                    ))}

                </div>
            </div>
        </div>
    );

};

export default Messages;
