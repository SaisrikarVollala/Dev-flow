package main

import (
	"log"
	"os"

	"github.com/SaisrikarVollala/Dev-flow/internal/routes"
	appvalidator "github.com/SaisrikarVollala/Dev-flow/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err:=godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")


	appvalidator.Init()
	r := gin.Default()
	

api := r.Group("/api/v1")
routes.RegisterAuthRoute(api)



r.Run(":"+port)

}
