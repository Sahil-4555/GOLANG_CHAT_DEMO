package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"chat-demo-golang/configs/database"
	"chat-demo-golang/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserByUSername function will return the user from the database as per username
func GetUserByUsername(username string) bool {
	var userDetails models.User
	conn := database.NewConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := conn.UserCollection().FindOne(ctx, bson.M{
		"user_name":  username,
		"deleted_at": bson.M{"$eq": nil},
	}).Decode(&userDetails)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true
		}
	}
	return false
}

func GetUserByID(userId primitive.ObjectID) models.User {
	var userDetails models.User
	conn := database.NewConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = conn.UserCollection().FindOne(ctx, bson.M{
		"_id":        userId,
		"deleted_at": bson.M{"$eq": nil},
	}).Decode(&userDetails)
	fmt.Print("---")
	return userDetails
}

// IsUsernameAvailable function will check if the username is available or not
func IsUsernameAvailable(username string) bool {
	return GetUserByUsername(username)
}

// GenerateSlug function will generate the slug
func GenerateSlug(channelname string) string {
	var slug string
	channelname = strings.ToLower(channelname)
	for i := 0; i < len(channelname); i++ {
		if channelname[i] == ' ' {
			slug += "-"
		} else {
			slug += string(channelname[i])
		}
	}
	return slug
}

// GetSlug fucntion will check the slug is exist in database or not by returning bool
func GetSlug(slug string) bool {
	var channel models.Channel
	conn := database.NewConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := conn.ChannelCollection().FindOne(ctx, bson.M{
		"slug": slug,
	}).Decode(&channel)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true
		}
	}
	return false
}

// IsSlugAvailable function will check the slug is available or not
func IsSlugAvailable(slug string) bool {
	return GetSlug(slug)
}
