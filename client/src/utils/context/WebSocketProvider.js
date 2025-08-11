import useWebSocket from 'react-use-websocket';
import WebSocketContext from './WebSocketContext';
import React, { useEffect, useState } from 'react';

const WebSocketProvider = ({ children }) => {
    const [url, setUrl] = useState(null);

    useEffect(() => {
        const userId = localStorage.getItem("userId");
        const token = localStorage.getItem('token');
        if (userId && token) {
            setUrl(`${process.env.REACT_APP_WEB_SOCKET}?_id=${userId}&token=Bearer ${token}`);
        }
    }, []);

    const { sendMessage, lastMessage, readyState } = useWebSocket(url, {
        shouldReconnect: (closeEvent) => true,
    });

    return (
        <WebSocketContext.Provider value={{ sendMessage, lastMessage, readyState }}>
            {children}
        </WebSocketContext.Provider>
    );
};

export default WebSocketProvider;
