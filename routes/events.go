package routes

import (
	"fmt"
	"github.com/agrism/go-event-booking-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetEvents(context *gin.Context) {

	events, err := models.GetAllEvents()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events, Try again later!"})
		fmt.Println(err)
		return
	}

	context.JSON(http.StatusOK, events)
}

func GetEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Id param incorrect in path!"})
		fmt.Println(err)
		return
	}

	event, err := models.GetEventByID(id)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		fmt.Println(err)
		return
	}

	context.JSON(http.StatusOK, event)
}

func CreateEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data!"})
		fmt.Println(err)
		return
	}

	event.UserID = context.GetInt64("userId")

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save event!"})
		fmt.Println(err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}

func UpdateEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Id param incorrect in path!"})
		fmt.Println(err)
		return
	}

	event, err := models.GetEventByID(id)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		fmt.Println(err)
		return
	}

	if event.UserID != context.GetInt64("userId") {
		context.JSON(http.StatusForbidden, gin.H{"message": "Not authorized to update event"})
		return
	}

	err = context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data!"})
		fmt.Println(err)
		return
	}

	event.ID = id

	err = event.UpdateEvent()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Event update error"})
		fmt.Println(err)
		return
	}

	context.JSON(http.StatusOK, event)
}

func DeleteEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Id param incorrect in path!"})
		fmt.Println(err)
		return
	}

	event, err := models.GetEventByID(id)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		fmt.Println(err)
		return
	}

	if event.UserID != context.GetInt64("userId") {
		context.JSON(http.StatusForbidden, gin.H{"message": "Not authorized to delete event"})
		return
	}

	err = event.DeleteEvent()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error, please try later!"})
		fmt.Println(err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted!"})
}
