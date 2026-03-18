import { request } from "../http"

export async function createInvitation(args: {
  team_id: string
  invitee_qq: string
  roles: number
}) {
  return request("POST", "/invitations", args)
}