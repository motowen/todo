package model

import (
	"net/http"
)

type ServiceResp struct {
	Status int           `json:"status"`
	ErrMsg ServiceErrMsg `json:"errMsg"`
}

type ServiceErrMsg struct {
	Msg string `json:"msg"`
}

type serviceError struct {
	OK                    ServiceResp
	Accepted              func(string) ServiceResp
	NoContent             ServiceResp
	Found                 func(string) ServiceResp
	NotModified           func(string) ServiceResp
	BadRequestError       func(string) ServiceResp
	ForbiddenError        func(string) ServiceResp
	NotFoundError         ServiceResp
	FailedDependencyError func(string) ServiceResp
	InternalServiceError  func(string) ServiceResp
}

var ServiceError = serviceError{
	OK: ServiceResp{
		http.StatusOK, ServiceErrMsg{http.StatusText(http.StatusOK)},
	},
	Accepted: func(msg string) ServiceResp {
		return ServiceResp{http.StatusAccepted, ServiceErrMsg{msg}}
	},
	NoContent: ServiceResp{
		http.StatusNoContent, ServiceErrMsg{http.StatusText(http.StatusNoContent)},
	},
	Found: func(uri string) ServiceResp {
		return ServiceResp{http.StatusFound, ServiceErrMsg{uri}}
	},
	NotModified: func(msg string) ServiceResp {
		return ServiceResp{http.StatusNotModified, ServiceErrMsg{msg}}
	},
	BadRequestError: func(msg string) ServiceResp {
		return ServiceResp{http.StatusBadRequest, ServiceErrMsg{msg}}
	},
	ForbiddenError: func(msg string) ServiceResp {
		return ServiceResp{http.StatusForbidden, ServiceErrMsg{msg}}
	},
	NotFoundError: ServiceResp{
		http.StatusNotFound, ServiceErrMsg{http.StatusText(http.StatusNotFound)},
	},
	FailedDependencyError: func(msg string) ServiceResp {
		return ServiceResp{http.StatusFailedDependency, ServiceErrMsg{msg}}
	},
	InternalServiceError: func(msg string) ServiceResp {
		return ServiceResp{http.StatusInternalServerError, ServiceErrMsg{msg}}
	},
}
