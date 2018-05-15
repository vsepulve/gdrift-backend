package models

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"time"

	"github.com/vsepulve/gdrift-backend/db"
	"gdrift/utils"

)

type Projects struct {
	Id               int           `gorm:"column:id;not null;primary_key"`
	Name             string        `gorm:"column:name;not null"`
	Owner_id         int           `gorm:"column:users_id;not null"`
	Owner            Users         `gorm:"ForeignKey:Owner_id;AssociationForeignKey:Id"`
	Date_created     *time.Time    `gorm:"column:date_created"`
	Date_finished    *time.Time    `gorm:"column:date_finished"`
}

func ProjectsCRUD(app *gin.Engine) {
	app.GET("/projects/:id", ProjectsFetchOne)
	app.GET("/projects/", ProjectsFetchAll)
	app.POST("/projects/", ProjectsCreate)
	app.PUT("/projects/:id", ProjectsUpdate)
	app.DELETE("/projects/:id", ProjectsRemove)
}

func ProjectsFetchOne(c *gin.Context) {
	id := c.Param("id")

	db := db.Database()
	defer db.Close()

	var proyecto Projects
	if err := db.Find(&proyecto, id).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		db.Model(&proyecto).Related(&proyecto.Owner, "Owner_id")
		c.JSON(http.StatusOK, proyecto)
	}
}

func ProjectsFetchAll(c *gin.Context) {

	db := db.Database()
	defer db.Close()

	var proyectos []Projects
	if err := db.Find(&proyectos).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		for i := range proyectos {
			db.Model(&proyectos[i]).Related(&proyectos[i].Owner, "Owner_id")
		}
		c.JSON(http.StatusOK, proyectos)
	}
}

func ProjectsCreate(c *gin.Context) {
	var proyecto Projects
	e := c.BindJSON(&proyecto)
	utils.Check(e)

	db := db.Database()
	defer db.Close()
	
	// Por ahora fijo la fecha actual como creacion
	fecha_actual := time.Now()
	proyecto.Date_created = &fecha_actual
	
	if err := db.Create(&proyecto).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		db.Model(&proyecto).Related(&proyecto.Owner, "Owner_id")
		c.JSON(http.StatusCreated, proyecto)
	}
}

func ProjectsUpdate(c *gin.Context) {
	var proyectos Projects
	id := c.Params.ByName("id")

	db := db.Database()
	defer db.Close()

	if err := db.Where("id = ?", id).First(&proyectos).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		c.BindJSON(&proyectos)

		db.Save(&proyectos)
		c.JSON(200, proyectos)
	}
}

func ProjectsRemove(c *gin.Context) {
	var proyectos Projects

	db := db.Database()
	defer db.Close()

	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&proyectos).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		db.Delete(&proyectos)
		c.JSON(http.StatusOK, proyectos)
	}
}



