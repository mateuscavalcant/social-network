// src/components/MessageList.js

import React from 'react';
import '../styles/Chat.css';

const MessageList = ({ messages }) => {
  return (
    <ul className="message-list">
      {messages.map((message, index) => (
        <li key={index} className="message-item">
          <span className="message-content">{message.content}</span>
        </li>
      ))}
    </ul>
  );
};

export default MessageList;
