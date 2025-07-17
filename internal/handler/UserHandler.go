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

// GetAll godoc
// @Summary Ambil semua user
// @Description Mengambil daftar user dengan pagination
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Halaman saat ini"
// @Param limit query int false "Jumlah item per halaman"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/users [get]
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

// Create godoc
// @Summary Tambah user baru
// @Description Membuat user baru oleh super admin
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.UserCreateRequest true "Data user baru"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/users [post]
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

// FindByID godoc
// @Summary Ambil user berdasarkan ID
// @Description Mengambil informasi user berdasarkan ID
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID user"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) FindByID(c *gin.Context) {
	res := response.NewResponder(c)
	user, err := h.userService.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	res.OK(user, "query ok", nil)
}

// EmailUpdate godoc
// @Summary Perbarui email user
// @Description Memperbarui email user berdasarkan ID
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID user"
// @Param request body request.UserUpdateEmailRequest true "Email baru"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/users/{id}/email [patch]
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

// RoleUpdate godoc
// @Summary Perbarui role user
// @Description Menambahkan atau mengubah role user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID user"
// @Param request body request.RoleRequest true "Daftar role baru"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/users/{id}/roles-update [patch]
func (h *UserHandler) RoleUpdate(c *gin.Context) {
	res := response.NewResponder(c)
	var req request.RoleRequest
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

	// roles update
	rolesUpdate, err := h.userService.RolesUpdate(ctx, user, req.Roles)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	if !rolesUpdate {
		res.OK(nil, "tidak ada perubahan pada roles", nil)
		return
	}

	res.Created(nil, "role berhasil di update")
}

// Delete godoc
// @Summary Hapus user
// @Description Menghapus user berdasarkan ID
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID user"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/users/{id} [delete]
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
