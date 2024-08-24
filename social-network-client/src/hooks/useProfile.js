import { useState, useEffect } from 'react';
import axios from 'axios';

const useProfile = (username) => {
    const [profile, setProfile] = useState(null);
    const [posts, setPosts] = useState([]);
    const [userInfos, setuserInfos] = useState({ name: '', iconBase64: '' });
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
                setuserInfos(response.data.userInfos || { name: '', iconBase64: '' });
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
            .then(() => window.location.href = `/chat/${username}`)
            .catch(error => {
                console.error("Failed to start chat:", error.response ? error.response.data : error.message);
            });
    };

    return { profile, posts, userInfos, isCurrentUser, handleFollow, handleMessage };
};

export default useProfile;
