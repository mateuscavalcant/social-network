import React, { useEffect, useState } from 'react';
import axios from 'axios';

const HomeTeste = () => {
    const [message, setMessage] = useState('');

    useEffect(() => {
        fetchProtectedMessage();
    }, []);

    const fetchProtectedMessage = async () => {
        try {
            const token = await localStorage.getItem('token'); // Use localStorage instead of AsyncStorage
            console.log('Token:', token);
            const response = await axios.post("http://localhost:8000/protected", {}, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });

            if (response.status === 200) {
                setMessage(response.data.message);
            } else {
                setMessage('User not logged');
            }
        } catch (error) {
            console.error(error);
            setMessage('User not logged');
        }
    };

    return (
        <div style={styles.container}>
            <p style={styles.message}>{message}</p>
        </div>
    );
};

const styles = {
    container: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh', // Ensure full height of the viewport
    },
    message: {
        fontSize: 20,
    },
};

export default HomeTeste;
