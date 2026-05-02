package svc_res

type SvcRes struct {
	code ErrCode
	ok   bool
	msg  string
}

func Accept() SvcRes {
	return SvcRes{
		ok: true,
	}
}

func Reject(code ErrCode, msg string) SvcRes {
	return SvcRes{
		ok:   false,
		code: code,
		msg:  msg,
	}
}

func (r *SvcRes) Code() ErrCode {
	return r.code
}

func (r *SvcRes) Msg() string {
	return r.msg
}

func (r *SvcRes) IsAccept() bool {
	return r.ok
}

func (r *SvcRes) IsReject() bool {
	return !r.ok
}
