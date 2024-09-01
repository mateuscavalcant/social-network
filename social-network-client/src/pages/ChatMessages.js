// src/pages/ChatMessages.js
import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useChatMessages } from '../hooks/useChatMessages';
import { handleProfile} from '../components/utils';

import Message from '../components/Message';
import MessageForm from '../components/ChatForm';
import '../styles/Chat.css';

const ChatMessages = () => {
  const { username } = useParams(); 
  const { messages, content, setContent, userInfos, handleCreateMessage } = useChatMessages(username);
  const navigate = useNavigate();


  return (
    <div className="chat-page">
      <div className="chat-header">
        <header>
        <div className="header-back">
              <img
              src="/images/arrow-back.png"
              alt='back'
                onClick={() => window.location.replace(`/chats`)}
                style={{ cursor: 'pointer' }}
              >
                
              </img>
            </div>
          <div className="header-home-screen">
            {userInfos.iconBase64 && (
              <img
                src={`data:image/jpeg;base64,${userInfos.iconBase64}`}
                className="profile-icon"
                alt="profile"
                onClick={() => handleProfile(userInfos.username)}
                style={{ cursor: 'pointer' }}
              />
            )}
            <div className="header-name">
              <p
                onClick={() => handleProfile(userInfos.username)}
                style={{ cursor: 'pointer' }}
              >
                {userInfos.name}
              </p>
            </div>
          </div>
        </header>
      </div>
      <div classname= 'messages-container' id="messages-container">
        {Array.isArray(messages) && messages.map((post) => (
          <Message key={post.postid} post={post} />
        ))}
      </div>
      <MessageForm
        content={content}
        setContent={setContent}
        handleCreateMessage={handleCreateMessage}
      />
    </div>
  );
};

export default ChatMessages;
