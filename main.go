package main

import (
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/labstack/echo/v4"
)

type DeviceRequest struct {
	Name               string   `json:"name"`
	RegistrationTokens []string `json:"registration_token"`
}

type SendNotificationRequest struct {
	DeviceRequest []DeviceRequest   `json:"devices"`
	Topics        []string          `json:"topics"`
	ImageURL      string            `json:"image_url"`
	Data          map[string]string `json:"data"`
}

func main() {
	e := echo.New()
	e.POST("/send-notification", sendNoti)
	e.Logger.Fatal(e.Start(":8080"))

}

func sendNoti(c echo.Context) error {
	var request SendNotificationRequest
	err := c.Bind(&request)
	if err != nil {
		fmt.Printf("invalid request: %v", err)
		return err
	}

	fmt.Printf("Request %v", request.ImageURL)

	ctx := c.Request().Context()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
		return err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		fmt.Printf("error create client: %v", err)
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    "test noti 01",
			Body:     "hello world",
			ImageURL: "https://img.freepik.com/free-photo/front-view-sad-girl-being-bullied_23-2149748403.jpg",
		},
		Data: request.Data,
	}

	if len(request.Topics) > 0 {
		message.Condition = getTopicCondition(request.Topics)
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		fmt.Printf("error create client: %v", err)
		return err
	}

	fmt.Printf("send message successfully: %s", response)

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "send notification successfully",
	})
}

func getTopicCondition(topics []string) string {
	var topicArr []string
	for _, topic := range topics {
		topicArr = append(topicArr, fmt.Sprintf("'%s' in topics", topic))
	}

	return strings.Join(topicArr, " || ")
}
