package handler

import (
	"net/http"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/service"
	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type OptionHandler struct {
	optionService *service.OptionService
}

func NewOptionHandler(optionService *service.OptionService) *OptionHandler {
	return &OptionHandler{optionService: optionService}
}

// POST /api/options
func (h *OptionHandler) CreateOption(c *gin.Context) {
	var req model.Option
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}
	req.ID = utils.GenerateUUID()

	if err := h.optionService.CreateOption(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create option", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": req.ID, "message": "Option created successfully"})
}

// GET /api/options/:id
func (h *OptionHandler) GetOption(c *gin.Context) {
	id := c.Param("id")
	option, err := h.optionService.GetOptionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.NewErrorResponse("Option not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, option)
}

// PUT /api/options/:id
func (h *OptionHandler) UpdateOption(c *gin.Context) {
	id := c.Param("id")

	var req model.Option
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}
	req.ID = id

	if err := h.optionService.UpdateOption(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update option", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Option updated successfully"})
}

// DELETE /api/options/:id
func (h *OptionHandler) DeleteOption(c *gin.Context) {
	id := c.Param("id")

	if err := h.optionService.DeleteOption(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to delete option", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Option deleted successfully"})
}

