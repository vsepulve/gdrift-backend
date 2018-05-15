package models

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"github.com/vsepulve/gdrift-backend/db"
	"gdrift/utils"
//	"fmt"

)

type Users struct {
	Id               int           `gorm:"column:id;not null;primary_key"`
	Username         string        `gorm:"column:username;not null"`
	Password         string        `gorm:"column:password;not null"`
	Name             string        `gorm:"column:name;not null"`
}

func UsersCRUD(app *gin.Engine) {
	app.GET("/users/:id", UsersFetchOne)
	app.GET("/users/", UsersFetchAll)
	app.POST("/users/", UsersCreate)
	app.PUT("/users/:id", UsersUpdate)
	app.DELETE("/users/:id", UsersRemove)
}

func UsersFetchOne(c *gin.Context) {
	id := c.Param("id")

	db := db.Database()
	defer db.Close()

	var usuario Users
	if err := db.Find(&usuario, id).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
//		db.Model(&usuario).Related(&usuario.Usuario_tipo, "Usuario_tipos_id")
//		db.Model(&usuario).Related(&usuario.Institucion, "Institucion_id")
		c.JSON(http.StatusOK, usuario)
	}
}

func UsersFetchAll(c *gin.Context) {

	db := db.Database()
	defer db.Close()

	var usuarios []Users
//	if err := db.Where("usuario_tipos_id != 6").Order("institucion asc").Find(&usuarios).Error; err != nil {
	if err := db.Find(&usuarios).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
//		for i := range usuarios {
//			db.Model(&usuarios[i]).Related(&usuarios[i].Usuario_tipo, "Usuario_tipos_id")
//			db.Model(&usuarios[i]).Related(&usuarios[i].Institucion, "Institucion_id")
//		}
		c.JSON(http.StatusOK, usuarios)
	}
}

func UsersCreate(c *gin.Context) {
	var usuario Users
	e := c.BindJSON(&usuario)
	utils.Check(e)

	db := db.Database()
	defer db.Close()

	if err := db.Create(&usuario).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusCreated, usuario)
	}
}

func UsersUpdate(c *gin.Context) {
	var usuario Users
	id := c.Params.ByName("id")

	db := db.Database()
	defer db.Close()

	if err := db.Where("id = ?", id).First(&usuario).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		c.BindJSON(&usuario)

		db.Save(&usuario)
		c.JSON(200, usuario)
	}
}

func UsersRemove(c *gin.Context) {
	var usuario Users

	db := db.Database()
	defer db.Close()

	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&usuario).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		db.Delete(&usuario)
		c.JSON(http.StatusOK, usuario)
	}
}



