/**
 * 文件用途：提供全局主题状态管理，统一控制亮色/暗黑模式与 Ant Design Vue 主题配置。
 */
import { computed, ref, watch, type ComputedRef, type Ref } from "vue";
import { theme as antdTheme, type ConfigProviderProps } from "ant-design-vue";

/**
 * Ant Design Vue 主题配置类型。
 */
type AntdThemeConfig = NonNullable<ConfigProviderProps["theme"]>;

/**
 * 主题持久化键名，用于在本地存储中记忆用户主题偏好。
 */
const THEME_STORAGE_KEY = "poprako_theme_mode";

/**
 * 主题模式枚举值。
 */
type ThemeMode = "light" | "dark";

/**
 * 全局主题 Provider 暴露的数据结构。
 */
export interface ThemeProviderState {
  /** 当前是否处于暗黑模式。 */
  isDarkMode: Ref<boolean>;
  /** Ant Design Vue 的主题配置对象。 */
  antdThemeConfig: ComputedRef<AntdThemeConfig>;
  /** 切换亮色/暗黑模式。 */
  toggleThemeMode: () => void;
  /** 显式设置主题模式。 */
  setThemeMode: (themeMode: ThemeMode) => void;
}

/**
 * 从浏览器环境读取初始主题模式。
 */
function resolveInitialThemeMode(): ThemeMode {
  const cachedThemeMode = localStorage.getItem(THEME_STORAGE_KEY);
  if (cachedThemeMode === "light" || cachedThemeMode === "dark") {
    return cachedThemeMode;
  }
  const preferDarkMode = window.matchMedia(
    "(prefers-color-scheme: dark)",
  ).matches;
  return preferDarkMode ? "dark" : "light";
}

/**
 * 创建全局主题 Provider。
 */
export function useThemeProvider(): ThemeProviderState {
  const isDarkMode = ref(resolveInitialThemeMode() === "dark");

  const antdThemeConfig = computed<AntdThemeConfig>(() => ({
    algorithm: isDarkMode.value
      ? antdTheme.darkAlgorithm
      : antdTheme.defaultAlgorithm,
    token: {
      colorPrimary: "#006d77",
      borderRadius: 12,
    },
  }));

  /**
   * 设置主题模式并触发持久化。
   */
  function setThemeMode(themeMode: ThemeMode): void {
    isDarkMode.value = themeMode === "dark";
  }

  /**
   * 在亮色与暗黑模式之间切换。
   */
  function toggleThemeMode(): void {
    setThemeMode(isDarkMode.value ? "light" : "dark");
  }

  watch(
    isDarkMode,
    (nextDarkMode) => {
      const currentThemeMode: ThemeMode = nextDarkMode ? "dark" : "light";
      localStorage.setItem(THEME_STORAGE_KEY, currentThemeMode);
      document.documentElement.dataset.theme = currentThemeMode;
    },
    { immediate: true },
  );

  return {
    isDarkMode,
    antdThemeConfig,
    toggleThemeMode,
    setThemeMode,
  };
}
