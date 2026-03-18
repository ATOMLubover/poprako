/**
 * 文件用途：封装认证与当前用户信息相关接口。
 */
import { httpClient } from "../http";
import type { UserInfo } from "../../types/domain";

/**
 * 登录参数，对应 swagger 的 value.LoginUserArgs。
 */
export interface LoginUserArgs {
  /** 登录用 QQ 账号。 */
  qq: string;
  /** 登录密码。 */
  password: string;
}

/**
 * 登录请求体类型。
 */
export type LoginUserRequest = LoginUserArgs;

/**
 * 注册参数，对应 swagger 的 value.RegisterUserArgs。
 */
export interface RegisterUserArgs {
  /** 注册用户名。 */
  username: string;
  /** 注册用 QQ 账号。 */
  qq: string;
  /** 注册密码。 */
  password: string;
}

/**
 * 注册请求体类型。
 */
export type RegisterUserRequest = RegisterUserArgs;

/**
 * 登录结果，对应 swagger 的 value.LoginUserResult。
 */
export interface LoginUserResult {
  /** 登录成功后签发的访问令牌。 */
  access_token: string;
  /** 当前登录用户 ID。 */
  user_id: string;
}

/**
 * 登录响应体类型。
 */
export type LoginUserResponse = LoginUserResult;

/**
 * 注册结果，对应 swagger 的 value.RegisterUserResult。
 */
export interface RegisterUserResult {
  /** 注册后自动登录的访问令牌。 */
  access_token: string;
  /** 新注册用户 ID。 */
  user_id: string;
}

/**
 * 注册响应体类型。
 */
export type RegisterUserResponse = RegisterUserResult;

/**
 * 获取当前用户信息的响应体类型。
 */
export type GetCurrentUserProfileResponse = UserInfo;

/**
 * 预留用户头像上传响应体，对应 value.ReserveUserAvatarResult。
 */
export interface ReserveUserAvatarResult {
  /** 后端签发的预签名 PUT URL。 */
  put_url: string;
}

/**
 * 预留用户头像上传请求体。
 */
export interface ReserveUserAvatarArgs {
  /**
   * 上传请求使用的 Content-Type，需要与后续 PUT 保持一致。
   */
  content_type: string;
}

/**
 * 预留用户头像上传响应类型。
 */
export type ReserveUserAvatarResponse = ReserveUserAvatarResult;

/**
 * 调用 swagger 的 POST /auth/login。
 * 请求类型：LoginUserRequest。
 * 返回类型：LoginUserResponse。
 */
export async function loginUser(
  loginUserArgs: LoginUserRequest,
): Promise<LoginUserResponse> {
  return httpClient.post<LoginUserResponse, LoginUserRequest>(
    "/auth/login",
    loginUserArgs,
  );
}

/**
 * 调用 swagger 的 POST /auth/register。
 * 请求类型：RegisterUserRequest。
 * 返回类型：RegisterUserResponse。
 */
export async function registerUser(
  registerUserArgs: RegisterUserRequest,
): Promise<RegisterUserResponse> {
  return httpClient.post<RegisterUserResponse, RegisterUserRequest>(
    "/auth/register",
    registerUserArgs,
  );
}

/**
 * 调用 swagger 的 GET /users/mine。
 * 请求类型：无。
 * 返回类型：GetCurrentUserProfileResponse。
 */
export async function getCurrentUserProfile(): Promise<GetCurrentUserProfileResponse> {
  return httpClient.get<GetCurrentUserProfileResponse>("/users/mine");
}

/**
 * 调用 swagger 的 POST /users/{user_id}/avatar。
 * 请求类型：无。
 * 返回类型：ReserveUserAvatarResponse。
 */
export async function reserveUserAvatar(
  userID: string,
  reserveUserAvatarArgs: ReserveUserAvatarArgs,
): Promise<ReserveUserAvatarResponse> {
  return httpClient.post<ReserveUserAvatarResponse, ReserveUserAvatarArgs>(
    `/users/${userID}/avatar`,
    reserveUserAvatarArgs,
  );
}

/**
 * 调用 swagger 的 POST /users/{user_id}/avatar/confirm。
 * 请求类型：无。
 * 返回类型：无。
 */
export async function confirmUserAvatarUploaded(userID: string): Promise<void> {
  await httpClient.post<void>(`/users/${userID}/avatar/confirm`);
}
