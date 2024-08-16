// src/components/CreatePostForm.js

import React, { useState } from 'react';
import axios from 'axios';

const CreatePostForm = ({ token, loadPosts }) => {
    const [content, setContent] = useState('');

    const handleCreatePostFormSubmit = (event) => {
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

    return (
        <div className="create-post-container">
            <form onSubmit={handleCreatePostFormSubmit}>
                <textarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="What's happening?"
                    required
                />
                <button type="submit">Postar</button>
            </form>
        </div>
    );
};

export default CreatePostForm;
