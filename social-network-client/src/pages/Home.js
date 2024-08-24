// src/pages/Home.js
import React from 'react';
import { useNavigate } from 'react-router-dom';
import { usePosts } from '../hooks/usePosts';
import VerticalNavBar from '../components/VerticalNavBar';
import CreatePostForm from '../components/CreatePostForm';
import Post from '../components/Post';
import '../styles/Home.css';


const Home = () => {
    const navigate = useNavigate();
    const { posts, content, setContent, userInfos, handleCreatePost } = usePosts(navigate);

    if (!posts) {
        return <div>Loading...</div>;
    }

    return (
        <div className="home-page">
            <div className="bar-btn-container">
                <VerticalNavBar userInfos={userInfos} />
            </div>
            <div className="home-container">
                <div className="home-page-header">
                    <header>
                        <p>Home</p>
                    </header>
                </div>
                <CreatePostForm
                    content={content}
                    setContent={setContent}
                    handleCreatePost={handleCreatePost}
                    userInfos={userInfos}
                />
                <div id="posts-container">
                    {posts.map((post) => (
                        <Post key={post.postID} post={post} />
                    ))}
                </div>
            </div>
        </div>
    );
};

export default Home;
