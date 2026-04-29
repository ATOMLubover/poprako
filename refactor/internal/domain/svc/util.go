package svc

import (
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
)

func classifyRepoErr(
	err repo_iface.RepoErr,
	clsf repo_iface.ErrClsf,
	notFoundCode svc_res.ErrCode,
	notFoundMsg, timeoutMsg, unavailableMsg, serverMsg string,
) svc_res.SvcRes {
	if clsf.IsNotFound(err) {
		return svc_res.Reject(notFoundCode, notFoundMsg)
	}
	if clsf.IsTimeout(err) {
		return svc_res.Reject(svc_res.Timeout, timeoutMsg)
	}
	if clsf.IsUnavailable(err) {
		return svc_res.Reject(svc_res.Unavailable, unavailableMsg)
	}

	return svc_res.Reject(svc_res.ServerError, serverMsg)
}
