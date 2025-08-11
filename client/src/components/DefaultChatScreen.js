import React from 'react';
import '../assets/style/components/DefaultChatScreen.css'; 

const DefaultChatScreen = () => {
    return (
        <div className='default-chat-screen'>
            <div className='welcome-text-container'>
                <h2 className='welcome-title'>
                    Welcome to Chat Application
                </h2>
                <p className='welcome-subtitle'>
                    Select a channel from the sidebar to start chatting.
                </p>
            </div>
            <div className='icon-container-default-screen'>
                <svg
                    xmlns='http://www.w3.org/2000/svg'
                    fill='none'
                    viewBox='0 0 24 24'
                    strokeWidth={1.5}
                    stroke='currentColor'
                    className='chat-icon'
                >
                    <path
                        strokeLinecap='round'
                        strokeLinejoin='round'
                        d='M12 11c1.104 0 2-.896 2-2s-.896-2-2-2-2 .896-2 2 .896 2 2 2zm0 2c-2.667 0-8 1.333-8 4v2h16v-2c0-2.667-5.333-4-8-4z'
                    />
                </svg>
            </div>
        </div>
    );
}

export default DefaultChatScreen;

