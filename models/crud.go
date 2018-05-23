package models

import (
	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	// Usuarios
	UsersCRUD(app)
	// Consultas SQL normales
	ProjectsCRUD(app)
	SimulationsCRUD(app)
	// Comandos
	ProjectCommandsCRUD(app)
	SimulationCommandsCRUD(app)
}



