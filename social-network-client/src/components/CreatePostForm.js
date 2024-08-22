// src/components/CreatePostForm.js
import React from 'react';
import { handleProfile } from './utils';

const CreatePostForm = ({ content, setContent, handleCreatePost, chatPartner }) => {
    return (
        <div className="create-post-container">
            <form onSubmit={handleCreatePost}>
                {chatPartner.iconBase64 && (
                    <img
                        src={`data:image/jpeg;base64,${chatPartner.iconBase64}`}
                        alt="Icon"
                        className="create-post-container-icon"
                        onClick={() => handleProfile(chatPartner.username)}
                        style={{ cursor: 'pointer' }}
                    />
                )}
                <textarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="What's happening?"
                    required
                />
                <button type="submit">Post</button>
            </form>
        </div>
    );
};

export default CreatePostForm;
