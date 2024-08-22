import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import '../styles/Profile.css';
import VerticalNavBar from '../components/VerticalNavBar';
import Post from '../components/Post';

const useProfileData = (username) => {
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
            headers: { Authorization: `Bearer ${token}` }
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

    return { profile, posts, chatPartner, isCurrentUser, loadProfile };
};

const Profile = () => {
    const { username } = useParams();
    const navigate = useNavigate();
    const { profile, posts, chatPartner, isCurrentUser, loadProfile } = useProfileData(username);

    const handleFollow = (action) => {
        if (!username) {
            console.error("Username is missing");
            return;
        }
        const url = action === 'follow' ? 'http://localhost:8080/follow' : 'http://localhost:8080/unfollow';
        axios.post(url, { username }, {
            headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        })
            .then(() => loadProfile(username, localStorage.getItem('token')))
            .catch(error => {
                console.error(`Error ${action}ing user:`, error.response ? error.response.data : error.message);
            });
    };

    const handleMessage = (username) => {
        const token = localStorage.getItem('token');
        axios.post(`http://localhost:8080/chat/${username}`, {}, {
            headers: { Authorization: `Bearer ${token}` }
        })
            .then(() => navigate(`/chat/${username}`))
            .catch(error => {
                console.error("Failed to start chat:", error.response ? error.response.data : error.message);
            });
    };

    if (!profile) {
        return <div>Loading...</div>;
    }

    return (
        <div className="profile-page">
            <div className="bar-btn-container">
                <VerticalNavBar chatPartner={chatPartner} />
            </div>
            <div className="profile-container">
                <div className="profile-header-container">
                    <img id="profile-icon" src={`data:image/jpeg;base64,${profile.iconbase64}`} className="user-icon" alt="Profile" />
                    <div className="page-header">
                        <header>
                            <div className="user-name">
                                <p>{profile.name}</p>
                            </div>
                            {profile.followto && <p className='follow-you'>Follows you</p>}
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
                                            onClick={() => handleMessage(profile.username)}
                                        >
                                            Message
                                        </button>
                                    </>
                                )}
                            </div>
                        </footer>
                    </main>
                </div>
                <div id="posts-container">
                    {posts.map((post) => (
                        <Post key={post.postID} post={post} />
                    ))}
                </div>
            </div>
        </div>
    );
};

export default Profile;
