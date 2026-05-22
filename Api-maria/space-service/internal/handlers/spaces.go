package handlers

import (
	"net/http"
	"strconv"

	"coworking/space-service/internal/services"

	"github.com/gin-gonic/gin"
)

// SpaceHandler groups HTTP handlers for the /spaces resource.
type SpaceHandler struct {
	svc *services.SpaceService
}

// NewSpaceHandler constructs a SpaceHandler.
func NewSpaceHandler(svc *services.SpaceService) *SpaceHandler {
	return &SpaceHandler{svc: svc}
}

// ListSpaces handles GET /spaces
// @Summary List spaces
// @Description Returns all spaces; ?active_only=true filters to active ones only.
// @Tags Spaces
// @Accept json
// @Produce json
// @Param active_only query bool false "Active only"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /spaces [get]
func (h *SpaceHandler) ListSpaces(c *gin.Context) {
	activeOnly := c.DefaultQuery("active_only", "true") == "true"
	spaces, err := h.svc.GetSpaces(activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve spaces: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    spaces,
		"message": "Spaces retrieved successfully",
	})
}

// CreateSpace handles POST /spaces (admin only)
// @Summary Create space
// @Description Create a new space (admin only)
// @Tags Spaces
// @Accept json
// @Produce json
// @Param space body services.CreateSpaceInput true "Space to create"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /spaces [post]
func (h *SpaceHandler) CreateSpace(c *gin.Context) {
	var input services.CreateSpaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	space, err := h.svc.CreateSpace(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    space,
		"message": "Space created successfully",
	})
}

// GetSpace handles GET /spaces/:id
// @Summary Get space
// @Description Get a space by ID
// @Tags Spaces
// @Accept json
// @Produce json
// @Param id path string true "Space ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /spaces/{id} [get]
func (h *SpaceHandler) GetSpace(c *gin.Context) {
	id := c.Param("id")
	space, err := h.svc.GetSpace(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "space not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    space,
		"message": "Space retrieved successfully",
	})
}

// UpdateSpace handles PUT /spaces/:id (admin only)
// @Summary Update space
// @Description Update an existing space by ID (admin only)
// @Tags Spaces
// @Accept json
// @Produce json
// @Param id path string true "Space ID"
// @Param space body services.UpdateSpaceInput true "Space update info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /spaces/{id} [put]
func (h *SpaceHandler) UpdateSpace(c *gin.Context) {
	id := c.Param("id")

	var input services.UpdateSpaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	space, err := h.svc.UpdateSpace(id, input)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "space not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    space,
		"message": "Space updated successfully",
	})
}

// DeleteSpace handles DELETE /spaces/:id (admin only)
// @Summary Delete space
// @Description Delete/Deactivate a space by ID (admin only)
// @Tags Spaces
// @Accept json
// @Produce json
// @Param id path string true "Space ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /spaces/{id} [delete]
func (h *SpaceHandler) DeleteSpace(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteSpace(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "space not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    nil,
		"message": "Space deactivated successfully",
	})
}

// GetAvailableSpaces handles GET /spaces/available
// @Summary Get available spaces
// @Description Get available spaces filtered by type and min_capacity
// @Tags Spaces
// @Accept json
// @Produce json
// @Param type query string false "Space type"
// @Param min_capacity query int false "Minimum capacity"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /spaces/available [get]
func (h *SpaceHandler) GetAvailableSpaces(c *gin.Context) {
	spaceType := c.Query("type")
	minCapacityStr := c.DefaultQuery("min_capacity", "0")

	minCapacity, err := strconv.Atoi(minCapacityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "min_capacity must be a valid integer",
		})
		return
	}

	spaces, err := h.svc.GetAvailableSpaces(spaceType, minCapacity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    spaces,
		"message": "Available spaces retrieved successfully",
	})
}
