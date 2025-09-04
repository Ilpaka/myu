package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Message struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	FromID    int    `json:"from_id"`
	ToID      int    `json:"to_id"`
	FromName  string `json:"from_name"`
	ToName    string `json:"to_name"`
	Timestamp string `json:"timestamp"`
	IsRead    bool   `json:"is_read"`
}

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var (
	users      = make(map[int]Person)
	messages   = make(map[int]Message)
	nextUserID = 1
	nextMsgID  = 1
)

func main() {
	e := echo.New()

	// Middleware для production
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS - ВАЖНО для фронтенда
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // В production лучше указать конкретные домены
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":   "ok",
			"time":     time.Now().Format(time.RFC3339),
			"users":    len(users),
			"messages": len(messages),
		})
	})

	// API endpoints
	e.POST("/users", CreateUserHandler)
	e.GET("/users", GetUsersHandler)
	e.GET("/messages", GetHandler)
	e.POST("/messages", PostHandler)
	e.PATCH("/messages/:id", PatchHandler)
	e.DELETE("/messages/:id", DeleteHandler)
	e.GET("/messages/user/:id", GetUserMessagesHandler)
	e.PATCH("/messages/:id/read", MarkAsReadHandler)

	// Graceful shutdown
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Server error: ", err)
		}
	}()

	fmt.Println("🚀 Messenger API запущен на порту 8080")
	fmt.Println("📋 Health check: http://localhost:8080/health")

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("🔄 Завершение сервера...")
	e.Logger.Info("Server stopped")
}

// Handlers - те же самые
func CreateUserHandler(c echo.Context) error {
	var person Person
	if err := c.Bind(&person); err != nil || person.Name == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not add the user",
		})
	}
	person.ID = nextUserID
	nextUserID++
	users[person.ID] = person
	return c.JSON(http.StatusOK, person)
}

func GetUsersHandler(c echo.Context) error {
	res := make([]Person, 0, len(users))
	for _, u := range users {
		res = append(res, u)
	}
	return c.JSON(http.StatusOK, res)
}

func GetHandler(c echo.Context) error {
	res := make([]Message, 0, len(messages))
	for _, m := range messages {
		res = append(res, m)
	}
	return c.JSON(http.StatusOK, res)
}

func PostHandler(c echo.Context) error {
	var msg Message
	if err := c.Bind(&msg); err != nil || msg.Text == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not add the message",
		})
	}

	from, okFrom := users[msg.FromID]
	to, okTo := users[msg.ToID]
	if !okFrom || !okTo {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Sender or receiver not found",
		})
	}

	msg.ID = nextMsgID
	nextMsgID++
	msg.FromName = from.Name
	msg.ToName = to.Name
	msg.Timestamp = time.Now().Format("15:04:05")
	msg.IsRead = false

	messages[msg.ID] = msg
	return c.JSON(http.StatusOK, msg)
}

func GetUserMessagesHandler(c echo.Context) error {
	idParam := c.Param("id")
	uid, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid user id",
		})
	}
	res := make([]Message, 0)
	for _, m := range messages {
		if m.ToID == uid {
			res = append(res, m)
		}
	}
	return c.JSON(http.StatusOK, res)
}

func MarkAsReadHandler(c echo.Context) error {
	idParam := c.Param("id")
	msgID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid message ID",
		})
	}
	msg, ok := messages[msgID]
	if !ok {
		return c.JSON(http.StatusNotFound, Response{
			Status:  "Error",
			Message: "Message not found",
		})
	}
	msg.IsRead = true
	messages[msgID] = msg
	return c.JSON(http.StatusOK, msg)
}

// Остальные handlers (PatchHandler, DeleteHandler) - аналогично
func PatchHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid id",
		})
	}

	var updated Message
	if err := c.Bind(&updated); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Could not update the message",
		})
	}
	existing, ok := messages[id]
	if !ok {
		return c.JSON(http.StatusNotFound, Response{
			Status:  "Error",
			Message: "Message not found",
		})
	}

	// Обновляем только разрешённые поля
	if updated.Text != "" {
		existing.Text = updated.Text
	}
	existing.IsRead = updated.IsRead

	messages[id] = existing
	return c.JSON(http.StatusOK, existing)
}

func DeleteHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid id",
		})
	}
	if _, ok := messages[id]; !ok {
		return c.JSON(http.StatusNotFound, Response{
			Status:  "Error",
			Message: "Message not found",
		})
	}
	delete(messages, id)
	return c.NoContent(http.StatusNoContent)
}
