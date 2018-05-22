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
	"time"
	"strconv"
)

func CommandsCRUD(app *gin.Engine) {
	
	// Servicios reales del nuevo model
	
	// Crear Proyecto
	//   - Recibe un json con los datos generales del proyecto (inlutendo samples)
	//   - Por ahora supongo que los samples estan en archivos y que el json incluye las rutas
	//   - El json recivido puede ser de tipo "Projects" (revisar "Individual_data" para los datos de la especie)
	//   - Responde el json agregando datos adicionales (id primero que nada)
	//   - Activa el servicio C++ de creacion de target del proyecto
	//   - Notar que el C++ almacena el json del proyecto (con el id como nombre), luego el Factory usa por separado ese json para generar el perfil y el de la simulacion para los eventos
	app.POST("/project/", CreateProject)
	
	// Detener Proyecto
//	app.DELETE("/project/:id", StopProject)
	
	// Consultar Proyecto
//	app.GET("/project/:id", QueryProject)
	
	
	// Iniciar Simulacion
	//   - Recibe un json con el escenario y los datos del proyecto
	//   - La idea es que bastara con el id o que los datos vengan bien estructurados
	//   - Agrega una Simulacion asociada al proyecto
	//   - Activa el servicio C++ de inicio de simulaciones
	//   - Si todo sale bien, agrega la simulacion a la BD
	//   - Retorna el mismo json agregando datos adicionales (id primero)
	app.POST("/simulation/", StartSimulation)
	
	// Detener Simulacion
	app.DELETE("/simulation/:id", StopSimulation)
	
	// Consultar Proyecto
//	app.GET("/simulation/:id", QuerySimulation)
	
	
	// Quizas sea razonable agregar servicios relacionados a projects y simulations por separado
	//   - Algo como POST a /project/ para iniciar
	//   - DELETE para detener el projecto completo
	//   - GET para consultar estado del proyecto (esto es opcional)
	//   - POST de /simulation/ para agregar y lanzar simulacion
	//   - DELETE de /simulation/ para detener simulacion
	//   - GET de /simulation/ para consultar resultados de simulacion
	// Los resultados de varias simulaciones se pueden componer en el json de respuesta del proyecto
	// Las simulaciones de un proyecto las puedo consultar a la BD
	// Los valores (o graficos) por simulacion los puedo obtener por socket del C++
	// Quzias todos los valores de graficos los pueda obtener uno a uno en float de 4 bytes
	
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
	
	// Datos que deben ser enviados:
	//   - Requsst type (1 byte, value = 1)
	//   - Project Id (4 bytes)
	//   - Json completo de proyecto (como string largo + bytes)
	request_type := []byte{1}
	bytes_int := make([]byte, 4)
	
	// Envio los datos del proyecto para su inicializacion (crear target)
	fmt.Printf("CreateProject - Enviando datos (request type)\n")
	connection.Write(request_type)
	
	fmt.Printf("CreateProject - Enviando Id de Proyecto (%d)\n", proyecto.Id)
	binary.LittleEndian.PutUint32(bytes_int, uint32(proyecto.Id))
	connection.Write(bytes_int)
	
	json_text, _ := json.MarshalIndent(proyecto, "", "\t")
	fmt.Printf("CreateProject - Data: %s\n", json_text)
	message := string(json_text)
	length := len(message)
	
	fmt.Printf("CreateProject - Enviando length (%d)\n", length)
	binary.LittleEndian.PutUint32(bytes_int, uint32(length))
	connection.Write(bytes_int)
	
	fmt.Printf("CreateProject - Enviando mensaje\n")
	connection.Write([]byte(strings.TrimRight(message, "\n")))
	
	// Espero respuesta
	fmt.Printf("CreateProject - Recibiendo respuesta\n")
	var buf bytes.Buffer
	io.Copy(&buf, connection)
	resp_code := binary.LittleEndian.Uint32(buf.Bytes())
	fmt.Printf("CreateProject - resp_code: %d\n", resp_code)
	
	// Si hay problemas, envio codigo y salgo
	if resp_code != 1 {
		fmt.Printf("CreateProject - Error al recibir respuesta\n")
		c.String(http.StatusInternalServerError, "Error")
	
	} else{
		// Respondo con el proyecto actualizado
		fmt.Printf("CreateProject - Terminando\n")
		c.JSON(http.StatusCreated, proyecto)
	}
	
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
	//   - Requsst type (1 byte, value = 2)
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

	fmt.Printf("StopSimulation - Inicio\n")
	
	id := c.Param("id")
	sim_id, _ := strconv.Atoi(id)
	
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
	//   - Requsst type (1 byte, value = 3)
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



