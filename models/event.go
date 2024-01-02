package models

import (
	"errors"
	"fmt"
	"github.com/agrism/go-event-booking-api/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

var events = []Event{}

func (event *Event) Save() error {
	query := `
INSERT INTO events (name, description, location, dateTime, user_id) 
VALUES (?,?,?,?,?)
`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.UserID)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	event.ID = id

	if err != nil {
		return err
	}

	return nil
}

func GetAllEvents() ([]Event, error) {

	query := "SELECT * FROM events"

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (Event, error) {
	query := "SELECT * FROM events WHERE id = ?"

	row := db.DB.QueryRow(query, id)

	var event Event

	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func (event *Event) UpdateEvent() error {

	query := `
UPDATE events SET name = ?, description = ?, location = ?, dateTime = ?, user_id = ? WHERE id = ?
`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.UserID, event.ID)

	if err != nil {
		return err
	}

	return nil
}

func (event *Event) DeleteEvent() error {
	query := "DELETE FROM events WHERE id = ?"

	smtp, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer smtp.Close()

	_, err = smtp.Exec(event.ID)

	return err
}

func (event Event) RegisterUserToEvent(userId int64) error {
	query := "INSERT OR IGNORE INTO user_events (event_id, user_id) VALUES (?, ?)"

	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.ID, userId)

	return err
}

func (event Event) CancelRegistrationFromEvent(userId int64) error {
	query := "SELECT id FROM user_events WHERE event_id = ? AND user_id = ?"

	row := db.DB.QueryRow(query, event.ID, userId)

	var registrationId int64

	err := row.Scan(&registrationId)

	if err != nil {
		fmt.Println(err)
		return errors.New("User not registered to event!")
	}

	query = "DELETE FROM user_events WHERE user_id = ? AND event_id = ?"
	statement, err := db.DB.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return errors.New("Cannot delete registration")
	}

	defer statement.Close()

	_, err = statement.Exec(userId, event.ID)

	if err != nil {
		fmt.Println(err)
		return errors.New("Cannot delete registration")
	}

	return nil
}
