// src/components/ChatHeader.js

import React from 'react';
import '../styles/Chat.css';

const ChatHeader = ({ chatPartner }) => {
  return (
    <header className="chat-header">
      <img src={chatPartner.iconBase64} alt="Chat partner icon" className="chat-icon" />
      <h2 className="chat-partner-name">{chatPartner.name}</h2>
    </header>
  );
};

export default ChatHeader;
