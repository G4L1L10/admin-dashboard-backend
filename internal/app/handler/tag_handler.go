package handler

import (
	"net/http"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/service"
	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagService *service.TagService
}

func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// POST /api/tags
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req model.Tag
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}
	req.ID = utils.GenerateUUID()

	if err := h.tagService.CreateTag(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to create tag", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": req.ID, "message": "Tag created successfully"})
}

// GET /api/tags/:id
func (h *TagHandler) GetTag(c *gin.Context) {
	id := c.Param("id")
	tag, err := h.tagService.GetTagByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.NewErrorResponse("Tag not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) SearchTags(c *gin.Context) {
	keyword := c.Query("search")

	tags, err := h.tagService.SearchTags(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to search tags", err.Error()))
		return
	}

	c.JSON(http.StatusOK, tags)
}

// PUT /api/tags/:id
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id := c.Param("id")

	var req model.Tag
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request", err.Error()))
		return
	}
	req.ID = id

	if err := h.tagService.UpdateTag(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to update tag", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag updated successfully"})
}

// DELETE /api/tags/:id
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id := c.Param("id")

	if err := h.tagService.DeleteTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Failed to delete tag", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

