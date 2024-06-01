package middleware

import (
	"context"
	"fmt"

	"strings"
	"time"

	"github.com/Sahil-4555/mvc/configs"
	"github.com/Sahil-4555/mvc/configs/database"
	"github.com/Sahil-4555/mvc/models"
	"github.com/Sahil-4555/mvc/shared/common"
	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/Sahil-4555/mvc/shared/message"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var JWT_API_KEY = configs.JwtApiAuthKey()

type UserTokenData struct {
	Id primitive.ObjectID `json:"id" bson:"id"`
}

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		if !strings.HasPrefix(bearerToken, "Bearer ") {
			log.GetLog().Info("ERROR : ", "Authorized Request Invalid.")
			common.Respond(c, common.STATUS_UNAUTHORIZED, map[string]interface{}{
				"message": message.AuthorizationRequestInvalid,
				"code":    common.META_FAILED,
			})
			c.Abort()
			return
		}

		token := strings.Split(bearerToken, "Bearer ")
		if len(token) < 2 {
			log.GetLog().Info("ERROR : ", "Authorized Token Not Supplied.")
			common.Respond(c, common.STATUS_UNAUTHORIZED, map[string]interface{}{
				"message": message.AuthorizationTokenNotSupplied,
				"code":    common.META_FAILED,
			})
			c.Abort()
			return
		}

		valid, err := ValidateToken(token[1], JWT_API_KEY)
		if err != nil {
			log.GetLog().Info("ERROR : ", "Authorized Token Invalid.")
			common.Respond(c, common.STATUS_UNAUTHORIZED, map[string]interface{}{
				"message": message.AuthorizationTokenInvalid,
				"code":    common.META_FAILED,
			})
			c.Abort()
			return
		}

		userInfo := valid.Claims.(jwt.MapClaims)["userData"]
		data := userInfo.(map[string]interface{})

		userId, err := primitive.ObjectIDFromHex(data["ID"].(string))
		if err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
		}

		c.Set("userId", userId)

		var expTime int
		expInfo := valid.Claims.(jwt.MapClaims)["exp"]
		if expInfo != nil {
			val := expInfo.(float64)
			expTime = int(val)
		}
		c.Set("exp", expTime)

		c.Next()
	}
}

func AuthHandlerWebsocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Query("token")

		if !strings.HasPrefix(bearerToken, "Bearer ") {
			log.GetLog().Info("ERROR : ", "Authorized Request Invalid.")
			common.Respond(c, common.STATUS_UNAUTHORIZED, map[string]interface{}{
				"message": message.AuthorizationRequestInvalid,
				"code":    common.META_FAILED,
			})
			c.Abort()
			return
		}

		token := strings.Split(bearerToken, "Bearer ")
		if len(token) < 2 {
			log.GetLog().Info("ERROR : ", "Authorized Token Not Supplied.")
			common.Respond(c, common.STATUS_UNAUTHORIZED, map[string]interface{}{
				"message": message.AuthorizationTokenNotSupplied,
				"code":    common.META_FAILED,
			})
			c.Abort()
			return
		}

		valid, err := ValidateToken(token[1], JWT_API_KEY)
		if err != nil {
			log.GetLog().Info("ERROR : ", "Authorized Token Invalid.")
			common.Respond(c, common.STATUS_UNAUTHORIZED, map[string]interface{}{
				"message": message.AuthorizationTokenInvalid,
				"code":    common.META_FAILED,
			})
			c.Abort()
			return
		}

		userInfo := valid.Claims.(jwt.MapClaims)["userData"]
		data := userInfo.(map[string]interface{})

		userId, err := primitive.ObjectIDFromHex(data["ID"].(string))
		if err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
		}

		c.Set("userId", userId)

		var expTime int
		expInfo := valid.Claims.(jwt.MapClaims)["exp"]
		if expInfo != nil {
			val := expInfo.(float64)
			expTime = int(val)
		}
		c.Set("exp", expTime)

		c.Next()
	}
}

func GenerateToken(userData interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := make(jwt.MapClaims)
	claims["userData"] = userData
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(JWT_API_KEY))
	return tokenString, err
}

func ValidateToken(t string, k string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
		}

		return []byte(k), nil
	})
	if err != nil {
		return nil, err
	}

	return token, err
}

func GetUserData(c *gin.Context) (id primitive.ObjectID, err error) {
	userData, userExists := c.Get("userId")

	if !userExists {
		return primitive.NilObjectID, fmt.Errorf("user not exist")
	}

	data := userData.(primitive.ObjectID)

	return data, nil
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func GetUserById(userId string) models.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, _ := primitive.ObjectIDFromHex(userId)
	conn := database.NewConnection()

	var user models.User
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}
	conn.UserCollection().FindOne(ctx, filter).Decode(&user)

	return user
}

func GetChannelById(channelId primitive.ObjectID) models.Channel {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	var channel models.Channel
	filter := bson.M{"_id": channelId, "deleted_at": bson.M{"$eq": nil}}
	conn.ChannelCollection().FindOne(ctx, filter).Decode(&channel)

	return channel
}
