package main

import (
	"github.com/vsepulve/gdrift-backend/db"
	"github.com/vsepulve/gdrift-backend/models"
	"github.com/vsepulve/gdrift-backend/utils"

	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadConfig("config/config.yaml")

	db.Setup()

	app := gin.Default()
	//app.Use(cors.Default())
	app.Use(utils.CorsMiddleware())

	models.Setup(app)
//	routes.Setup(app)

	app.Run(":" + utils.Config.Server.Port)
}
