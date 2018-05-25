package models

import (
	"github.com/vsepulve/gdrift-backend/db"
	"github.com/vsepulve/gdrift-backend/utils"
	"github.com/gin-gonic/gin"
	"fmt"
	"net"
	"strings"
	"encoding/binary"
	"io"
	"bytes"
	"net/http"
	"encoding/json"
//	"time"
	"strconv"
)

func SimulationCommandsCRUD(app *gin.Engine) {
	
	// Servicios reales del nuevo model
	
	// Iniciar Simulacion
	//   - Recibe un json con el escenario y los datos del proyecto
	//   - La idea es que bastara con el id o que los datos vengan bien estructurados
	//   - Agrega una Simulacion asociada al proyecto
	//   - Activa el servicio C++ de inicio de simulaciones
	//   - Si todo sale bien, agrega la simulacion a la BD
	//   - Retorna el mismo json agregando datos adicionales (id primero)
	app.POST("/simulation-command/", StartSimulation)
	
	// Detener Simulacion
	app.DELETE("/simulation-command/:id", StopSimulation)
	
	// Consultar Simulacion
	app.GET("/simulation-command/:id", QuerySimulation)
	
	// Quizas sea razonable agregar servicios relacionados a projects y simulations por separado, algo como:
	//   - (1) POST a /project/ para iniciar
	//   - (2) DELETE para detener el projecto completo
	//   - (3) GET para consultar estado del proyecto (esto es opcional)
	//   - (4) POST de /simulation/ para agregar y lanzar simulacion
	//   - (5) DELETE de /simulation/ para detener simulacion
	//   - (6) GET de /simulation/ para consultar resultados de simulacion
	// Los resultados de varias simulaciones se pueden componer en el json de respuesta del proyecto
	// Las simulaciones de un proyecto las puedo consultar a la BD
	// Los valores (o graficos) por simulacion los puedo obtener por socket del C++
	// Quzias todos los valores de graficos los pueda obtener uno a uno en float de 4 bytes
	// Tambien deberia definir los tipos de request desde aca y considerando la pagina
	// Por ahora para ser razonable usar los 6 tipos de arriba (de 1 a 6)
	
}

func StartSimulation(c *gin.Context) {

	fmt.Printf("StartSimulation - Inicio\n")
	
	// Json con la simulacion de entrada (de la pagina: paso 2)
	var simulacion Simulations
	c.BindJSON(&simulacion)
	
	// Conexion a BD
	db := db.Database()
	defer db.Close()
	
	// Agrego la simulacion a la BD
	// Notar que de este modo, quizas seria razonable borrar de la BD si hay errores
	fmt.Printf("StartSimulation - Agregando a la BD (para obtener Id)\n")
	if err := db.Create(&simulacion).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} else {
//		db.Model(&simulacion).Related(&simulacion.Project, "Project_id")
	}
	
	// Comunicacion con el demonio c++
	fmt.Printf("StartSimulation - Comunicando con C++ (%s, %s)\n", utils.Config.Daemon.Ip, utils.Config.Daemon.Port)
	connection, err := net.Dial("tcp", utils.Config.Daemon.Ip+":"+utils.Config.Daemon.Port)
	if err != nil {
//		fmt.Println(error)
		fmt.Printf("StartSimulation - Error al conectar con Daemon\n")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer connection.Close()
	
	// Datos que deben ser enviados:
	//   - Requsst type (1 byte, value = 4)
	//   - Project Id (4 bytes)
	//   - Simulation Id (4 bytes)
	//   - Json completo de simulacion (como string largo + bytes)
	request_type := []byte{2}
	bytes_int := make([]byte, 4)
	
	// Envio los datos de simulacion para agregarla a la cola de trabajo
	fmt.Printf("StartSimulation - Enviando datos (request type)\n")
	connection.Write(request_type)
	
	fmt.Printf("StartSimulation - Enviando Id de Proyecto (%d)\n", simulacion.Id)
	binary.LittleEndian.PutUint32(bytes_int, uint32(simulacion.Project_id))
	connection.Write(bytes_int)
	
	fmt.Printf("StartSimulation - Enviando Id de Simulacion (%d)\n", simulacion.Id)
	binary.LittleEndian.PutUint32(bytes_int, uint32(simulacion.Id))
	connection.Write(bytes_int)
	
	json_text, _ := json.MarshalIndent(simulacion, "", "\t")
	fmt.Printf("StartSimulation - Data: %s\n", json_text)
	message := string(json_text)
	length := len(message)
	
	fmt.Printf("StartSimulation - Enviando length (%d)\n", length)
	binary.LittleEndian.PutUint32(bytes_int, uint32(length))
	connection.Write(bytes_int)
	
	fmt.Printf("StartSimulation - Enviando mensaje\n")
	connection.Write([]byte(strings.TrimRight(message, "\n")))
	
	// Espero respuesta
	fmt.Printf("StartSimulation - Recibiendo respuesta\n")
	var buf bytes.Buffer
	io.Copy(&buf, connection)
	resp_code := binary.LittleEndian.Uint32(buf.Bytes())
	fmt.Printf("StartSimulation - resp_code: %d\n", resp_code)
	
	// Si hay problemas, envio codigo y salgo
	if resp_code != 1 {
		fmt.Printf("StartSimulation - Error al recibir respuesta\n")
		c.String(http.StatusInternalServerError, "Error")
	
	} else{
		// Respondo con el proyecto actualizado
		fmt.Printf("StartSimulation - Terminando\n")
		c.JSON(http.StatusCreated, simulacion)
	}
	
	fmt.Printf("StartSimulation - Fin\n")
}

