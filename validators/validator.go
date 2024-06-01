package validators

import (
	"regexp"
	"strings"

	"github.com/Sahil-4555/mvc/models"
	"github.com/go-playground/validator"
)

func GetError(field, tag string) string {
	if tag == "required" {
		return field + " is " + tag
	} else if tag == "email" {
		return field + " not valid"
	}
	return ""
}

func InArrayString(val string, array []string) (exists bool) {
	exists = false
	for _, v := range array {
		if val == v {
			exists = true
			return
		}
	}
	return
}

func ValidateStruct(req interface{}, key string) (string, bool) {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		valErrs := err.(validator.ValidationErrors)
		for _, v := range valErrs {
			fieldName := strings.Replace(strings.Replace(v.Namespace(), key+".", "", 1), ".", " ", 3)
			reg, _ := regexp.Compile("[^A-Z`[]]+")
			fieldName = strings.Replace(reg.ReplaceAllString(fieldName, ""), "[", "", 2)
			errorString := GetError(fieldName, v.Tag())
			if errorString == "" {
				errorString = "Some of required field are missing or invalid"
			}
			return errorString, false
		}
	}
	return "", true
}

func ValidateUserStruct(req models.User) (string, bool) {
	errorString := ""
	if req.UserName == "" {
		errorString = "Username is missing or invalid"
		return errorString, false
	}
	if len(req.UserName) > 50 {
		errorString = "Username is max 50 character limited"
		return errorString, false
	}
	if req.Email != "" {
		if !IsEmailValid(req.Email) {
			errorString = "Please use valid email address for users"
			return errorString, false
		}
	} else {
		errorString = "User email is missing."
		return errorString, false
	}
	return "", true
}

func ValidateChannelStruct(req models.Channel) (string, bool) {
	errorString := ""
	if req.ChannelName == "" {
		errorString = "Channelname is missing or invalid"
		return errorString, false
	}
	if len(req.ChannelName) > 50 {
		errorString = "Channelname is max 50 character limited"
		return errorString, false
	}
	channelType := []string{"one-to-one", "private"}
	if !InArrayString(req.ChannelType, channelType) {
		errorString = "Channel type is missing or invalid"
		return errorString, false
	}
	return "", true
}

func ValidateMessageStruct(req models.Message) (string, bool) {
	errorString := ""
	contentType := []string{"text", "media", "both"}
	if !InArrayString(req.ContentType, contentType) {
		errorString = "Content type is missing or invalid"
		return errorString, false
	}
	return "", true
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
