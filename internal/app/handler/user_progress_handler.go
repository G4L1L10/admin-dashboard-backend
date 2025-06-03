// internal/app/handler/user_progress_handler.go
package handler

import (
	"net/http"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserProgressHandler struct {
	Service *service.UserProgressService
}

// ✅ Constructor
func NewUserProgressHandler(service *service.UserProgressService) *UserProgressHandler {
	return &UserProgressHandler{
		Service: service,
	}
}

// ✅ GET /api/user-progress/:user_id
func (h *UserProgressHandler) GetUserProgress(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	progress, err := h.Service.GetUserProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// ✅ POST /api/user-progress
func (h *UserProgressHandler) MarkLessonCompleted(c *gin.Context) {
	var input model.MarkProgressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := h.Service.MarkLessonCompleted(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark lesson as completed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "lesson marked as completed"})
}

