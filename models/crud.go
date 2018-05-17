package models

import (
	"github.com/vsepulve/gdrift-backend/db"
//	"github.com/vsepulve/gdrift-backend/utils"
	"github.com/gin-gonic/gin"
	"fmt"
	"net"
	"strings"
	"encoding/binary"
	"io"
	"bytes"
	"net/http"
//	"encoding/json"
//	"strconv"
	"time"
)

func Setup(app *gin.Engine) {
	UsersCRUD(app)
	ProjectsCRUD(app)
	
	// Servicios para pruebas de comunicacion con C++
	app.GET("/simulate/", TestSimulation)
	app.POST("/simulate/", StartSimulation)
	
	// Servicios reales del nuevo model
	
	// Crear Proyecto
	//   - Recibe un json con los datos generales del proyecto (inlutendo samples)
	//   - Por ahora supongo que los samples estan en archivos y que el json incluye las rutas
	//   - El json recivido puede ser de tipo "Projects" (revisar "Individual_data" para los datos de la especie)
	//   - Responde el json agregando datos adicionales (id primero que nada)
	//   - Activa el servicio C++ de creacion de target del proyecto
	app.POST("/create-project/", CreateProject)
	
	// Iniciar Simulacion
	//   - Recibe un json con el escenario y los datos del proyecto
	//   - La idea es que bastara con el id o que los datos vengan bien estructurados
	//   - Agrega una Simulacion asociada al proyecto
	//   - Activa el servicio C++ de inicio de simulaciones
	//   - Si todo sale bien, agrega la simulacion a la BD
	//   - Retorna el mismo json agregando datos adicionales (id primero)
	app.POST("/start-simulation/", StartSimulation)
	
	
	
}

func CreateProject(c *gin.Context) {

	fmt.Printf("CreateProject - Inicio\n")
	
	// Json con el proyecto de entrada (de la pagina: paso 1)
	var proyecto Projects
	c.BindJSON(&proyecto)
	
	// Por ahora fijo la fecha actual como creacion
	fmt.Printf("CreateProject - Agregando Date_created (now)\n")
	fecha_actual := time.Now()
	proyecto.Date_created = &fecha_actual
	
	// Conexion a BD
	db := db.Database()
	defer db.Close()
	
	// Agrego el proyecto a la BD
	// Notar que de este modo, quizas seria razonable borrar de la BD si hay errores
	fmt.Printf("CreateProject - Agregando a la BD (para obtener Id)\n")
	if err := db.Create(&proyecto).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} else {
		db.Model(&proyecto).Related(&proyecto.Owner, "Owner_id")
	}
	
	/*
	// Desactivado mientras trabajo en el Daemon
	
	// Comunicacion con el demonio c++
	fmt.Printf("CreateProject - Comunicando con C++ (%s, %s)\n", utils.Config.Daemon.Ip, utils.Config.Daemon.Port)
	connection, err := net.Dial("tcp", utils.Config.Daemon.Ip+":"+utils.Config.Daemon.Port)
	if err != nil {
//		fmt.Println(error)
		fmt.Printf("CreateProject - Error al conectar con Daemon\n")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer connection.Close()
	
	// Envio los datos del proyecto para su inicializacion (crear target)
	fmt.Printf("CreateProject - Enviando datos (request type)\n")
	request_type := []byte{1}
	connection.Write(request_type)
	
	// Otros datos...
	
	// Espero respuesta
	fmt.Printf("CreateProject - Recibiendo respuesta\n")
	var buf bytes.Buffer
	io.Copy(&buf, connection)
	resp_code := binary.LittleEndian.Uint32(buf.Bytes())
	fmt.Printf("CreateProject - resp_code: %d\n", resp_code)
	*/
	
	// Si hay problemas, envio codigo y salgo
	
	// Respondo con el proyecto actualizado
	fmt.Printf("CreateProject - Terminando\n")
	c.JSON(http.StatusCreated, proyecto)
	
	fmt.Printf("CreateProject - Fin\n")
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
	
	/*
	// Desactivado mientras trabajo en el Daemon
	
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
	
	// Envio los datos de simulacion para su inicializacion (crear target)
	fmt.Printf("StartSimulation - Enviando datos (request type)\n")
	request_type := []byte{2}
	connection.Write(request_type)
	
	// Otros datos...
	
	// Espero respuesta
	fmt.Printf("StartSimulation - Recibiendo respuesta\n")
	var buf bytes.Buffer
	io.Copy(&buf, connection)
	resp_code := binary.LittleEndian.Uint32(buf.Bytes())
	fmt.Printf("StartSimulation - resp_code: %d\n", resp_code)
	*/
	
	// Si hay problemas, envio codigo y salgo
	
	// Respondo con la simulacion actualizado
	fmt.Printf("StartSimulation - Terminando\n")
	c.JSON(http.StatusCreated, simulacion)
	
	fmt.Printf("StartSimulation - Fin\n")
}


