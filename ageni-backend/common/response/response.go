package response

import (
	"errors"
	"github.com/readonme/open-studio/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Resp .
type Resp struct {
	Code    int32       `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// JSONResponse .
func JSONResponse(c *gin.Context, httpStatus int, ecode int32, msg string, data interface{}) {
	resp := &Resp{
		Code:    ecode,
		Data:    data,
		Message: msg,
	}
	c.JSON(httpStatus, resp)
}

// JSONSuccess .
func JSONSuccess(c *gin.Context, data interface{}) {
	JSONResponse(c, http.StatusOK, 0, "success", data)
}

func JSONFail(c *gin.Context, err error, data interface{}) {
	if err == nil {
		return
	}

	var customErr *common.CustomError
	ok := errors.As(err, &customErr)
	if !ok {
		JSONResponse(c, http.StatusBadRequest, -1, err.Error(), data)
		return
	}

	JSONResponse(c, http.StatusOK, customErr.Code, customErr.Message, data)
}

// AbortWithJSONResponse .
func AbortWithJSONResponse(c *gin.Context, httpStatus int, ecode int32, msg string, data interface{}) {
	resp := &Resp{
		Code:    ecode,
		Data:    data,
		Message: msg,
	}
	c.Abort()
	c.JSON(httpStatus, resp)
}

// AbortWithJSONSuccess .
func AbortWithJSONSuccess(c *gin.Context, data interface{}) {
	AbortWithJSONResponse(c, http.StatusOK, 0, "success", data)
}

// AbortWithJSONFail .
func AbortWithJSONFail(c *gin.Context, ecode int32, msg string) {
	AbortWithJSONResponse(c, http.StatusOK, ecode, msg, nil)
}
