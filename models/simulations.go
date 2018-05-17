package models

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"github.com/vsepulve/gdrift-backend/db"
	"github.com/vsepulve/gdrift-backend/utils"

)

// Notar que para un numero fijo, basta: "fixed" + [value]
//type Distribution struct {
//	Type                 string
//	Parameters           []float32
//}

// Notar que aqui no esta claro que hacer con tipos de marcadores diferentes
//type Event struct {
//	Id                   int
//	Type                 int
//	Generation           Distribution
//	// ...lo ideal seria usar listas de parametros como en el nuevo Simulator
//	// Pero por ahora creo que los dejare como interface sin tipo
//}

type Simulations struct {
	Id               int           `gorm:"column:id;not null;primary_key"`
	Project_id       int           `gorm:"column:projects_id;not null"`
	// Notar que, de hecho, basta con el id, no se necesita el projecto
//	Project          Projects      `gorm:"ForeignKey:Project_id;AssociationForeignKey:Id"`
	Model            int           `gorm:"column:model"`
	Events           interface{}   `gorm:"-"`
}


func SimulationsCRUD(app *gin.Engine) {
	app.GET("/simulations/:id", SimulationsFetchOne)
	app.GET("/simulations/", SimulationsFetchAll)
	app.POST("/simulations/", SimulationsCreate)
	app.PUT("/simulations/:id", SimulationsUpdate)
	app.DELETE("/simulations/:id", SimulationsRemove)
}

func SimulationsFetchOne(c *gin.Context) {
	id := c.Param("id")

	db := db.Database()
	defer db.Close()

	var simulacion Simulations
	if err := db.Find(&simulacion, id).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
//		db.Model(&simulacion).Related(&simulacion.Project, "Project_id")
		c.JSON(http.StatusOK, simulacion)
	}
}

func SimulationsFetchAll(c *gin.Context) {

	db := db.Database()
	defer db.Close()

	var simulaciones []Simulations
	if err := db.Find(&simulaciones).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
//		for i := range simulaciones {
//			db.Model(&simulaciones[i]).Related(&simulaciones[i].Project, "Project_id")
//		}
		c.JSON(http.StatusOK, simulaciones)
	}
}

func SimulationsCreate(c *gin.Context) {
	var simulacion Simulations
	e := c.BindJSON(&simulacion)
	utils.Check(e)

	db := db.Database()
	defer db.Close()
	
	if err := db.Create(&simulacion).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
//		db.Model(&simulacion).Related(&simulacion.Project, "Project_id")
		c.JSON(http.StatusCreated, simulacion)
	}
}

func SimulationsUpdate(c *gin.Context) {
	var simulaciones Simulations
	id := c.Params.ByName("id")

	db := db.Database()
	defer db.Close()

	if err := db.Where("id = ?", id).First(&simulaciones).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		c.BindJSON(&simulaciones)

		db.Save(&simulaciones)
		c.JSON(200, simulaciones)
	}
}

func SimulationsRemove(c *gin.Context) {
	var simulaciones Simulations

	db := db.Database()
	defer db.Close()

	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&simulaciones).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		db.Delete(&simulaciones)
		c.JSON(http.StatusOK, simulaciones)
	}
}



