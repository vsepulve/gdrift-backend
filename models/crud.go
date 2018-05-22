package models

import (
	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	UsersCRUD(app)
	ProjectsCRUD(app)
	CommandsCRUD(app)
}



