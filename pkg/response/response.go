package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Responder interface {
	OK(data interface{}, message string, meta *MetaData)
	Created(data interface{}, message string)
	NoContent()
	BadRequest(errors interface{}, message string)
	Unauthorized(message string)
	Forbidden(message string)
	NotFound(message string)
	ServerError(message string)
}

type responder struct {
	c *gin.Context
}

func NewResponder(c *gin.Context) Responder {
	return &responder{c: c}
}

func (r *responder) OK(data interface{}, message string, meta *MetaData) {
	response := APIResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	r.c.JSON(http.StatusOK, response)
}

func (r *responder) Created(data interface{}, message string) {
	r.c.JSON(http.StatusCreated, APIResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func (r *responder) NoContent() {
	r.c.Status(http.StatusNoContent)
}

func (r *responder) BadRequest(errors interface{}, message string) {
	r.c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
		Code:    http.StatusBadRequest,
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}

func (r *responder) Unauthorized(message string) {
	r.c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
		Code:    http.StatusUnauthorized,
		Status:  "error",
		Message: message,
	})
}

func (r *responder) Forbidden(message string) {
	r.c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
		Code:    http.StatusForbidden,
		Status:  "error",
		Message: message,
	})
}

func (r *responder) NotFound(message string) {
	r.c.AbortWithStatusJSON(http.StatusNotFound, APIResponse{
		Code:    http.StatusNotFound,
		Status:  "error",
		Message: message,
	})
}

func (r *responder) ServerError(message string) {
	r.c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
		Code:    http.StatusInternalServerError,
		Status:  "error",
		Message: message,
	})
}
