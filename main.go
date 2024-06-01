package main

import (
	"context"

	"github.com/Sahil-4555/mvc/configs/database"
	"github.com/Sahil-4555/mvc/routes"
	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/Sahil-4555/mvc/shared/utils"
)

// /*
// *
// @title Chat Application API
// @version 1.0
// @description Chat Application API
// @securityDefinitions.apikey BearerAuth
// @BasePath /
// @in header
// @name Authorization
// @description Bearer token authentication
// *
// */
func main() {

	log.Init()

	database.Init()
	log.GetLog().Info("", "DB connected")

	go routes.Run()

	utils.GracefulStop(log.GetLog(), func(ctx context.Context) error {
		var err error
		if err = routes.Close(ctx); err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
			return err
		}
		if err = database.Close(); err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
			return err
		}
		return nil
	})
}
