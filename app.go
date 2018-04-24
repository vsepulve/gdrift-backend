package main

import (
	"gdrift/models"
	"gdrift/utils"

	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadConfig("config/config.yaml")

//	db.Setup()

	app := gin.Default()
	//app.Use(cors.Default())
	app.Use(utils.CorsMiddleware())

	models.Setup(app)
//	routes.Setup(app)

	app.Run(":" + utils.Config.Server.Port)
}
