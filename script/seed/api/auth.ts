import { request, setToken } from "../http";

export async function login(qq: string, password: string) {
  const data = await request("POST", "/auth/login", {
    qq,
    password,
  });

  // API returns `access_token` according to swagger, not `token`.
  setToken(data.access_token);

  return data;
}

export async function registerUser(args: {
  invitation_code: string;
  qq: string;
  name: string;
  password: string;
}) {
  return request("POST", "/auth/register", args);
}
