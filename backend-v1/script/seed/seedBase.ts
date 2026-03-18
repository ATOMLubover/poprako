import { CONFIG } from "./config";

import { login, registerUser } from "./api/auth";
import { createTeam } from "./api/team";
import { createInvitation } from "./api/invitation";
import { createMember } from "./api/member";

import { basicRoles } from "./util/roles";
import { genQQ } from "./util/qq";

async function main() {
  console.log("login super admin");

  const admin = await login(CONFIG.superAdmin.qq, CONFIG.superAdmin.password);

  console.log("admin login ok");

  const team = await createTeam(CONFIG.team.name, CONFIG.team.description);

  console.log("team created", team.id);

  // Ensure super admin is a team admin (roles=64) inside the team
  const superMember = await createMember(
    team.id,
    admin.user_id ?? admin.userId ?? admin.id,
    64,
  );
  console.log(
    "super admin added to team as member (roles=64)",
    superMember.id ??
      superMember.member_id ??
      superMember.memberId ??
      superMember,
  );

  const roles = basicRoles();

  const invitations = [];

  for (let i = 0; i < roles.length; i++) {
    const qq = genQQ(i + 1);

    const inv = await createInvitation({
      team_id: team.id,
      invitee_qq: qq,
      roles: roles[i],
    });

    console.log(
      "invitation created",
      inv.id ?? inv.invitation_id ?? inv.invitationId ?? inv.invitation_code,
    );

    invitations.push({
      id: inv.id ?? inv.invitation_id ?? inv.invitationId,
      code: inv.invitation_code,
      qq,
      role: roles[i],
    });
  }

  console.log("invitations created");

  for (const inv of invitations) {
    const user = await registerUser({
      invitation_code: inv.code,
      qq: inv.qq,
      name: `user_${inv.role}`,
      password: CONFIG.defaultPassword,
    });

    console.log(
      `user created qq=${inv.qq} role=${inv.role}`,
      user.id ?? user.user_id ?? user.userId ?? user,
    );
  }

  console.log("seed finished");
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
