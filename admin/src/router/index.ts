import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue')
    },
    {
      path: '/',
      name: 'Index',
      component: () => import('@/views/Index.vue'),
      redirect: '/admin',
      children: [
        {
          path: '/admin',
          name: 'Admin',
          component: () => import('@/views/Admin.vue'),
        },
        {
          path: '/user',
          name: 'User',
          component: () => import('@/views/User.vue'),
        },
      ],
    },
  ],
})

export default router
