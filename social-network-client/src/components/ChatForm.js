// src/components/ChatForm.js

import React from 'react';
import '../styles/Chat.css';

const ChatForm = ({ content, onContentChange, onSubmit }) => {
  return (
    <form className="chat-form" onSubmit={onSubmit}>
      <input
        type="text"
        value={content}
        onChange={(e) => onContentChange(e.target.value)}
        className="chat-input"
        placeholder="Type a message..."
      />
      <button type="submit" className="chat-submit-btn">Send</button>
    </form>
  );
};

export default ChatForm;
