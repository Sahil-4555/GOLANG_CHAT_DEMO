import React, { useState, useEffect, Suspense } from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import { thunk } from 'redux-thunk';
import reducers from './redux';
import './index.css';
import { legacy_createStore, applyMiddleware, compose } from 'redux';
import { Provider } from 'react-redux';
import WebSocketProvider from './utils/context/WebSocketProvider';
import Loading from './components/Loading';

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

const store = legacy_createStore(
    reducers,
    composeEnhancers(applyMiddleware(thunk))
);

const AppWrapper = () => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            setIsAuthenticated(true);
        }
    }, []);

    return (
        <Provider store={store}>
            {isAuthenticated ? (
                <WebSocketProvider>
                    <Suspense fallback={<Loading />}>
                        <App />
                    </Suspense>
                </WebSocketProvider>
            ) : (
                <Suspense fallback={<Loading />}>
                    <App />
                </Suspense>
            )}
        </Provider>
    );
};

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<AppWrapper />);
