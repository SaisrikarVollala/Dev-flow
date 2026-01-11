package controller

import (
	"github.com/SaisrikarVollala/Dev-flow/internal/auth"
	appvalidator "github.com/SaisrikarVollala/Dev-flow/internal/validator"
	"github.com/gin-gonic/gin"
)


type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Username        string `json:"username" validate:"required,min=3"`
    Password        string `json:"password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

func Login(c *gin.Context) {
    var req LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }

    if err := appvalidator.Validate.Struct(req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // TODO:
    // 1. Find user by email
    // 2. Compare bcrypt password

    userID := "12345" // from DB

    token, err := auth.GenerateJWT(userID)
    if err != nil {
        c.JSON(500, gin.H{"error": "token generation failed"})
        return
    }

    c.JSON(200, gin.H{
        "access_token": token,
        "token_type":   "Bearer",
    })
}

func Register(c *gin.Context) {
    var req RegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }

    if err := appvalidator.Validate.Struct(req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // TODO:
    // 1. Check if email exists
    // 2. Hash password (bcrypt)
    // 3. Save user to DB

    c.JSON(201, gin.H{
        "message": "user registered",
    })
}

func Info(c *gin.Context){

}