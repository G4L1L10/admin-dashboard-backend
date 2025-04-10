package handler

import (
	"net/http"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/service"
	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionService *service.QuestionService
}

func NewQuestionHandler(questionService *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService}
}

// POST /api/questions
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req struct {
		LessonID     string   `json:"lesson_id" binding:"required"`
		QuestionText string   `json:"question_text" binding:"required"`
		QuestionType string   `json:"question_type" binding:"required"`
		ImageURL     *string  `json:"image_url"`
		AudioURL     *string  `json:"audio_url"`
		Answer       string   `json:"answer" binding:"required"`
		Explanation  string   `json:"explanation"`
		Options      []string `json:"options" binding:"required,dive,required"`
		Tags         []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}

	question := &model.Question{
		ID:           utils.GenerateUUID(),
		LessonID:     req.LessonID,
		QuestionText: req.QuestionText,
		QuestionType: req.QuestionType,
		ImageURL:     req.ImageURL,
		AudioURL:     req.AudioURL,
		Answer:       req.Answer,
		Explanation:  req.Explanation,
	}

	var options []*model.Option
	for _, optText := range req.Options {
		options = append(options, &model.Option{
			ID:         utils.GenerateUUID(),
			QuestionID: question.ID,
			OptionText: optText,
		})
	}

	if err := h.questionService.CreateQuestion(question, options, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create question", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": question.ID, "message": "Question created successfully"})
}

// GET /api/questions/:id
func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	id := c.Param("id")
	question, err := h.questionService.GetQuestionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.NewErrorResponse("Question not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) GetQuestionsByLesson(c *gin.Context) {
	lessonID := c.Param("lesson_id")

	questions, err := h.questionService.GetQuestionsByLessonID(lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to fetch questions", err.Error()))
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (h *QuestionHandler) GetQuestionsByTag(c *gin.Context) {
	tag := c.Query("tag")

	questions, err := h.questionService.GetQuestionsByTag(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to fetch questions by tag", err.Error()))
		return
	}

	c.JSON(http.StatusOK, questions)
}

// PUT /api/questions/:id
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	id := c.Param("id")

	var req model.Question
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}
	req.ID = id

	if err := h.questionService.UpdateQuestion(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update question", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question updated successfully"})
}

// DELETE /api/questions/:id
func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id := c.Param("id")

	if err := h.questionService.DeleteQuestion(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to delete question", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

