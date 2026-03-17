<template>
  <a-config-provider :theme="antdThemeConfig">
    <router-view />
    <div class="theme-switcher">
      <a-switch
        :checked="isDarkMode"
        checked-children="暗"
        un-checked-children="亮"
        @change="handleThemeModeChange"
      />
    </div>
  </a-config-provider>
</template>

<script setup lang="ts">
/**
 * 根组件负责挂载全局主题 Provider 与路由承载容器。
 */
import { useThemeProvider } from "./theme/provider";

const { isDarkMode, antdThemeConfig, setThemeMode } = useThemeProvider();

/**
 * 处理主题开关状态变化。
 */
function handleThemeModeChange(checked: boolean): void {
  setThemeMode(checked ? "dark" : "light");
}
</script>

<style scoped lang="scss">
/* 主题切换浮层，颜色统一由全局色板变量控制。 */
.theme-switcher {
  position: fixed;
  right: 20px;
  bottom: 20px;
  z-index: 1100;
  padding: 8px 10px;
  border-radius: 999px;
  background: var(--float-bg);
  box-shadow: var(--float-shadow);
  backdrop-filter: blur(10px);
}
</style>
