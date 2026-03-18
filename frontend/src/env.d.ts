/**
 * 文件用途：声明 Vite 客户端环境类型，提供 import.meta 的类型提示。
 */
/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
