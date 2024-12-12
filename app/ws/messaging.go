package ws

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repositories"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
	"go.elastic.co/apm"
)

func ServeWSMessaging(app *fiber.App) {
	clients := make(map[*websocket.Conn]bool)
	broadcast := make(chan models.MessagePayload)

	app.Get("/message/send", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()

		clients[c] = true

		for {
			var msg models.MessagePayload
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Println("error payload: ", err)
				break
			}

			tx := apm.DefaultTracer.StartTransaction("Send message", "ws")
			ctx := apm.ContextWithTransaction(context.Background(), tx)

			msg.Date = time.Now()
			if err := repositories.InsertNewMessage(ctx, msg); err != nil {
				log.Println(err)
				break
			}
			tx.End()
			broadcast <- msg
		}
	}))

	go func() {
		for {
			msg := <-broadcast
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Println("Failed to write json: ", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT_SOCKET", "8080"))))
}
