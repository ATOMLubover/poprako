/**
 * 文件用途：Vite 构建配置入口，负责注册 Vue 插件与开发/构建时的解析行为。
 */
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { fileURLToPath, URL } from "node:url";

/**
 * 创建 Vite 配置对象。
 * 这里统一定义 Vue 插件和路径别名，确保工程结构清晰。
 */
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
});