func StopSimulation(c *gin.Context) {

	id := c.Param("id")
	sim_id, _ := strconv.Atoi(id)

	fmt.Printf("StopSimulation - Inicio (id: \"%s\" -> %d)\n", id, sim_id)
	
	// Conexion a BD (para verificar que la simulacion siga corriendo?, lo dejo para despues)
	// Quzias tambien para marcar la simulacion como detenida
//	db := db.Database()
//	defer db.Close()

	// Comunicacion con el demonio c++
	fmt.Printf("StopSimulation - Comunicando con C++ (%s, %s)\n", utils.Config.Daemon.Ip, utils.Config.Daemon.Port)
	connection, err := net.Dial("tcp", utils.Config.Daemon.Ip+":"+utils.Config.Daemon.Port)
	if err != nil {
//		fmt.Println(error)
		fmt.Printf("StopSimulation - Error al conectar con Daemon\n")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer connection.Close()
	
	// Datos que deben ser enviados:
	//   - Requsst type (1 byte, value = 5)
	//   - Simulation Id (4 bytes)
	request_type := []byte{3}
	bytes_int := make([]byte, 4)
	
	// Envio los datos de simulacion para agregarla a la cola de trabajo
	fmt.Printf("StopSimulation - Enviando datos (request type)\n")
	connection.Write(request_type)
	
	fmt.Printf("StopSimulation - Enviando Id de Simulacion (%d)\n", sim_id)
	binary.LittleEndian.PutUint32(bytes_int, uint32(sim_id))
	connection.Write(bytes_int)
	
	// Espero respuesta
	fmt.Printf("StopSimulation - Recibiendo respuesta\n")
	var buf bytes.Buffer
	io.Copy(&buf, connection)
	resp_code := binary.LittleEndian.Uint32(buf.Bytes())
	fmt.Printf("StopSimulation - resp_code: %d\n", resp_code)
	
	// Si hay problemas, envio codigo y salgo
	if resp_code != 1 {
		fmt.Printf("StopSimulation - Error al recibir respuesta\n")
		c.String(http.StatusInternalServerError, "Error")
	
	} else{
		// Respondo con el proyecto actualizado
		fmt.Printf("StopSimulation - Terminando\n")
		c.String(http.StatusOK, "")
	}
	
	fmt.Printf("StopSimulation - Fin\n")
}

type Point struct {
	X float32
	Y float32
}

type Curve struct {
	Name string
	Size int
	// Version 1: arreglos independientes (la otra opcion serian pares)
//	Data_x []float32
//	Data_y []float32
	// Version 2: Pares
	Data []Point
}

type Graph struct {
	Title string
	Curves []Curve
}

func QuerySimulation(c *gin.Context) {

	id := c.Param("id")
	sim_id, _ := strconv.Atoi(id)

	fmt.Printf("QuerySimulation - Inicio (id: \"%s\" -> %d)\n", id, sim_id)
	
	// Por ahora preparo un grafico dummy para retornar
	var grafico Graph
	grafico.Title = "Test Graph"
	grafico.Curves = make([]Curve, 0)
	
	var curva Curve
	curva.Name = "Curve 1"
	curva.Size = 3
	// Version 1
//	curva.Data_x = make([]float32, 0)
//	curva.Data_x = append(curva.Data_x, 0.1)
//	curva.Data_x = append(curva.Data_x, 0.5)
//	curva.Data_x = append(curva.Data_x, 0.9)
//	curva.Data_y = make([]float32, 0)
//	curva.Data_y = append(curva.Data_y, 0.1)
//	curva.Data_y = append(curva.Data_y, 0.3)
//	curva.Data_y = append(curva.Data_y, 0.85)
	// Version 2
	curva.Data = make([]Point, 0)
	var punto Point
	punto.X = 0.1
	punto.Y = 0.1
	curva.Data = append(curva.Data, punto)
	punto.X = 0.5
	punto.Y = 0.3
	curva.Data = append(curva.Data, punto)
	punto.X = 0.9
	punto.Y = 0.85
	curva.Data = append(curva.Data, punto)
	
	grafico.Curves = append(grafico.Curves, curva)
	
	curva.Name = "Curve 2"
	curva.Size = 4
	// Version 1
//	curva.Data_x = make([]float32, 0)
//	curva.Data_x = append(curva.Data_x, 0.1)
//	curva.Data_x = append(curva.Data_x, 0.5)
//	curva.Data_x = append(curva.Data_x, 0.9)
//	curva.Data_y = make([]float32, 0)
//	curva.Data_y = append(curva.Data_y, 0.1)
//	curva.Data_y = append(curva.Data_y, 0.3)
//	curva.Data_y = append(curva.Data_y, 0.85)
	// Version 2
	curva.Data = make([]Point, 0)
	punto.X = 0.2
	punto.Y = 0.6
	curva.Data = append(curva.Data, punto)
	punto.X = 0.4
	punto.Y = 0.5
	curva.Data = append(curva.Data, punto)
	punto.X = 0.6
	punto.Y = 0.3
	curva.Data = append(curva.Data, punto)
	punto.X = 0.8
	punto.Y = 0.1
	curva.Data = append(curva.Data, punto)
	
	grafico.Curves = append(grafico.Curves, curva)
	
	c.JSON(http.StatusOK, grafico)

}


