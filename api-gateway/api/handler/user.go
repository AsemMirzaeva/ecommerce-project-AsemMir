package handler

import (
	"context"
	"io"
	"log"
	"net/http"

	usr "api-gateway/protos/user-service"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type UserHandler struct {
	client usr.UserServiceClient
}

func NewUserHandler(conn *grpc.ClientConn) *UserHandler {
	client := usr.NewUserServiceClient(conn)
	return &UserHandler{client: client}
}

// CreateUser ...
// @Summary CreateUser
// @Security ApiKeyAuth
// @Description Api for creating a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param User body models.CreateUser true "createUserModel"
// @Success 200 {object} models.User
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req usr.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.CreateUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUser gets user by id
// @Summary GetUser
// @Security ApiKeyAuth
// @Description Api for getting user by id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	req := usr.GetUserRequest{Id: id}
	resp, err := h.client.GetUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateUser updates user by id
// @Summary UpdateUser
// @Security ApiKeyAuth
// @Description Api for updating users by id
// @Tags Users
// @Accept json
// @Produce json
// @Param  id path string true "user_id"
// @Param User body models.UpdateUser true "updateUserModel"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req = usr.UpdateUserRequest{
		Id: id,
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.UpdateUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteUser deletes user by id
// @Summary DeleteUser
// @Security ApiKeyAuth
// @Description Api for deleting users by id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.User
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	req := usr.DeleteUserRequest{Id: id}
	resp, err := h.client.DeleteUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieve a list of all users
// @Tags Users
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	req := usr.ListUsersRequest{
		Limit: limit,
		Page:  page,
	}

	stream, err := h.client.ListUsers(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var users []usr.ListUsersResponse
	for {
		user, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				log.Printf("Error receiving user: %v", err)
			}
			break
		}
		users = append(users, *user)
	}

	// Log the number of users retrieved (optional)
	log.Printf("Retrieved %d users", len(users))

	c.JSON(http.StatusOK, users)
}
