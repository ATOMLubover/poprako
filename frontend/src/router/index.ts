import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';

/**
 * 路由表定义。
 */
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue'),
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('@/views/DashboardView.vue'),
  },
];

/**
 * 创建路由实例。
 */
const router = createRouter({
  history: createWebHistory(),
  routes,
});

/**
 * 路由前置守卫，用于校验登录态。
 */
router.beforeEach((to) => {
  const token = localStorage.getItem('access_token');
  if (to.path !== '/login' && !token) {
    return '/login';
  }
  if (to.path === '/login' && token) {
    return '/dashboard';
  }
  return true;
});

export default router;
