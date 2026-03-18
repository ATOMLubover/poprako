/**
 * 文件用途：定义全局认证状态 Store，统一维护访问令牌与登录态。
 */
import { computed, ref } from "vue";
import { defineStore } from "pinia";

/**
 * 访问令牌在本地存储中的键名。
 */
const ACCESS_TOKEN_STORAGE_KEY = "access_token";

/**
 * 认证 Store。
 */
export const useAuthStore = defineStore("auth", () => {
  /**
   * 当前访问令牌。
   */
  const accessToken = ref<string>(
    localStorage.getItem(ACCESS_TOKEN_STORAGE_KEY) ?? "",
  );

  /**
   * 是否已登录。
   */
  const isLoggedIn = computed<boolean>(() => accessToken.value.length > 0);

  /**
   * 设置访问令牌并同步到本地存储。
   */
  function setAccessToken(nextAccessToken: string | undefined): void {
    const normalizedAccessToken = nextAccessToken ?? "";
    accessToken.value = normalizedAccessToken;
    localStorage.setItem(ACCESS_TOKEN_STORAGE_KEY, normalizedAccessToken);
  }

  /**
   * 清空访问令牌并同步到本地存储。
   */
  function clearAccessToken(): void {
    accessToken.value = "";
    localStorage.removeItem(ACCESS_TOKEN_STORAGE_KEY);
  }

  return {
    accessToken,
    isLoggedIn,
    setAccessToken,
    clearAccessToken,
  };
});