func TestSimulation(c *gin.Context) {
	
	fmt.Printf("TestSimulation - Inicio\n")

	for i := 0; i < 10; i++ {
		
	
		connection, error := net.Dial("tcp", "localhost:12345")
		if error != nil {
			fmt.Println(error)
		}
		defer connection.Close()
	
		// Envio el tipo de request (un byte con un valor arbitrario, 7 en este caso)
		fmt.Printf("\nTestSimulation - Prueba %d, enviando request_type (7)\n", i)
		request_type := []byte{7}
		connection.Write(request_type)
		
		// Envio el mensaje de prueba (string en el formato length + chars)
		// Notar que esto, si o si, implica convertir un numero en binario
		message := "Prueba"
		length := len(message)
		
		fmt.Printf("TestSimulation - Enviando length (%d)\n", length)
		bytes_int := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes_int, uint32(length))
		connection.Write(bytes_int)
		
		fmt.Printf("TestSimulation - Enviando mensaje (\"%s\")\n", message)
		connection.Write([]byte(strings.TrimRight(message, "\n")))
		
		
		// Recibo la respuesta de prueba (un entero en 4 bytes)
		fmt.Printf("TestSimulation - Leyendo codigo de respuesta\n")
		
		// Metodo 1 para leer bytes de connection
//		n, err := connection.Read(bytes_int);
//		if err != nil {
//			if err != io.EOF {
//				fmt.Println("read error:", err)
//			}
//			break
//		}
//		resp_code := binary.LittleEndian.Uint32(bytes_int)
//		fmt.Printf("TestSimulation - Respuesta: %d (%d bytes)\n", resp_code, n)
		
		// Metodo 2 para leer bytes de connection
		var buf bytes.Buffer
		io.Copy(&buf, connection)
		resp_code := binary.LittleEndian.Uint32(buf.Bytes())
		fmt.Printf("TestSimulation - Respuesta: %d (%d bytes)\n", resp_code, buf.Len())
		
		
		
		
	}
	
	fmt.Printf("\nTestSimulation - Fin\n")
	
	data := make(map[string]string)
	
	data["Campo Prueba"] = "Valor Prueba"
	data["Estado"] = "Ok"
	c.JSON(http.StatusOK, data)
	
}

/*
func StartSimulation(c *gin.Context) {
	
	fmt.Printf("StartSimulation - Inicio\n")
	
//	var data interface{}
	var data map[string]interface{}
	e := c.BindJSON(&data)
	if e != nil {
		panic(e)
	}
	
	json_text, e := json.MarshalIndent(data, "", "\t")
//	json_text, e := json.Marshal(data)
	
	sim_id, e := strconv.Atoi( (data["id"]).(string) )
	n_sims, e := strconv.Atoi( (data["batch-size"]).(string) )
	
	fmt.Printf("StartSimulation - sim_id: %d, n_sims: %d\n", sim_id, n_sims)
	fmt.Printf("StartSimulation - Data: %s\n", json_text)
	
	connection, error := net.Dial("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}
	defer connection.Close()

	// Envio el tipo de request (un byte con un valor arbitrario, 1 en este caso)
	fmt.Printf("\nStartSimulation - Enviando request_type (1)\n")
	request_type := []byte{1}
	connection.Write(request_type)
	
	// Bytes para entero
	bytes_int := make([]byte, 4)
	
	// Envio sim_id
	binary.LittleEndian.PutUint32(bytes_int, uint32(sim_id))
	connection.Write(bytes_int)
	
	// Envio n_sims
	binary.LittleEndian.PutUint32(bytes_int, uint32(n_sims))
	connection.Write(bytes_int)
	
	// Envio del json en el formato len + chars
//	message := "Prueba"
	message := string(json_text)
	length := len(message)
	
	fmt.Printf("StartSimulation - Enviando length (%d)\n", length)
	binary.LittleEndian.PutUint32(bytes_int, uint32(length))
	connection.Write(bytes_int)
	
	fmt.Printf("StartSimulation - Enviando mensaje\n")
	connection.Write([]byte(strings.TrimRight(message, "\n")))
	
	// Recibo la respuesta de prueba (un entero en 4 bytes)
	fmt.Printf("StartSimulation - Leyendo codigo de respuesta\n")
	
	var buf bytes.Buffer
	io.Copy(&buf, connection)
	resp_code := binary.LittleEndian.Uint32(buf.Bytes())
	fmt.Printf("TestSimulation - Respuesta: %d (%d bytes)\n", resp_code, buf.Len())
	
	fmt.Printf("\n TestSimulation - Fin\n")
	
	resp := make(map[string]string)
	
	resp["Campo Prueba"] = "Valor Prueba"
	resp["Estado"] = "Ok"
	c.JSON(http.StatusOK, resp)

}
*/



