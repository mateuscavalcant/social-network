import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../styles/Home.css'
import { handleLogout, handleProfile } from '../components/utils';

const Home = () => {
    const [posts, setPosts] = useState([]);
    const [content, setContent] = useState('');
    const [chatPartner, setChatPartner] = useState({ name: '', iconBase64: '' });
    const [token, setToken] = useState('');
    const navigate = useNavigate();


    useEffect(() => {
        const storedToken = localStorage.getItem('token');
        if (storedToken) {
            setToken(storedToken);
            loadPosts(storedToken);
        } else {
            navigate('/login')
        }
    }, [navigate]);

    const loadPosts = (token) => {
        axios.post("http://localhost:8080/feed", {}, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
            .then(response => {
                console.log(token)
                setPosts(response.data.posts);
                setChatPartner(response.data.chatPartner || { name: '', iconBase64: '' });
            })
            .catch(error => {
                console.error("Failed to load posts:", error.response ? error.response.data : error.message);
            });
    };

    const handleCreatePostFormSubmit = (event) => {
        event.preventDefault();

        if (!content.trim()) {
            console.error("Failed to create post: Content is missing!");
            return;
        }

        axios.post("http://localhost:8080/create-post", { content }, {
            headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
            .then(() => {
                setContent('');
                loadPosts(token);
            })
            .catch(error => {
                console.error("Failed to create post:", error.response ? error.response.data : error.message);
            });
    };

    if (!posts) {
        return <div>Loading...</div>;
    }

    return (

        <div className="home-page">
            <div className="bar-btn-container">
                <div className="vertical-bar">
                    <button id="home-btn">
                        <img src="/images/home-solid.png" alt="Home" />
                    </button>
                    <button id="profile-btn">
                        <img
                            src="/images/profile.png"
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
                            alt="Messages"
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

            <div className="home-container">

                <div className="create-post-container2">
                    <p>Home</p>

                </div>


                <div className="create-post-container">
                    <form onSubmit={handleCreatePostFormSubmit}>

                        {chatPartner.iconBase64 && (
                            <img
                                src={`data:image/jpeg;base64,${chatPartner.iconBase64}`}
                                alt="Icon"
                                className="create-post-container-icon"
                                onClick={() => handleProfile(chatPartner.username)}
                                style={{ cursor: 'pointer' }}
                            />
                        )}
                        <textarea
                            value={content}
                            onChange={(e) => setContent(e.target.value)}
                            placeholder="What's happening?"
                            required
                        />
                        <button type="submit">Postar</button>
                    </form>

                </div>
                <div id="posts-container">
                    {posts.map(post => (
                        <div className="post" key={post.postID}>
                            <header>
                                {post.iconbase64 ? (
                                    <img
                                        src={`data:image/jpeg;base64,${post.iconbase64}`}
                                        alt="Profile"
                                        className="profile-icon"
                                        onClick={() => handleProfile(post.createdby)}
                                        style={{ cursor: 'pointer' }}
                                    />
                                ) : (
                                    <img
                                        src="default-profile-icon.png"
                                        alt="Profile"
                                        className="profile-icon"
                                        onClick={() => handleProfile(post.createdby)}
                                        style={{ cursor: 'pointer' }}
                                    />
                                )}
                                <div className="post-title">
                                    <div className="user-name" onClick={() => handleProfile(post.createdby)} style={{ cursor: 'pointer' }}>
                                        <p>{post.createdbyname}</p>
                                    </div>
                                    <div className="user-username" onClick={() => handleProfile(post.createdby)} style={{ cursor: 'pointer' }}>
                                        <p>@{post.createdby}</p>
                                    </div>
                                </div>
                            </header>
                            <main>
                                <div className="post-content">
                                    <p>{post.content}</p>
                                </div>
                                <div className="post-links">

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
export default Home;
