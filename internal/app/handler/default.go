package handler

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/util"
)

// HealthHandler is health checker API
// @Tags     Default
// @Success  200  {string}  string  "ok"
// @Router   /health [get]
func HealthHandler(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

// VersionHandler is version checker API
// @Tags     Default
// @Success  200  {string}  string  "0.4.12"
// @Router   /version [get]
func VersionHandler(c *gin.Context) {
	version := config.Env.Version
	c.String(http.StatusOK, version)
}

func result(c *gin.Context, data interface{}, err model.ServiceResp) {
	switch err.Status {
	case http.StatusOK:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusOK, util.StructToJsonString(err.ErrCode))
		c.JSON(http.StatusOK, data)

	case http.StatusAccepted:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusAccepted, util.StructToJsonString(err.ErrCode))
		c.JSON(http.StatusAccepted, err.ErrCode)

	case http.StatusNoContent:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusNoContent, util.StructToJsonString(err.ErrCode))
		c.JSON(http.StatusNoContent, nil)

	case http.StatusFound:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusFound, util.StructToJsonString(err.ErrCode))
		location := url.URL{Path: err.ErrCode.Code}
		c.Redirect(http.StatusFound, location.RequestURI())

	case http.StatusNotModified:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusNotModified, util.StructToJsonString(err.ErrCode))
		c.JSON(http.StatusNotModified, err.ErrCode)

	case http.StatusBadRequest:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusBadRequest, util.StructToJsonString(err.ErrCode))
		c.JSON(http.StatusBadRequest, err.ErrCode)

	case http.StatusForbidden:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusForbidden, util.StructToJsonString(err.ErrCode))
		c.JSON(http.StatusForbidden, err.ErrCode)

	case http.StatusNotFound:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusNotFound, util.StructToJsonString(err.ErrCode))
		c.JSON(http.StatusNotFound, err.ErrCode)

	case http.StatusFailedDependency:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusFailedDependency, util.StructToJsonString(err.ErrCode))
		c.JSON(http.StatusFailedDependency, err.ErrCode)

	default:
		logger.Info.Printf("status=%+v, resp=%+v\n", http.StatusInternalServerError, util.StructToJsonString(err.ErrCode))
		c.JSON(http.StatusInternalServerError, err.ErrCode)
	}
}
