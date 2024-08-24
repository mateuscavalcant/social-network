// src/hooks/useChatMessages.js
import { useState, useEffect, useCallback } from 'react';
import axios from 'axios';

export const useChatMessages = (username) => {
  const [messages, setMessages] = useState([]);
  const [content, setContent] = useState('');
  const [autoScroll, setAutoScroll] = useState(true);
  const [userInfos, setuserInfos] = useState({ name: '', iconBase64: '', username: '' });
  const token = localStorage.getItem('token');
  const cookie = document.cookie = `token=${token}; path=/; Secure; SameSite=Strict`;

  const loadMessages = useCallback(async () => {
    try {
      const response = await axios.post(`http://localhost:8080/chat/${username}`, {}, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
      setMessages(response.data.messages || []);
      setuserInfos(response.data.userInfos || { name: '', iconBase64: '' });
      
      if (autoScroll) {
        window.scrollTo(0, document.body.scrollHeight);
      }
    } catch (error) {
      console.error('Failed to load messages:', error);
    }
  }, [username, token, autoScroll]);

  const setupWebSocket = useCallback(() => {
    if (!cookie) return;
    const wsURL = `ws://localhost:8080/websocket/${username}`;
    const ws = new WebSocket(wsURL);

    ws.onopen = () => console.log('WebSocket connection established.');

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setMessages(data.messages || []);
      setuserInfos(data.userInfos || { name: '', iconBase64: '' });

    };

    ws.onclose = () => {
      console.log('WebSocket connection closed. Reconnecting...');
      setTimeout(setupWebSocket, 1000);
      loadMessages();
    };
  }, [username, loadMessages]);

  useEffect(() => {
    loadMessages();
    setupWebSocket();

    const handleScroll = () => {
      const scrollTop = window.scrollY;
      const scrollHeight = document.body.scrollHeight;
      const windowHeight = window.innerHeight;

      setAutoScroll(scrollTop + windowHeight >= scrollHeight - 5);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [loadMessages, setupWebSocket]);

  const handleCreateMessage = async (event) => {
    event.preventDefault();
    try {
      await axios.post(`http://localhost:8080/create-message/${username}`, { content }, {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      });
      setContent('');
      loadMessages();
    } catch (error) {
      console.error('Failed to create message:', error);
    }
  };

  return {
    messages,
    content,
    setContent,
    userInfos,
    handleCreateMessage
  };
};
