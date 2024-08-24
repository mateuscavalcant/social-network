import React, { useState, useEffect } from 'react';
import '../styles/loginForm.css';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { usePosts } from '../hooks/usePosts';
import { handleProfile } from '../components/utils';
import imageCompression from 'browser-image-compression';


const EditProfile = () => {
    const navigate = useNavigate();
    const { userInfos } = usePosts(navigate);
    const [name, setName] = useState('');
    const [bio, setBio] = useState('');
    const [icon, setIcon] = useState(null);
    const [iconPreview, setIconPreview] = useState(null);
    const [error, setError] = useState('');
    const storedToken = localStorage.getItem('token');

    useEffect(() => {
        if (userInfos) {
            setName(userInfos.name);

        }
    }, [userInfos]);

    useEffect(() => {
        if (userInfos.iconBase64) {
            setIconPreview(`data:image/jpeg;base64,${userInfos.iconBase64}`);
        }
    }, [userInfos.iconBase64]);

    const handleFileChange = async (e) => {
        const file = e.target.files[0];
        if (file) {
            try {
                const options = {
                    maxSizeMB: 1,
                    maxWidthOrHeight: 1024,
                };
                const compressedFile = await imageCompression(file, options);
                const previewUrl = URL.createObjectURL(compressedFile);
                setIcon(compressedFile);
                setIconPreview(previewUrl);
            } catch (error) {
                console.error('Error compressing the image:', error);
            }
        }
    };

    const handleSubmit = async () => {
        const formData = new FormData();

        if (name) formData.append("name", name);
        if (bio) formData.append("bio", bio);
        if (icon) formData.append("icon", icon);

        try {
            const response = await axios.post('http://localhost:8080/edit-profile', formData, {
                headers: {
                    'Authorization': `Bearer ${storedToken}`,
                    'Content-Type': 'multipart/form-data'
                }
            });

            if (response.data.error) {
                setError(response.data.error.name || response.data.error.bio || response.data.error.icon);
            } else {
                console.log(response.data.message);
                setName('');
                setBio('');
                setIcon(null);
                setIconPreview(`data:image/jpeg;base64,${response.data.iconBase64}`);
                handleProfile(userInfos.username);
            }
        } catch (error) {
            console.error(error);
            setError('An error occurred while updating the profile.');
        }
    };

    return (
        <div className="container">
            <div className="form">
                <h2 className="header">Edit Profile</h2>

                {iconPreview && (

                    <header>
                        <img
                            id="edit-profile-icon"
                            src={iconPreview || 'default-image-url'}
                            className="edit-profile-icon"
                            alt="Profile"
                        />
                        <input
                            type="file"
                            className="icon"
                            id="icon"
                            accept="image/*"
                            onChange={handleFileChange}
                        />
                        <span id="error-icon" className="error-message">{error}</span>
                    </header>
                )}

                <div className="field">
                    <input
                        className="input"
                        type="text"
                        placeholder="Name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                    />
                </div>
                <div className="field">
                    <input
                        className="input"
                        type="text"
                        placeholder="Bio"
                        value={bio}
                        onChange={(e) => setBio(e.target.value)}
                    />
                </div>

                <button className="button" onClick={handleSubmit}>
                    Save updates
                </button>
            </div>
        </div>
    );
};

export default EditProfile;
