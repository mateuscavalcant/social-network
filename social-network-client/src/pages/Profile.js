import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import '../styles/Profile.css';
import { handleLogout, handleProfile } from '../components/utils';

const Profile = () => {
    const { username } = useParams();
    // eslint-disable-next-line
    const navigate = useNavigate();
    const [profile, setProfile] = useState(null);
    const [posts, setPosts] = useState([]);
    const [chatPartner, setChatPartner] = useState({ name: '', iconBase64: '' });
    const [isCurrentUser, setIsCurrentUser] = useState(false);

    useEffect(() => {
        const storedToken = localStorage.getItem('token');
        if (storedToken) {
            loadProfile(username, storedToken);
        }
    }, [username]);

    const loadProfile = (username, token) => {
        axios.post(`http://localhost:8080/profile/${username}`, {}, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
            .then(response => {
                setProfile(response.data.profile);
                setIsCurrentUser(response.data.isCurrentUser);
                document.title = `${response.data.profile.name} / (@${response.data.profile.username})`;
                setPosts(response.data.posts);
                setChatPartner(response.data.chatPartner || { name: '', iconBase64: '' });
            })
            .catch(error => {
                console.error("Failed to fetch profile:", error.response ? error.response.data : error.message);
            });
    };

    const handleFollow = (action) => {
        if (!username) {
            console.error("Username is missing");
            return;
        }
        const url = action === 'follow' ? 'http://localhost:8080/follow' : 'http://localhost:8080/unfollow';
        axios.post(url, { username: username }, {
            headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`
            }
        })
            .then(response => {
                loadProfile(username, localStorage.getItem('token'));
                console.log(`${action}ed successfully:`, response);
            })
            .catch(error => {
                console.error(`Error ${action}ing user:`, error.response ? error.response.data : error.message);
            });
    };

    if (!profile) {
        return <div>Loading...</div>;
    }

    const HandleMessage = (username) => {
        const token = localStorage.getItem('token');
        axios.post(`http://localhost:8080/chat/${username}`, {}, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
            .then(response => {
                window.location.replace(`chat/${username}`);
            })
            .catch(error => {
                console.error("Failed to fetch Chat:", error.response ? error.response.data : error.message);  // Log detalhado de erro
            });
    };


    return (
        <div className="profile-page">
            <div className="bar-btn-container">
                <div className="vertical-bar">
                    <button id="home-btn">
                        <img
                            src="/images/home.png"
                            alt="Home"
                            onClick={() => window.location.replace('home')}
                            style={{ cursor: 'pointer' }}
                        />
                    </button>
                    <button id="profile-btn">
                        <img
                            src={isCurrentUser ? "/images/profile-solid.png" : "/images/profile.png"}
                            alt="Profile"
                            onClick={() => handleProfile(chatPartner.username)}
                            style={{ cursor: 'pointer' }}
                        />
                    </button>
                    <button id="search-btn">
                        <img src="/images/search.png" alt="Search" />
                    </button>
                    <button id='envelope-btn'>
                        <img
                            src="/images/envelope.png"
                            alt="Home"
                            onClick={() => window.location.replace('chats')}
                            style={{ cursor: 'pointer' }}
                        />
                    </button>
                    <button id="configure-btn">
                        <img src="/images/config.png" alt="Configure" />
                    </button>
                    <button id='logout-btn'>
                        <img
                            src="/images/logout.png"
                            alt="Messages"
                            onClick={() => handleLogout()}
                            style={{ cursor: 'pointer' }}
                        />
                    </button>
                </div>
            </div>
            <div className="profile-container">
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

                            <div className='items-profile'>
                                <p className='item-1'>Posts</p>
                            </div>

                        </footer>
                    </main>

                </div>
                <div id="posts-container">
                    {posts.map(post => (
                        <div className="post" key={post.postID}>
                            <header>
                                <img src={`data:image/jpeg;base64,${post.iconbase64}`} alt="Profile" className="profile-icon" />
                                <div className="post-title">
                                    <div className="user-name-post">
                                        <p className={`name-user${post.postID}`}>{post.createdbyname}</p>
                                    </div>
                                    <div className="user-username">
                                        <p className={`username-user${post.postID}`}>@{post.createdby}</p>
                                    </div>
                                </div>
                            </header>
                            <main>
                                <div className="post-content">
                                    <p>{post.content}</p>
                                </div>
                            </main>
                            <footer>

                            </footer>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default Profile;