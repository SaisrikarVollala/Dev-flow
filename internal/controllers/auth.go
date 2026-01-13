package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/SaisrikarVollala/Dev-flow/internal/auth"
	"github.com/SaisrikarVollala/Dev-flow/internal/db"
	"github.com/SaisrikarVollala/Dev-flow/internal/models"
	appvalidator "github.com/SaisrikarVollala/Dev-flow/internal/validator"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Name        string `json:"name" validate:"required,min=3"`
    Password        string `json:"password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}
func getUserCollection() *mongo.Collection {
	return db.Database.Collection("users")
}

func FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := getUserCollection().
		FindOne(ctx, bson.M{"email": email}).
		Decode(&user)
	return &user, err
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil 
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := appvalidator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := FindUserByEmail(ctx, req.Email)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	check := CheckPasswordHash(req.Password, user.Password)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
	})
    
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := appvalidator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Check if email exists
	_, err := FindUserByEmail(ctx, req.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
	}
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// 2. Hash password
	hashed, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to secure password"})
		return
	}

	// 3. Save user
	newUser := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashed,
		CreatedAt: time.Now().UTC(),
	}

	_, err = getUserCollection().InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered",
	})
}
