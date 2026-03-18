/**
 * 文件用途：前端应用入口，负责初始化 Vue、Pinia、路由与 Ant Design Vue。
 */
import { createApp } from "vue";
import { createPinia } from "pinia";
import Antd from "ant-design-vue";
import "ant-design-vue/dist/reset.css";
import App from "./App.vue";
import router from "./router";
import "./style.scss";

/**
 * 创建并挂载 Vue 应用实例。
 * 该入口统一注入路由与 Ant Design Vue 组件库。
 */
function bootstrapApplication(): void {
  const app = createApp(App);
  const pinia = createPinia();
  app.use(pinia);
  app.use(router);
  app.use(Antd);
  app.mount("#app");
}

bootstrapApplication();
