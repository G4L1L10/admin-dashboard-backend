package handler

import (
	"net/http"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/service"
	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService *service.CourseService
}

func NewCourseHandler(courseService *service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// POST /api/courses
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req model.Course
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}
	req.ID = utils.GenerateUUID()

	if err := h.courseService.CreateCourse(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create course", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": req.ID, "message": "Course created successfully"})
}

// GET /api/courses/:id
func (h *CourseHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")
	course, err := h.courseService.GetCourseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.NewErrorResponse("Course not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, course)
}

// GET /api/courses  ✅ NEW
func (h *CourseHandler) ListCourses(c *gin.Context) {
	courses, err := h.courseService.ListCourses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to fetch courses", err.Error()))
		return
	}

	c.JSON(http.StatusOK, courses)
}

// PUT /api/courses/:id
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id := c.Param("id")

	var req model.Course
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}
	req.ID = id

	if err := h.courseService.UpdateCourse(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update course", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

// DELETE /api/courses/:id
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	if err := h.courseService.DeleteCourse(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to delete course", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
