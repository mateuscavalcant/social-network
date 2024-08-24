import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/loginForm.css';


const LoginForm = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            navigate('/home');
        }
    }, [navigate]);

    const handleSubmit = async (e) => {
        e.preventDefault();

        const formData = new FormData();
        formData.append("identifier", email);
        formData.append("password", password);

        try {
            const response = await fetch("http://localhost:8080/login", {
                method: "POST",
                body: formData,
            });

            const data = await response.json();

            if (data.error) {
                setError(data.error.credentials || data.error.password);
            } else {
                console.log(data.message);
                setEmail('');
                setPassword('');
                let token = data.token;
                localStorage.setItem('token', token);

                console.log('Token:', token);
                navigate('/home');
            }
        } catch (error) {
            console.error(error);
        }
    };

    return (
        <div className="container">
            <div className="form">
                <h2 className="header">Sign In</h2>
                <form onSubmit={handleSubmit}>
                    <div className="field">
                        <input
                            type="text"
                            placeholder="Email or username"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                        />
                    </div>
                    <div className="field">
                        <input
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>
                    {error && <span className="error">{error}</span>}
                    <button type="submit" className="button">
                        Login
                    </button>
                </form>
                <div className="form-link">
                    <p>Don't have an account? <a href="/signup" className="link">Signup</a></p>
                </div>
            </div>
        </div>
    );
};

export default LoginForm;
