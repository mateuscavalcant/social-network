// src/hooks/usePosts.js
import { useState, useEffect } from 'react';
import axios from 'axios';

export const usePosts = (navigate) => {
    const [posts, setPosts] = useState([]);
    const [content, setContent] = useState('');
    const [chatPartner, setChatPartner] = useState({ name: '', iconBase64: '' });
    const [token, setToken] = useState('');

    useEffect(() => {
        const storedToken = localStorage.getItem('token');
        if (storedToken) {
            setToken(storedToken);
            loadPosts(storedToken);
        } else {
            navigate('/login');
        }
    }, [navigate]);

    const loadPosts = (token) => {
        axios.post("http://localhost:8080/feed", {}, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
        .then(response => {
            setPosts(response.data.posts);
            setChatPartner(response.data.chatPartner || { name: '', iconBase64: '' });
        })
        .catch(error => {
            console.error("Failed to load posts:", error.response ? error.response.data : error.message);
        });
    };

    const handleCreatePost = (event) => {
        event.preventDefault();

        if (!content.trim()) {
            console.error("Failed to create post: Content is missing!");
            return;
        }

        axios.post("http://localhost:8080/create-post", { content }, {
            headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
        .then(() => {
            setContent('');
            loadPosts(token);
        })
        .catch(error => {
            console.error("Failed to create post:", error.response ? error.response.data : error.message);
        });
    };

    return {
        posts,
        content,
        setContent,
        chatPartner,
        handleCreatePost
    };
};
