package routes

import (
	"fmt"
	"github.com/agrism/go-event-booking-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func RegisterForEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing event ID!"})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found!"})
		return
	}

	err = event.RegisterUserToEvent(context.GetInt64("userId"))

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error on registering to event"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Registered"})
}

func CancelRegistration(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing event ID!"})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found!"})
		return
	}

	err = event.CancelRegistrationFromEvent(context.GetInt64("userId"))

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusForbidden, gin.H{"message": "Cant cancel registration"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Registration canceled"})
}
