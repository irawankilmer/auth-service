package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/valigo"
	"github.com/irawankilmer/auth-service/internal/dto/request"
	"github.com/irawankilmer/auth-service/internal/service"
	"github.com/irawankilmer/auth-service/pkg/response"
	"strconv"
)

type UserHandler struct {
	userService service.UserService
	validate    *valigo.Valigo
}

func NewUserHandler(us service.UserService, v *valigo.Valigo) *UserHandler {
	return &UserHandler{userService: us, validate: v}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	res := response.NewResponder(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	users, total, err := h.userService.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	meta := response.MetaData{
		Total: total,
		Page:  page,
		Limit: limit,
	}

	res.OK(users, "query ok", &meta)
}

func (h *UserHandler) Create(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.UserCreateRequest
	if !h.validate.ValigoJSON(c, &req) {
		return
	}

	if err := h.userService.Create(c.Request.Context(), req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.Created(nil, "user baru berhasil dibuat")
}

func (h *UserHandler) FindByID(c *gin.Context) {
	res := response.NewResponder(c)
	user, err := h.userService.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(user, "query ok", nil)
}

func (h *UserHandler) UsernameUpdate(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.UserUpdateUsernameRequest
	ctx := c.Request.Context()

	// cek user
	user, err := h.userService.FindByID(ctx, c.Param("id"))
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	// validasi
	if !h.validate.ValigoJSON(c, &req) {
		return
	}

	// update username
	update, err := h.userService.UsernameUpdate(ctx, user, req.Username)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	// jika tidak ada perubahan pada username
	if !update {
		res.OK(nil, "tidak ada perubahan username", nil)
		return
	}

	res.Created(nil, "username berhasil di update")
}

func (h *UserHandler) EmailUpdate(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.UserUpdateEmailRequest
	ctx := c.Request.Context()

	// cek user
	user, err := h.userService.FindByID(ctx, c.Param("id"))
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	// validasi
	if !h.validate.ValigoJSON(c, &req) {
		return
	}

	// update email
	emailUpdate, err := h.userService.EmailUpdate(ctx, user, req.Email)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}
	if !emailUpdate {
		res.OK(nil, "tidak ada perubahan email", nil)
		return
	}

	res.OK(nil, "email berhasil di update", nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	res := response.NewResponder(c)
	user, err := h.userService.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	if err := h.userService.Delete(c.Request.Context(), user); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(nil, "user berhasil dihapus", nil)
}
