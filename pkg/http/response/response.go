package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

const (
	ERROR        = 500
	UNAUTHORIZED = 401
	BADREQUEST   = 400
	SUCCESS      = 200
)

func NewResponseWithData(code int, msg string, result interface{}) Response {
	return Response{
		code,
		msg,
		result,
	}
}

func NewResponse(code int, msg string) Response {
	return Response{
		code,
		msg,
		map[string]interface{}{},
	}
}

func Result(c *gin.Context, code int, msg string, result interface{}) {
	c.JSON(http.StatusOK, NewResponseWithData(code, msg, result))
}

func Ok(c *gin.Context, msg string, result interface{}) {
	Result(c, SUCCESS, msg, result)
}

func Fail(c *gin.Context, msg string, result interface{}) {
	Result(c, ERROR, msg, result)
}

func Unauthorized(c *gin.Context, msg string, result interface{}) {
	Result(c, UNAUTHORIZED, msg, result)
}

func BadRequest(c *gin.Context, msg string, result interface{}) {
	Result(c, BADREQUEST, msg, result)
}
