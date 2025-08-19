package bootstrap

import (
	"chat-demo-golang/configs/database"
	"chat-demo-golang/models"
	"chat-demo-golang/shared/common"
	"chat-demo-golang/shared/log"
	"context"
	"encoding/json"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type DummyUser struct {
	UserName string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoadDummyUsers() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn := database.NewConnection()

	data, err := os.ReadFile("configs/bootstrap/dummy_user.json")
	if err != nil {
		log.GetLog().Error("Failed to read dummy users file: ", err.Error())
		return
	}

	var users []DummyUser
	if err := json.Unmarshal(data, &users); err != nil {
		log.GetLog().Error("Failed to parse dummy users JSON: ", err.Error())
		return
	}

	for _, u := range users {
		// check if user already exists
		var existing models.User
		err := conn.UserCollection().FindOne(ctx, bson.M{"email": u.Email}).Decode(&existing)
		if err == nil {
			log.GetLog().Info("Skipping existing dummy user: ", u.Email)
			continue
		}

		// hash password
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.GetLog().Error("Failed to hash password for dummy user: ", err.Error())
			continue
		}

		user := models.User{
			UserName: u.UserName,
			Name:     u.Name,
			Email:    u.Email,
			Password: string(hashPassword),
			Status:   common.USER_STATUS_OFFLINE,
		}
		user.TimeStamp()
		user.NewUser()

		result, err := conn.UserCollection().InsertOne(ctx, user)
		if err != nil {
			log.GetLog().Error("Failed to insert dummy user: ", err.Error())
			continue
		}

		id := result.InsertedID.(primitive.ObjectID)
		log.GetLog().Info("Inserted dummy user: ", u.Email, " with ID: ", id.Hex())
	}
}
