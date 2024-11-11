package primary

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roxxers/surfe-techtest/internal/core/services"
)

var (
	ErrUserNotFound  = errors.New("User not found")
	ErrInvalidUserId = errors.New("Invalid user ID")
)

// Using one controller here but we would likely seperate them depending on how we structured our API layer.
type ActionRequest struct {
	ActionType string `json:"action"`
}

type FetchUserResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type UserActionCountResponse struct {
	Count int32 `json:"count"`
}

type Controller struct {
	service *services.Service
}

func NewController(service *services.Service) *Controller {
	return &Controller{service}
}

func (c *Controller) FetchUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// Would be better to used defined errors here
		ctx.JSON(400, gin.H{"error": ErrInvalidUserId.Error()})
		return
	}
	user, err := c.service.FetchUser(int64(userId))
	if err != nil {
		ctx.JSON(404, gin.H{"error": ErrUserNotFound.Error()})
		return
	}

	// Using a type here and explicitly naming wanted keys as to not leak whole user table
	response := FetchUserResponse{
		Id:        user.Id,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	ctx.JSON(200, response)
}

// FINISH DOCS
// ERROR HANDLING OUT OF RANGE EXECPTIONS!!!!!

func (c *Controller) GetUserActionCount(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// Would be better to used defined errors here
		ctx.JSON(400, gin.H{"error": ErrInvalidUserId.Error()})
		return
	}

	count, err := c.service.GetUserActionCount(int64(userId))
	if err != nil {
		ctx.JSON(404, gin.H{"error": ErrUserNotFound.Error()})
		return
	}

	response := UserActionCountResponse{
		Count: count,
	}
	ctx.JSON(200, response)
	return
}

func (c *Controller) CalculateNextActionProbablity(ctx *gin.Context) {
	var actionRequest ActionRequest
	ctx.BindJSON(&actionRequest)

	probabilities, err := c.service.CalculateNextActionProbablity(actionRequest.ActionType)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not find probablites: " + err.Error()})
		return
	}

	// Already formatted in the wanted way
	ctx.JSON(200, probabilities)
}

func (c *Controller) CalculateAllUserReferalIndexes(ctx *gin.Context) {
	indexes := c.service.CalculateAllUserReferalIndexes()
	ctx.JSON(200, indexes)
}
