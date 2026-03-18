package service

import "strings"

func GenerateTeamAvatarOSSKey(teamID string) string {
	return strings.Join([]string{"team-avatar", teamID}, "_")
}
