import { httpClient } from '@/api/http';
import type { UserInfo } from '@/types/domain';

/**
 * 登录参数，对应 swagger 的 value.LoginUserArgs。
 */
export interface LoginUserArgs {
  qq: string;
  password: string;
}

/**
 * 注册参数，对应 swagger 的 value.RegisterUserArgs。
 */
export interface RegisterUserArgs {
  username: string;
  qq: string;
  password: string;
}

/**
 * 登录结果，对应 swagger 的 value.LoginUserResult。
 */
export interface LoginUserResult {
  access_token: string;
  user: UserInfo;
}

/**
 * 注册结果，对应 swagger 的 value.RegisterUserResult。
 */
export interface RegisterUserResult {
  access_token: string;
  user: UserInfo;
}

/**
 * 调用 swagger 的 POST /auth/login。
 */
export async function loginUser(loginUserArgs: LoginUserArgs): Promise<LoginUserResult> {
  return httpClient.post<LoginUserResult, LoginUserArgs>('/auth/login', loginUserArgs);
}

/**
 * 调用 swagger 的 POST /auth/register。
 */
export async function registerUser(registerUserArgs: RegisterUserArgs): Promise<RegisterUserResult> {
  return httpClient.post<RegisterUserResult, RegisterUserArgs>('/auth/register', registerUserArgs);
}

/**
 * 调用 swagger 的 GET /users/mine。
 */
export async function getCurrentUserProfile(): Promise<UserInfo> {
  return httpClient.get<UserInfo>('/users/mine');
}
