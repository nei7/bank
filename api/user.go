package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nei7/bank/internal/db"
	"github.com/nei7/bank/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Password          string    `json:"password" `
	FullName          string    `json:"full_name" `
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func (server *Server) createUser(c *gin.Context) {

	var request createUserRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: request.Username,
		FullName: request.FullName,
		Email:    request.Email,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}

		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusCreated, newUserResponse(user))
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if err = util.CheckPassword(req.Password, user.Password); err != nil {
		c.JSON(http.StatusConflict, errResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	response := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	c.JSON(http.StatusOK, response)
}
