// src/components/ProfileHeader.js

import React from 'react';
import { handleProfile } from '../services/authService';

const ProfileHeader = ({ profile, isCurrentUser, handleFollow, HandleMessage }) => {
    return (
        <div className="profile-header-container">
            <img id="profile-icon" src={`data:image/jpeg;base64,${profile.iconbase64}`} className="user-icon" alt="Profile" />
            <div className="name">
                <header>
                    <div className="user-name">
                        <p>{profile.name}</p>
                    </div>
                    {profile.followto && <p className='follow-you'>Follow you</p>}
                </header>
            </div>
            <main>
                <div className="user-title">
                    <p>@{profile.username}</p>
                </div>
                <div className="user-bio">
                    <p>{profile.bio}</p>
                </div>
                <div className="create-btn">
                    <div className="posts-count">
                        <p>{profile.countposts}</p>
                        <p id="posts-name">Posts</p>
                    </div>
                    <div className="user-followby">
                        <p>{profile.followbycount}</p>
                        <p id="followers-name">Followers</p>
                    </div>
                    <div className="user-followto">
                        <p>{profile.followtocount}</p>
                        <p id="following-name">Following</p>
                    </div>
                </div>
                <footer>
                    {isCurrentUser && (
                        <button id="edit-profile-btn">Edit Profile</button>
                    )}
                    {!isCurrentUser && (
                        <>
                            <button 
                                id="follow-btn" 
                                onClick={() => handleFollow('follow')}
                                style={{ display: profile.followby ? 'none' : 'block' }}
                            >
                                Follow
                            </button>
                            <button 
                                id="following-btn" 
                                onClick={() => handleFollow('unfollow')}
                                style={{ display: profile.followby ? 'block' : 'none' }}
                            >
                                Following
                            </button>
                            <button 
                                id="message-btn" 
                                onClick={() => HandleMessage(profile.username)}
                            >
                                Message
                            </button>
                        </>
                    )}
                    <div className='items-profile'>
                        <p className='item-1'>Posts</p>
                    </div>
                </footer>
            </main>
        </div>
    );
};

export default ProfileHeader;
