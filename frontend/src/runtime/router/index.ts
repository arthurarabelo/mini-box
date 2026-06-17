import { createRouter as createVueRouter, createWebHistory } from 'vue-router'
import LoginPage from '../pages/LoginPage.vue'
import AccessDenied from '../pages/AccessDenied.vue'
import AppLayout from '../../AppLayout.vue'

export function createAppRouter() {
  return createVueRouter({
    history: createWebHistory(),
    routes: [
      {
        path: '/login',
        name: 'login',
        component: LoginPage,
        meta: { authContext: 'anonymous' },
      },
      {
        path: '/access-denied',
        name: 'accessDenied',
        component: AccessDenied,
        meta: { authContext: 'anonymous' },
      },
      {
        path: '/',
        name: 'layout',
        component: AppLayout,
        redirect: '/files',
        children: [],
      },
    ],
  })
}