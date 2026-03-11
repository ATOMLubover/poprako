import { request } from "../http";

export async function createMember(
  team_id: string,
  user_id: string,
  roles: number,
) {
  return request("POST", "/members", {
    team_id,
    user_id,
    roles,
  });
}
