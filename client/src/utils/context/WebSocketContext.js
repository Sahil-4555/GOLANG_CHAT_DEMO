import { createContext } from 'react';

const WebSocketContext = createContext({
    sendMessage: () => { },
    lastMessage: null,
    readyState: 0,
});

export default WebSocketContext;
