import axios from 'axios';
import { useNavigate } from 'react-router-dom';

// Função para lidar com perfil
export const handleProfile = (username, token) => {
    token = localStorage.getItem('token');
    axios.post(`http://localhost:8080/profile/${username}`, {}, {
        headers: {
            Authorization: `Bearer ${token}`
        }
    })
    .then(response => {
        window.location.replace(`/${username}`);
    })
    .catch(error => {
        console.error("Failed to fetch profile:", error.response ? error.response.data : error.message);
    });
};


// Função para lidar com logout
export const handleLogout = () => {
    try {
        localStorage.removeItem('token');
        window.location.replace('/login'); // Ajuste o caminho conforme necessário
    } catch (error) {
        console.error('Error during logout:', error);
    }
};


const useRedirectEditProfile = () => {
    const navigate = useNavigate();

    const handleEditProfile = () => {
        const token = localStorage.getItem('token');
        if (!token) return;

        navigate('/settings/account/edit-profile');
    };

    return handleEditProfile;
};

export default useRedirectEditProfile;