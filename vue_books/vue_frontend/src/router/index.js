import { createRouter, createWebHistory } from 'vue-router';
import PageBody from './../components/PageBody.vue';
import UserLogin from '../components/UserLogin.vue';
import BooksView from '../components/BooksView.vue';
import BookView from '../components/BookView.vue';
import BooksAdmin from './../components/BooksAdmin.vue';
import BookEdit from './../components/BookEdit.vue';
import UsersView from '../components/UsersView.vue';
import UserEdit from './../components/UserEdit.vue';
import Security from '@/components/security';

const routes = [
    {
        path: '/',
        name: 'Home',
        component: PageBody,
    },
    {
        path: '/login',
        name: 'UserLogin',
        component: UserLogin,
    },
    {
        path: '/books',
        name: 'BooksView',
        component: BooksView,
    },
    {
        path: '/books/:bookName',
        name: 'BookView',
        component: BookView,
    },
    {
        path: '/admin/books',
        name: 'BooksAdmin',
        component: BooksAdmin,
    },
    {
        path: '/admin/books/:bookId',
        name: 'BookEdit',
        component: BookEdit,
    },
    {
        path: '/admin/users',
        name: 'UsersView',
        component: UsersView,
    },
    {
        path: '/admin/users/:userId',
        name: 'UserEdit',
        component: UserEdit,
    },
];

const router = createRouter({ history: createWebHistory(), routes });

router.beforeEach(() => {
    Security.checkToken();
});

export default router;