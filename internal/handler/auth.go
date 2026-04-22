package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/project/kocokan/internal/service"
	"github.com/project/kocokan/pkg/response"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc} }

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	u, err := h.svc.Register(req.Name, req.Email, req.Password)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, u)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	token, user, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	c.SetCookie("koco_token", token, int(72*time.Hour.Seconds()), "/", "", false, true)
	response.OK(c, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("koco_token", "", -1, "/", "", false, true)
	response.Message(c, "Logout berhasil")
}

func (h *AuthHandler) Me(c *gin.Context) {
	// user ID already validated by middleware
	c.JSON(http.StatusOK, gin.H{"success": true})
}
