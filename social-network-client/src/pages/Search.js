import { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { handleProfile } from "../components/utils";
import VerticalNavBar from "../components/VerticalNavBar";
import { usePosts } from "../hooks/usePosts";

const Search = () => {
    const [users, setUsers] = useState([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [error, setError] = useState(null);
    const navigate = useNavigate();
    const userInfos = usePosts(navigate)

    useEffect(() => {
        const storedToken = localStorage.getItem('token');
        if (!storedToken) {
            navigate('/login');
        }
    }, [navigate]);

    const loadUsers = (token) => {
        axios.post("http://localhost:8080/search", { search: searchQuery }, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
        .then(response => {
            console.log("Users state:", users);
            console.log("Response data:", response.data); 
            setUsers(response.data.users || []);
            setError(null);
        })
        .catch(error => {
            console.error("Failed to load users:", error.response ? error.response.data : error.message);
            setError("Failed to load users");
        });
    };
    
    const handleSubmit = (e) => {
        e.preventDefault();
        const storedToken = localStorage.getItem('token');
        if (storedToken) {
            loadUsers(storedToken);
        }
    };

    return (
        <div className="home-page">
            <div className="bar-btn-container">
                <VerticalNavBar userInfos={userInfos} />
            </div>


<div className="home-container">
<div className="create-post-container">
            <form className="message-form-create" onSubmit={handleSubmit}>
                <input
                    type="text"
                    placeholder="Search"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                />
                <button type="submit">Search</button>
            </form>
            </div>

            {error && <p className="error-message">{error}</p>}

            <div id="posts-container">
            {users.length > 0 ? (
                users.map((user) => (
                    <div className="post" key={user.userID}>
                        <header>
                            {user.iconbase64 ? (
                                <img
                                    src={`data:image/jpeg;base64,${user.iconbase64}`}
                                    alt="Profile"
                                    className="profile-icon"
                                    onClick={() => handleProfile(user.username)}
                                    style={{ cursor: 'pointer' }}
                                />
                            ) : (
                                <img
                                    src="default-profile-icon.png"
                                    alt="Profile"
                                    className="profile-icon"
                                    onClick={() => handleProfile(user.username)}
                                    style={{ cursor: 'pointer' }}
                                />
                            )}
                            <div className="post-title">
                                <div className="user-name" onClick={() => handleProfile(user.username)} style={{ cursor: 'pointer' }}>
                                    <p>{user.name}</p>
                                </div>
                                <div className="user-username" onClick={() => handleProfile(user.username)} style={{ cursor: 'pointer' }}>
                                    <p>@{user.username}</p>
                                </div>
                            </div>
                        </header>
                        <main>
                            <div className="post-content">
                                <p>{user.bio}</p>
                            </div>
                        </main>
                    </div>
                    
                ))
            
            ) : (
                <p>No users found</p>
                
            )}
            </div>

            </div>
        </div>
    );
};

export default Search;
