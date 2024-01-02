package routes

import (
	"fmt"
	"github.com/agrism/go-event-booking-api/models"
	"github.com/agrism/go-event-booking-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data!"})
		fmt.Println(err)
		return
	}

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "User saving error!"})
		fmt.Println(err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created!"})

}

func Login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data!"})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials provided!"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating jwt token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login success!", "token": token})

}
