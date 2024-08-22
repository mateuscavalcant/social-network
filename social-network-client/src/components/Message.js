// src/components/Message.js
import React from 'react';

const Message = ({ post }) => {
  return (
    <div className={post.messagesession ? 'message-session' : 'message-to'}>
      <header>
        <div className="message-box">
          <div className="post-content">
            <p>{post.content}</p>
          </div>
        </div>
      </header>
      <div className="user-name-message">
        <p>{post.hourminute}</p>
      </div>
    </div>
  );
};

export default Message;
