package configs

import (
	"os"

	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/joho/godotenv"
)

func MongoURI() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	// Host := os.Getenv("MONGODB_HOST")
	// Username := os.Getenv("MONGODB_USER")
	// Password := os.Getenv("MONGODB_PASSWORD")
	// Port := os.Getenv("MONGODB_PORT")
	// connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", Username, Password, Host, Port)
	// return connectionString
	return os.Getenv("MONGO_URI")
}

func Database() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("MONGODB_DATABASE")
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

	return os.Getenv("JWT_SECRET")
}
