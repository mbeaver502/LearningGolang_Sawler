import { createRouter, createWebHistory } from 'vue-router';
import PageBody from './../components/PageBody.vue';
import UserLogin from '../components/UserLogin.vue';

const routes = [
    {
        path: '/',
        name: 'Home',
        component: PageBody,
    },
    {
        path: '/login',
        name: 'Login',
        component: UserLogin,
    },
];

const router = createRouter({ history: createWebHistory(), routes });

export default router;