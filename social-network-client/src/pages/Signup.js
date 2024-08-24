import React, { useState } from 'react';
import '../styles/loginForm.css';
import { useNavigate } from 'react-router-dom';


const SignupForm = () => {
    const [name, setName] = useState('');
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleSubmit = () => {
        const formData = new FormData();
        formData.append("name", name);
        formData.append("username", username);
        formData.append("email", email);
        formData.append("password", password);
        formData.append("confirm_password", confirmPassword);


        fetch("http://localhost:8080/signup", {
            method: "POST",
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    setError(data.error.name || data.error.username || data.error.email || data.error.password || data.error.confirm_password);
                } else {
                    console.log(data.message);
                    setName('');
                    setUsername('');
                    setEmail('');
                    setPassword('');
                    setConfirmPassword('');
                    navigate.push('/login');
                }
            })
            .catch(error => {
                console.error(error);
            });
    };

    return (
        <div className="container">
            <div className="form">
                <h2 className="header">Sign Up</h2>
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
                        placeholder="Username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                    />
                </div>
                <div className="field">
                    <input
                        className="input"
                        type="email"
                        placeholder="Email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                    />
                </div>
                <div className="field">
                    <input
                        className="input"
                        type="password"
                        placeholder="Password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>
                <div className="field">
                    <input
                        className="input"
                        type="password"
                        placeholder="Confirm password"
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                    />
                </div>

                {error && <p className="error">{error}</p>}
                <button className="button" onClick={handleSubmit}>
                    Create Account
                </button>
                <div className="form-link">
                    <p>Already have an account? <a href="/login" className="link">Login</a></p>
                </div>
            </div>
        </div>
    );
};

export default SignupForm;
