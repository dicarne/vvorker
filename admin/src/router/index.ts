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
      redirect: '/workers',
      children: [
        {
          path: '/workers',
          name: 'Workers',
          component: () => import('@/views/Workers.vue'),
        },
        {
          path: '/status',
          name: 'Status',
          component: () => import('@/views/Status.vue'),
        },
        {
          path: '/task',
          name: 'Task',
          component: () => import('@/views/Task.vue'),
        },
        {
          path: '/sql',
          name: 'SQL',
          component: () => import('@/views/SQL.vue'),
        },
        {
          path: '/oss',
          name: 'OSS',
          component: () => import('@/views/OSS.vue'),
        },
        {
          path: '/kv',
          name: 'KV',
          component: () => import('@/views/KV.vue'),
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
