/**
 * 文件用途：前端路由配置，定义页面入口与登录态守卫逻辑。
 */
import {
  createRouter,
  createWebHistory,
  type RouteRecordRaw,
} from "vue-router";
import { useAuthStore } from "../stores/auth";

/**
 * 路由表定义。
 */
const routes: RouteRecordRaw[] = [
  {
    path: "/",
    redirect: "/dashboard",
  },
  {
    path: "/login",
    name: "login",
    component: () => import("../views/LoginView.vue"),
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: () => import("../views/DashboardView.vue"),
  },
  {
    path: "/file-test",
    name: "file-test",
    component: () => import("../views/FileTransferTestView.vue"),
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
  const authStore = useAuthStore();
  if (to.path !== "/login" && !authStore.isLoggedIn) {
    return "/login";
  }
  if (to.path === "/login" && authStore.isLoggedIn) {
    return "/dashboard";
  }
  return true;
});

export default router;
