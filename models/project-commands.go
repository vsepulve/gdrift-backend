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

func ProjectCommandsCRUD(app *gin.Engine) {
	
	// Servicios reales del nuevo model
	
	// Crear Proyecto
	//   - Recibe un json con los datos generales del proyecto (inlutendo samples)
	//   - Por ahora supongo que los samples estan en archivos y que el json incluye las rutas
	//   - El json recivido puede ser de tipo "Projects" (revisar "Individual_data" para los datos de la especie)
	//   - Responde el json agregando datos adicionales (id primero que nada)
	//   - Activa el servicio C++ de creacion de target del proyecto
	//   - Notar que el C++ almacena el json del proyecto (con el id como nombre), luego el Factory usa por separado ese json para generar el perfil y el de la simulacion para los eventos
	app.POST("/project-command/", CreateProject)
	
	// Detener Proyecto
	app.DELETE("/project-command/:id", StopProject)
	
	// Consultar Proyecto
	app.GET("/project-command/:id", QueryProject)
	
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

func StopProject(c *gin.Context) {

	id := c.Param("id")
	proj_id, _ := strconv.Atoi(id)

	fmt.Printf("StopProject - Inicio (id: \"%s\" -> %d)\n", id, proj_id)

}
	

func QueryProject(c *gin.Context) {

	id := c.Param("id")
	proj_id, _ := strconv.Atoi(id)

	fmt.Printf("QueryProject - Inicio (id: \"%s\" -> %d)\n", id, proj_id)
	
}


