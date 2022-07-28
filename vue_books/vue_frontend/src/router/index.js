import { createRouter, createWebHistory } from 'vue-router';
import PageBody from './../components/PageBody.vue';

const routes = [
    {
        path: '/',
        name: 'Home',
        component: PageBody,
    }
];

const router = createRouter({ history: createWebHistory(), routes });

export default router;