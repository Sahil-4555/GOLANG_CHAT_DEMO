package service

import (
	"context"
	"regexp"

	"time"

	"github.com/Sahil-4555/mvc/configs/crypto"
	"github.com/Sahil-4555/mvc/configs/database"

	"github.com/Sahil-4555/mvc/models"
	"github.com/Sahil-4555/mvc/shared/common"
	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/Sahil-4555/mvc/shared/message"
	"github.com/Sahil-4555/mvc/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var ApiAuthKey = "$(*%S$FDd!3)96|12AP&LR"

// call service
func SignUp(req common.SignUpReq) map[string]interface{} {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := database.NewConnection()

	var IsValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z.-]+\.[a-zA-Z]{2,4}$`).MatchString
	if !IsValidEmail(req.Email) {
		log.GetLog().Info("ERROR : ", "Invalid email format")
		return map[string]interface{}{
			"message": message.EmailInvalid,
			"code":    common.META_FAILED,
		}
	}

	// Check if the user with same email is not there
	var user models.User
	err := conn.UserCollection().FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err == nil {
		log.GetLog().Info("ERROR : ", "Email in use.")
		return map[string]interface{}{
			"message": message.EmailInUse,
			"code":    common.META_FAILED,
		}
	}

	// Check if the username is only alphanumberic
	var IsAlphaNumeric = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$`).MatchString
	if !IsAlphaNumeric(req.UserName) {
		log.GetLog().Info("ERROR : ", "Username not valid.")
		return map[string]interface{}{
			"message": message.UserNameNotValid,
			"code":    common.META_FAILED,
		}
	}

	// Check if the username is available or not
	if ok := utils.IsUsernameAvailable(req.UserName); !ok {
		log.GetLog().Info("ERROR : ", "Username is already in use.")
		return map[string]interface{}{
			"message": message.UsernameIsNotAvailable,
			"code":    common.META_FAILED,
		}
	}

	// Converts the password to hash
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		return map[string]interface{}{
			"message":  message.FailedToHashPassword,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST}
	}

	// Insert the user details in DB
	user = models.User{UserName: req.UserName, Name: req.Name, Email: req.Email, Password: string(hashPassword), Status: common.USER_STATUS_ONLINE}
	user.TimeStamp()
	user.NewUser()
	result, err := conn.UserCollection().InsertOne(ctx, user)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message":  message.FailedToInsert,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST}
	}

	newId := result.InsertedID.(primitive.ObjectID)
	tokenData := crypto.UserTokenData{
		ID: newId.Hex(),
	}

	token := crypto.GenerateAuthToken(tokenData)
	loginData := common.LoginResponse{
		ID:        newId,
		UserName:  user.UserName,
		Name:      user.Name,
		Email:     user.Email,
		LastLogin: time.Now(),
	}
	data := map[string]interface{}{
		"token": token,
		"data":  loginData,
	}

	response := common.ResponseSuccessWithToken(message.SignUpSuccess, common.META_SUCCESS, data)

	return response
}

func SignIn(req common.SignInReq) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var IsValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z.-]+\.[a-zA-Z]{2,4}$`).MatchString
	if !IsValidEmail(req.Email) {
		log.GetLog().Info("ERROR : ", "Invalid email format.")
		return map[string]interface{}{
			"message": message.EmailInvalid,
			"code":    common.META_FAILED,
		}
	}

	conn := database.NewConnection()
	var user models.User
	err := conn.UserCollection().FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		log.GetLog().Info("ERROR(Query) : ", err.Error())
		return map[string]interface{}{
			"message": message.EmailOrPasswordNotMatched,
			"code":    common.META_FAILED,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		return map[string]interface{}{
			"message": message.EmailOrPasswordNotMatched,
			"code":    common.META_FAILED,
		}
	}

	tokenData := crypto.UserTokenData{
		ID: user.Id.Hex(),
	}

	token := crypto.GenerateAuthToken(tokenData)
	loginData := common.LoginResponse{
		ID:        user.Id,
		UserName:  user.UserName,
		Name:      user.Name,
		Email:     user.Email,
		LastLogin: time.Now(),
	}

	data := map[string]interface{}{
		"token": token,
		"data":  loginData,
	}

	response := common.ResponseSuccessWithToken(message.LoginSuccess, common.META_SUCCESS, data)

	return response
}
