package handler

import (
	"net/http"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/service"
	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	lessonService *service.LessonService
}

func NewLessonHandler(lessonService *service.LessonService) *LessonHandler {
	return &LessonHandler{lessonService: lessonService}
}

// POST /api/lessons
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	var req model.Lesson
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}
	req.ID = utils.GenerateUUID()

	if err := h.lessonService.CreateLesson(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create lesson", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": req.ID, "message": "Lesson created successfully"})
}

// GET /api/lessons/:id
func (h *LessonHandler) GetLesson(c *gin.Context) {
	id := c.Param("id")
	lesson, err := h.lessonService.GetLessonByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.NewErrorResponse("Lesson not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func (h *LessonHandler) GetFullLesson(c *gin.Context) {
	lessonID := c.Param("id")

	fullLesson, err := h.lessonService.GetFullLesson(lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to fetch full lesson", err.Error()))
		return
	}

	c.JSON(http.StatusOK, fullLesson)
}

func (h *LessonHandler) GetLessonsByCourseID(c *gin.Context) {
	courseID := c.Param("course_id")

	lessons, err := h.lessonService.GetLessonsByCourseID(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to fetch lessons", err.Error()))
		return
	}

	c.JSON(http.StatusOK, lessons)
}

// PUT /api/lessons/:id
func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	id := c.Param("id")

	var req model.Lesson
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}
	req.ID = id

	if err := h.lessonService.UpdateLesson(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update lesson", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lesson updated successfully"})
}

// DELETE /api/lessons/:id
func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	id := c.Param("id")

	if err := h.lessonService.DeleteLesson(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to delete lesson", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lesson deleted successfully"})
}
