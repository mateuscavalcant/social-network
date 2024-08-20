import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import '../styles/chats.css';
import { handleLogout } from '../components/utils';

const Messages = () => {
    const [chats, setChats] = useState([]);
    const [chatPartner, setChatPartner] = useState({ name: '', iconBase64: '' });
    const [token, setToken] = useState('');

    useEffect(() => {
        const storedToken = localStorage.getItem('token');
        if (storedToken) {
            setToken(storedToken);
            document.cookie = `token=${storedToken}; path=/; Secure; SameSite=Strict`;
            loadChats(storedToken);
        }
    }, []);

    const loadChats = useCallback(async () => {
        axios.post("http://localhost:8080/chats", {}, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
            .then(response => {
                setChats(response.data.chats);
                setChatPartner(response.data.chatPartner || { name: '', iconBase64: '' });
            })
            .catch(error => {
                console.error("Failed to load posts:", error.response ? error.response.data : error.message);
            });
    });

    const handleProfile = (username) => {
        axios.post(`http://localhost:8080/profile/${username}`, {}, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
            .then(response => {
                window.location.replace(`/${username}`);
            })
            .catch(error => {
                console.error("Failed to fetch profile:", error.response ? error.response.data : error.message);
            });
    };

    const setupWebSocket = useCallback(() => {
        if (!token) return;

        const wsURL = `ws://localhost:8080/chats`;
        const ws = new WebSocket(wsURL);

        ws.onopen = () => {
            console.log('WebSocket connection established.');
        };

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            setChats(data.chats);
        };

        ws.onclose = () => {
            console.log('WebSocket connection closed. Reconnecting...');
            setTimeout(setupWebSocket, 1000);
            loadChats();
        };

        return () => {
            ws.close();
        };
    }, [loadChats, token]);

    useEffect(() => {
        if (token) {
            loadChats();
            setupWebSocket();
        }
    }, [loadChats, token, setupWebSocket]);


    const HandleMessage = (username) => {
        const token = localStorage.getItem('token');
        axios.post(`http://localhost:8080/chat/${username}`, {}, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
            .then(response => {
                window.location.replace(`chat/${username}`);
            })
            .catch(error => {
                console.error("Failed to fetch Chat:", error.response ? error.response.data : error.message);  // Log detalhado de erro
            });
    };

    return (
        <div className='chats-page'>
            <div className="bar-btn-container">
                <div className="vertical-bar">
                    <button id="home-btn">
                        <img
                            src="/images/home.png"
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
                            src="/images/envelope-solid.png"
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
                            alt="Messages"
                            onClick={() => handleLogout()}
                            style={{ cursor: 'pointer' }}
                        />
                    </button>

                </div>
            </div>
            <div className="chats-container">
                <div id="chats-container">
                    {chats.map(post => (
                        <div className="chats" onClick={() => HandleMessage(post.createdby)}
                            style={{ cursor: 'pointer' }} key={post.postID}>
                            <header>
                                {post.iconbase64 ? (
                                    <img
                                        src={`data:image/jpeg;base64,${post.iconbase64}`}
                                        alt="Profile"
                                        className="chats-icon"
                                        onClick={() => HandleMessage(post.createdby)}
                                        style={{ cursor: 'pointer' }}
                                    />
                                ) : (
                                    <img
                                        src="default-profile-icon.png"
                                        alt="Profile"
                                        className="chats-icon"
                                        onClick={() => HandleMessage(post.createdby)}
                                        style={{ cursor: 'pointer' }}
                                    />
                                )}
                                <div className="chats-title">
                                    <div className="chats-name" onClick={() => HandleMessage(post.createdby)} style={{ cursor: 'pointer' }}>
                                        <p>{post.createdbyname}</p>
                                    </div>
                                </div>
                            </header>
                            <main>
                                <div className="chats-main">
                                    <div className="chats-content">
                                        <p>{post.content}</p>
                                    </div>
                                </div>
                                <div className="chats-links">
                                </div>
                            </main>
                            <footer>
                            </footer>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );

};

export default Messages;
