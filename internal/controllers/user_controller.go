package controllers

import (
	"net/http"

	"github.com/Koshsky/subs-service/internal/models"
	"github.com/Koshsky/subs-service/internal/services"
	"github.com/Koshsky/subs-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type UserController struct {
	UserService     *services.UserService
	JWTTokenManager *utils.JWTTokenManager
}

func NewUserController(service *services.UserService, jwtManager *utils.JWTTokenManager) *UserController {
	return &UserController{UserService: service, JWTTokenManager: jwtManager}
}

// Register handles user registration
func (uc *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := validate.Struct(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	if err := uc.UserService.RegisterUser(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// Login handles user authentication
func (uc *UserController) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials format"})
		return
	}

	user, err := uc.UserService.ValidateCredentials(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := uc.JWTTokenManager.GenerateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("auth_token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfull login",
	})
}
