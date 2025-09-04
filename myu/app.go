package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	e := echo.New()

	e.GET("/messages", GetHandler)
	e.POST("/messages", PostHander)
	e.PATCH("/messages/:id", PathcHandler)
	e.DELETE("/messages/:id", DeleteHandler)

	e.POST("/users", CreateUserHandler)
	e.GET("/users", GetUsersHandler)
	e.GET("/messages/user/:id", GetUserMessagesHandler)
	e.PATCH("/messages/:id/read", MarkAsReadHandler)

	e.Start(":8080")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

//CRUD - Creata, Read, Updete, Delete

type Message struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	FromID    int    `json:"from_id"`
	ToId      int    `json:"to_id"`
	FromName  string `json:"from_name"`
	ToName    string `json:"to_name"`
	TimeStamp string `json:"timestamp"`
	IsRead    bool   `json:"is_read"`
}

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var users = make(map[int]Person)
var messages = make(map[int]Message)
var nextId int = 1

func CreateUserHandler(c echo.Context) error {
	var person Person
	if err := c.Bind(&person); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not add the user",
		})
	}
	person.Id = nextId
	nextId++
	users[person.Id] = person

	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "User was successfully added " + person.Name,
	})
}
func GetUsersHandler(c echo.Context) error {
	var userSLice []Person
	for _, user := range users {
		userSLice = append(userSLice, user)
	}
	return c.JSON(http.StatusOK, &userSLice)
}

func GetUserMessagesHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not find the user",
		})
	}
	var msgSLice []Message
	for _, msg := range messages {
		if msg.FromID == id || msg.ToId == id {
			msgSLice = append(msgSLice, msg)
		}
	}
	return c.JSON(http.StatusOK, &msgSLice)
}

func GetHandler(c echo.Context) error {
	// if r.Method == http.MethodGet {
	// 	fmt.Fprint(w, "Counter = ", counter)
	// } else {
	// 	fmt.Fprint(w, "Это не Get запрос")
	// }
	var msgSLice []Message

	for _, msg := range messages {
		msgSLice = append(msgSLice, msg)
	}

	return c.JSON(http.StatusOK, &msgSLice)
}

func PostHander(c echo.Context) error {
	// if r.Method == http.MethodPost {
	// 	counter++
	// 	fmt.Fprint(w, "Counter увеличен на 1", counter)
	// } else {
	// 	fmt.Fprint(w, "Это не Post запрос")
	// }
	var message Message

	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not add the message",
		})
	}
	message.ID = nextId
	nextId++
	messages[message.ID] = message
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message was successfully added",
	})
}

func PathcHandler(c echo.Context) error {
	idPAram := c.Param("id")
	id, err := strconv.Atoi(idPAram)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Не смогли перевестив в число",
		})
	}

	var updatedMessage Message
	if err := c.Bind(&updatedMessage); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not update the message",
		})
	}

	// updated := false

	// for i, message := range messages {
	// 	if message.ID == id {
	// 		updatedMessage.ID = id
	// 		messages[i] = updatedMessage
	// 		updated = true
	// 		break
	// 	}

	// } для слайса

	if _, exists := messages[id]; !exists {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not find the message",
		})
	}

	updatedMessage.ID = id
	messages[id] = updatedMessage

	// if !updated {
	// 	return c.JSON(http.StatusBadRequest, Response{
	// 		Status:  "Error",
	// 		Message: "Could not find the message",
	// 	})
	// } // для слайса

	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message was successfully updated",
	})
}

func DeleteHandler(c echo.Context) error {
	idPAram := c.Param("id")
	id, err := strconv.Atoi(idPAram)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Не смогли перевестив в число",
		})
	}

	if _, exists := messages[id]; !exists {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not find the message",
		})
	}

	delete(messages, id)
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message was successfully deleted",
	})
}
