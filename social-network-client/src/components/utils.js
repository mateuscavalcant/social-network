import axios from 'axios';

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


export const handleEditProfile = (token) => {
    token = localStorage.getItem('token');
   if (!token) return;

   window.location.replace(`/editprofile`);


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
