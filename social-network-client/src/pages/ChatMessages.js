import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import '../styles/Chat.css';
import { handleLogout, handleProfile } from '../components/utils';

const ChatMessages = () => {
  const [messages, setMessages] = useState([]);
  const [content, setContent] = useState('');
  const [autoScroll, setAutoScroll] = useState(true);
  const [chatPartner, setChatPartner] = useState({ name: '', iconBase64: '', username: '' });
  const [chatPartnerUsername, setChatPartnerUsername] = useState({ username: '' });
  const username = window.location.pathname.split("/").pop();
  const token = localStorage.getItem('token');
  const cookie = document.cookie = `token=${token}; path=/; Secure; SameSite=Strict`;

  const loadPosts = useCallback(async () => {
    try {
      const response = await axios.post(`http://localhost:8080/chat/${username}`, {}, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
      setMessages(response.data.messages || []);
      setChatPartner(response.data.chatPartner || { name: '', iconBase64: '' });
      setChatPartnerUsername(response.data.chatPartnerUsername || { username: '' });

      if (autoScroll) {
        window.scrollTo(0, document.body.scrollHeight);
      }
    } catch (error) {
      console.error('Erro ao carregar posts:', error);
    }
  }, [username, token, autoScroll]);

  const setupWebSocket = useCallback(() => {
    if (!cookie) return;

    const wsURL = `ws://localhost:8080/websocket/${username}`;
    const ws = new WebSocket(wsURL);
    console.log(token);

    ws.onopen = () => {
      console.log('WebSocket connection established.');
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setMessages(data.messages || []);
      setChatPartner(data.chatPartner || { name: '', iconBase64: '' });

      if (autoScroll) {
        window.scrollTo(0, document.body.scrollHeight);
      }
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed. Reconnecting...');
      setTimeout(setupWebSocket, 1000);
      loadPosts();
    };
  }, [username, token, loadPosts, autoScroll]);

  useEffect(() => {
    loadPosts();
    setupWebSocket();
    window.addEventListener('scroll', handleScroll);
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, [loadPosts, setupWebSocket]);

  const handleCreatePost = async (event) => {
    event.preventDefault();
    try {
      await axios.post(`http://localhost:8080/create-message/${username}`, { content }, {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      });
      setContent('');
      loadPosts();
    } catch (error) {
      console.error('Erro ao criar post:', error);
    }
  };

  const handleScroll = () => {
    const scrollTop = window.scrollY;
    const scrollHeight = document.body.scrollHeight;
    const windowHeight = window.innerHeight;

    if (scrollTop + windowHeight < scrollHeight - 5) { // Permite um pequeno buffer
      setAutoScroll(false);
    } else {
      setAutoScroll(true);
    }
  };

  return (
    <div className='chat-page'>
      <div className="chat-bar-btn-container">
        <div className="chat-vertical-bar">
          <button id="home-btn">
            <img
              src="/images/home.png"
              alt="Home"
              onClick={() => window.location.replace('http://localhost:3000/home')}
              style={{ cursor: 'pointer' }}
            />
          </button>
          <button id="profile-btn">
            <img
              src="/images/profile.png"
              alt="Profile"
              onClick={() => handleProfile(chatPartnerUsername.username)}

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
              onClick={() => window.location.replace('http://localhost:3000/chats')}
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
      <div className="home-header">
        <header>
          <div className="header-home-screen">
            {chatPartner.iconBase64 && (
              <img
                src={`data:image/jpeg;base64,${chatPartner.iconBase64}`}
                className="profile-icon"
                alt="profile"
                onClick={() => handleProfile(chatPartner.username)}
                style={{ cursor: 'pointer' }}
              />
            )}
            <div className='header-name'>
              <p onClick={() => handleProfile(chatPartner.username)}
                style={{ cursor: 'pointer' }}>{chatPartner.name}</p>
            </div>
          </div>
        </header>
      </div>
      <div id="messages-container">
        {Array.isArray(messages) && messages.map((post) => (
          <div key={post.postid} className={post.messagesession ? 'message-session' : 'message-to'}>
            <header>
              <div className="message-box">
                <div className="post-content">
                  <p>{post.content}</p>
                </div>
              </div>
            </header>
            <div className="user-name">
              <p>{post.hourminute}</p>
            </div>
          </div>
        ))}
      </div>
      <form className="message-form-create" onSubmit={handleCreatePost}>
        <input
          type="text"
          placeholder="What's happening?"
          value={content}
          onChange={(e) => setContent(e.target.value)}
        />
        <button type="submit">Post</button>
      </form>
    </div>
  );
};

export default ChatMessages;
