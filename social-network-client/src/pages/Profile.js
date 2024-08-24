import axios from 'axios';
import { useParams } from 'react-router-dom';
import '../styles/Profile.css';
import VerticalNavBar from '../components/VerticalNavBar';
import Post from '../components/Post';
import useProfile from '../hooks/useProfile';
import useRedirectEditProfile from '../components/utils';


const Profile = () => {
    const redirectEditProfile = useRedirectEditProfile()
    const { username } = useParams();
    const { profile, posts, userInfos, isCurrentUser, handleFollow, handleMessage } = useProfile(username);


    if (!profile) {
        return <div>Loading...</div>;
    }

    return (
        <div className="profile-page">
            <div className="bar-btn-container">
                <VerticalNavBar userInfos={userInfos} />
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
                                <button id="edit-profile-btn"
                                    onClick={redirectEditProfile}
                                    style={{ cursor: 'pointer' }}>Edit Profile
                                </button>
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
