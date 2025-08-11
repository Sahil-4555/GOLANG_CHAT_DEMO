import { lazy } from 'react';

const SomethingwentWrong = lazy(() => import('../pages/SomethingwentWrong'));
const Login = lazy(() => import('../pages/Login'));
const Signup = lazy(() => import('../pages/Signup'));
const ChatParent = lazy(() => import('../pages/ChatParent'));

const routes = [
    {
        path: '/login',
        exact: true,
        name: 'Log In',
        component: Login,
        private: false,
    },
    {
        path: '/signup',
        exact: true,
        name: 'Sign Up',
        component: Signup,
        private: false,
    },
    {
        path: '/error',
        exact: true,
        name: 'Error',
        component: SomethingwentWrong,
        private: false,
    },
    {
        path: '/chat',
        exact: true,
        name: 'Chat Screen',
        component: ChatParent,
        private: true,
    },
    {
        path: '/',
        exact: true,
        name: 'Chat Screen',
        component: ChatParent,
        private: true,
    }
];

export default routes;