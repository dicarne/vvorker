import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/register',
      component: () => import('@/views/Register.vue')
    },
    {
      path: '/',
      // component: () => import('@/views/Index.vue'),
      redirect: '/admin',
      children: [
        {
          path: '/admin',
          component: () => import('@/views/Admin.vue'),
        },
      ],
    },
  ],
})

export default router
