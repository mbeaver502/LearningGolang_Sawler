import { createRouter, createWebHistory } from 'vue-router';
import PageBody from './../components/PageBody.vue';
import UserLogin from '../components/UserLogin.vue';
import BooksView from '../components/BooksView.vue';
import BookView from '../components/BookView.vue';
import BooksAdmin from './../components/BooksAdmin.vue';
import BookEdit from './../components/BookEdit.vue';
import UsersView from '../components/UsersView.vue';
import UserEdit from './../components/UserEdit.vue';

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
    {
        path: '/books',
        name: 'Books',
        component: BooksView,
    },
    {
        path: '/books/:bookName',
        name: 'Book',
        component: BookView,
    },
    {
        path: '/admin/books',
        name: 'Books Admin',
        component: BooksAdmin,
    },
    {
        path: '/admin/books/:bookId',
        name: 'Book Edit',
        component: BookEdit,
    },
    {
        path: '/admin/users',
        name: 'Users',
        component: UsersView,
    },
    {
        path: '/admin/users/:userId',
        name: 'User Edit',
        component: UserEdit,
    },
];

const router = createRouter({ history: createWebHistory(), routes });

export default router;