package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project/kocokan/internal/service"
	"github.com/project/kocokan/pkg/middleware"
	"github.com/project/kocokan/pkg/response"
)

type GroupHandler struct{ svc *service.GroupService }

func NewGroupHandler(svc *service.GroupService) *GroupHandler { return &GroupHandler{svc} }

func (h *GroupHandler) List(c *gin.Context) {
	groups, err := h.svc.List(middleware.UserID(c))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.OK(c, groups)
}

func (h *GroupHandler) Get(c *gin.Context) {
	id := parseID(c, "id")
	g, err := h.svc.Get(id, middleware.UserID(c))
	if err != nil {
		response.Error(c, 404, "Grup tidak ditemukan")
		return
	}
	response.OK(c, g)
}

func (h *GroupHandler) Create(c *gin.Context) {
	var req struct {
		Name            string `json:"name" binding:"required"`
		Emoji           string `json:"emoji"`
		Description     string `json:"description"`
		PeriodType      string `json:"period_type"`
		NumParticipants int    `json:"num_participants" binding:"required,min=2"`
		TotalRounds     int    `json:"total_rounds" binding:"required,min=1"`
		PrizeAmount     int64  `json:"prize_amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if req.PeriodType == "" {
		req.PeriodType = "monthly"
	}
	g, err := h.svc.Create(middleware.UserID(c), req.Name, req.Emoji, req.Description, req.PeriodType, req.NumParticipants, req.TotalRounds, req.PrizeAmount)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.OK(c, g)
}

func (h *GroupHandler) Update(c *gin.Context) {
	id := parseID(c, "id")
	var req struct {
		Name            string `json:"name" binding:"required"`
		Emoji           string `json:"emoji"`
		Description     string `json:"description"`
		PeriodType      string `json:"period_type"`
		NumParticipants int    `json:"num_participants" binding:"required,min=2"`
		TotalRounds     int    `json:"total_rounds" binding:"required,min=1"`
		PrizeAmount     int64  `json:"prize_amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if req.PeriodType == "" {
		req.PeriodType = "monthly"
	}
	g, err := h.svc.Update(id, middleware.UserID(c), req.Name, req.Emoji, req.Description, req.PeriodType, req.NumParticipants, req.TotalRounds, req.PrizeAmount)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.OK(c, g)
}

func (h *GroupHandler) Delete(c *gin.Context) {
	id := parseID(c, "id")
	if err := h.svc.Delete(id, middleware.UserID(c)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Message(c, "Grup dihapus")
}

func (h *GroupHandler) AddParticipant(c *gin.Context) {
	groupID := parseID(c, "id")
	var req struct {
		Name  string `json:"name" binding:"required"`
		Phone string `json:"phone"`
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	p, err := h.svc.AddParticipant(groupID, middleware.UserID(c), req.Name, req.Phone, req.Notes)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, p)
}

func (h *GroupHandler) UpdateParticipant(c *gin.Context) {
	groupID := parseID(c, "id")
	pid := parseID(c, "pid")
	var req struct {
		Name  string `json:"name" binding:"required"`
		Phone string `json:"phone"`
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	p, err := h.svc.UpdateParticipant(pid, groupID, middleware.UserID(c), req.Name, req.Phone, req.Notes)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, p)
}

func (h *GroupHandler) DeleteParticipant(c *gin.Context) {
	groupID := parseID(c, "id")
	pid := parseID(c, "pid")
	if err := h.svc.DeleteParticipant(pid, groupID, middleware.UserID(c)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Message(c, "Peserta dihapus")
}

func (h *GroupHandler) Draw(c *gin.Context) {
	id := parseID(c, "id")
	round, err := h.svc.Draw(id, middleware.UserID(c))
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, round)
}

func (h *GroupHandler) UpdateWinner(c *gin.Context) {
	groupID := parseID(c, "id")
	roundID := parseID(c, "rid")
	var req struct {
		WinnerID uint   `json:"winner_id" binding:"required"`
		Notes    string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	round, err := h.svc.UpdateWinner(roundID, groupID, middleware.UserID(c), req.WinnerID, req.Notes)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, round)
}

func parseID(c *gin.Context, key string) uint {
	id, _ := strconv.ParseUint(c.Param(key), 10, 64)
	return uint(id)
}
