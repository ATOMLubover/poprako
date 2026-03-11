import { request } from "../http";

export async function createTeam(name: string, description: string) {
  return request("POST", "/teams", {
    name,
    description,
  });
}
