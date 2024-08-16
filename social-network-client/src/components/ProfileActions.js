// src/components/ProfileActions.js

import React from 'react';

const ProfileActions = ({ handleFollow, HandleMessage, isCurrentUser, profile }) => {
    return (
        <div className='profile-btn'>
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
        </div>
    );
};

export default ProfileActions;
