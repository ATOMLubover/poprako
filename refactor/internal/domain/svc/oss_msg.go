package svc

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/pkg/util"
)

const OSS_CRE_EXP = 30 * time.Minute

type OssMsgSvc struct{}

func NewOssMsgSvc() OssMsgSvc {
	return OssMsgSvc{}
}

func (OssMsgSvc) SavePendingCre(
	ossMsgRepo repo_iface.OssMsgRepo,
	resTyp enum.OssResTyp,
	resId string,
	keys []string,
) error {
	id := util.GenId("oss_msg")
	exp := time.Now().Add(OSS_CRE_EXP)

	return ossMsgRepo.SavePendingCre(&aggr.OssCreMsg{
		Id:        id,
		ResTyp:    resTyp,
		ResId:     resId,
		Status:    enum.OssMsgStatePend,
		ObjKeys:   keys,
		VisibleAt: time.Now(),
		ExpireAt:  exp,
		CreatedAt: time.Now(),
	})
}
