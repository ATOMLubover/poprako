import { createApp } from 'vue';
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/reset.css';
import App from './App.vue';
import router from './router';
import './style.css';

/**
 * 创建并挂载 Vue 应用实例。
 * 该入口统一注入路由与 Ant Design Vue 组件库。
 */
function bootstrapApplication(): void {
  const app = createApp(App);
  app.use(router);
  app.use(Antd);
  app.mount('#app');
}

bootstrapApplication();
