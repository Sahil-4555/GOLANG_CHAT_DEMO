package configs

import (
	"fmt"
	"os"

	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/joho/godotenv"
)

func MongoURI() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	Username := os.Getenv("DB_USER_MONGO")
	Password := os.Getenv("DB_PASSWORD_MONGO")
	Host := os.Getenv("DB_HOST_MONGO")
	Port := os.Getenv("DB_PORT_MONGO")
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", Username, Password, Host, Port)
	return connectionString
	// return os.Getenv("MONGO_URI")
}

func Database() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("DB_NAME_MONGO")
}

func Port() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("SERVER_PORT")
}

func AccessKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("AWS_ACCESS_KEY")
}

func SecretKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("AWS_SECRET_KEY")
}

func Region() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("AWS_REGION")
}

func Bucket() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("S3_BUCKET")
}

func JwtApiAuthKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("JWT_API_AUTH_KEY")
}
