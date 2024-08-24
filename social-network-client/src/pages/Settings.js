import React from "react";
import { useNavigate } from "react-router-dom";
import { handleLogout } from "../components/utils";
import VerticalNavBar from "../components/VerticalNavBar";
import { usePosts } from "../hooks/usePosts";
import '../styles/Settings.css';
import useRedirectEditProfile from "../components/utils";


const Settings = () => {
    const redirectEditProfile = useRedirectEditProfile()
    const navigate = useNavigate()
    const { userInfos } = usePosts(navigate);
    return (

        <div className="home-page">
            <div className="bar-btn-container">
                <VerticalNavBar userInfos={userInfos} />
            </div>
            <div className="home-container">
                <div className="home-page-header">
                    <header>
                        <p>Settings</p>
                    </header>
                </div>
                <div className="settings-container">
                    <div className="bar-btn-settings-container">
                        <div className="vertical-bar-settings">
                            <button
                                id="home-btn"
                                style={{ cursor: 'pointer' }}
                            > Account information
                            </button>
                            <button
                                id="home-btn"
                                style={{ cursor: 'pointer' }}
                            > Change your password
                            </button>
                            <button
                                id="search-btn"
                                style={{ cursor: 'pointer' }}
                                onClick={redirectEditProfile}
                            > Edit your profile

                            </button>
                            <button
                                id='envelope-btn'
                                style={{ cursor: 'pointer' }}
                            > Privacy and safety

                            </button>
                            <button id="configure-btn"
                                style={{ cursor: 'pointer' }}
                            > Notifications
                            </button>
                            <button
                                id='logout-btn'
                                onClick={() => handleLogout()}
                                style={{ cursor: 'pointer' }}
                            >
                                Logout
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>


    );
}


export default Settings;