// src/components/PostItem.js

import React from 'react';
import { handleProfile } from '../services/authService';

const PostItem = ({ post }) => {
    return (
        <div className="post" key={post.postID}>
            <header>
                <img 
                    src={`data:image/jpeg;base64,${post.iconbase64 || 'default-profile-icon.png'}`} 
                    alt="Profile" 
                    className="profile-icon" 
                    onClick={() => handleProfile(post.createdby)} 
                    style={{ cursor: 'pointer' }}
                />
                <div className="post-title">
                    <div className="user-name" onClick={() => handleProfile(post.createdby)} style={{ cursor: 'pointer' }}>
                        <p>{post.createdbyname}</p>
                    </div>
                    <div className="user-username" onClick={() => handleProfile(post.createdby)} style={{ cursor: 'pointer' }}>
                        <p>@{post.createdby}</p>
                    </div>
                </div>
            </header>
            <main>
                <div className="post-content">
                    <p>{post.content}</p>
                </div>
            </main>
            <footer>
                {/* Additional footer content */}
            </footer>
        </div>
    );
};

export default PostItem;
