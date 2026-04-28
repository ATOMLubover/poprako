package svc

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/pkg/util"
)

// `OSS_CRE_EXP` defines the expiration duration for OSS creation messages.
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
	now := time.Now()

	return ossMsgRepo.SavePendingCre(&aggr.OssCreMsg{
		Id:        id,
		ResTyp:    resTyp,
		ResId:     resId,
		Status:    enum.OssMsgStatePending,
		ObjKeys:   keys,
		VisibleAt: now,
		ExpireAt:  exp,
		CreatedAt: now,
	})
}
