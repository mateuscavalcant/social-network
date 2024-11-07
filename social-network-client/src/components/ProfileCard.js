// src/components/Post.js
import React from 'react';
import { handleProfile } from './utils';

const ProfileCard = ({ user }) => (
    <div className="post" key={user.postID}>
        <header>
            {user.iconbase64 ? (
                <img
                    src={`data:image/jpeg;base64,${user.iconbase64}`}
                    alt="Profile"
                    className="profile-icon"
                    onClick={() => handleProfile(user.username)}
                    style={{ cursor: 'pointer' }}
                />
            ) : (
                <img
                    src="default-profile-icon.png"
                    alt="Profile"
                    className="profile-icon"
                    onClick={() => handleProfile(user.username)}
                    style={{ cursor: 'pointer' }}
                />
            )}
            <div className="post-title">
                <div className="user-name" onClick={() => handleProfile(user.username)} style={{ cursor: 'pointer' }}>
                    <p>{user.name}</p>
                </div>
                <div className="user-username" onClick={() => handleProfile(user.username)} style={{ cursor: 'pointer' }}>
                    <p>@{user.username}</p>
                </div>
            </div>
        </header>
        <main>
            <div className="post-content">
                <p>{user.bio}</p>
            </div>
        </main>
    </div>
);

export default ProfileCard;
