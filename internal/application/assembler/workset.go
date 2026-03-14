package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/value"
)

func AssembleWorksetInfo(worksetInfo model.WorksetInfo, onLoadURL OnLoadURL) value.WorksetInfo {
	result := value.WorksetInfo{
		ID:          worksetInfo.ID,
		TeamID:      worksetInfo.TeamID,
		Index:       worksetInfo.Index,
		Name:        worksetInfo.Name,
		Description: worksetInfo.Description,
		ComicCount:  worksetInfo.ComicCount,
		CreatedAt:   worksetInfo.CreatedAt.UnixMilli(),
		UpdatedAt:   worksetInfo.UpdatedAt.UnixMilli(),
	}

	if worksetInfo.Team != nil {
		teamInfo := AssembleTeamInfo(*worksetInfo.Team, onLoadURL)
		result.Team = &teamInfo
	}

	return result
}
