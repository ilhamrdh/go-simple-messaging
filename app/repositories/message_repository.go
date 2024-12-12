package repositories

import (
	"context"
	"fmt"

	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/pkg/database"
	"go.elastic.co/apm"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertNewMessage(ctx context.Context, data models.MessagePayload) error {
	span, _ := apm.StartSpan(ctx, "InsertNewMessage", "repository")
	defer span.End()
	_, err := database.MongoDB.InsertOne(ctx, data)
	return err
}

func GetAllMessage(ctx context.Context) ([]models.MessagePayload, error) {
	data := []models.MessagePayload{}

	cursor, err := database.MongoDB.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find message: %v", err)
	}

	for cursor.Next(ctx) {
		payload := models.MessagePayload{}
		if err := cursor.Decode(&payload); err != nil {
			return nil, err
		}
		data = append(data, payload)
	}
	return data, nil
}
