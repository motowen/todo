package model

import (
	"net/http"
)

type ServiceResp struct {
	Status  int            `json:"status"`
	ErrCode ServiceErrCode `json:"errCode"`
}

type ServiceErrCode struct {
	Code string `json:"code"`
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
		http.StatusOK, ServiceErrCode{http.StatusText(http.StatusOK)},
	},
	Accepted: func(code string) ServiceResp {
		return ServiceResp{http.StatusAccepted, ServiceErrCode{code}}
	},
	NoContent: ServiceResp{
		http.StatusNoContent, ServiceErrCode{http.StatusText(http.StatusNoContent)},
	},
	Found: func(uri string) ServiceResp {
		return ServiceResp{http.StatusFound, ServiceErrCode{uri}}
	},
	NotModified: func(code string) ServiceResp {
		return ServiceResp{http.StatusNotModified, ServiceErrCode{code}}
	},
	BadRequestError: func(code string) ServiceResp {
		return ServiceResp{http.StatusBadRequest, ServiceErrCode{code}}
	},
	ForbiddenError: func(code string) ServiceResp {
		return ServiceResp{http.StatusForbidden, ServiceErrCode{code}}
	},
	NotFoundError: ServiceResp{
		http.StatusNotFound, ServiceErrCode{http.StatusText(http.StatusNotFound)},
	},
	FailedDependencyError: func(code string) ServiceResp {
		return ServiceResp{http.StatusFailedDependency, ServiceErrCode{code}}
	},
	InternalServiceError: func(code string) ServiceResp {
		return ServiceResp{http.StatusInternalServerError, ServiceErrCode{code}}
	},
}

// DB
const DBCreateTodoFail = "1001"
const DBFindTodoFail = "1002"
const DBUpdateTodoFail = "1003"
const DBDeleteTodoFail = "1004"
const DBTimeoutFail = "1005"
const DBGetIconPresignedURLFail = "1006"

// External
const ExternalGetAuthTokenFail = "2001"
const ExternalGetAuthTokenParseFail = "2002"
const ExternalCreateVendorFail = "2011"
const ExternalCreateVendorParseFail = "2012"
const ExternalGetVendorFail = "2013"
const ExternalUpdateVendorFail = "2014"
const ExternalDeleteVendorFail = "2015"

// HTTP
const HttpMethodInvalid = "3001"
