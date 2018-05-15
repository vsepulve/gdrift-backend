package models

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net"
	"strings"
	"encoding/binary"
	"io"
	"bytes"
	"net/http"
	"encoding/json"
	"strconv"
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
//	app.POST("/create-project/", CreateProject)
	
	// Iniciar Simulacion
	//   - Recibe un json con el escenario y los datos del proyecto
	//   - La idea es que bastara con el id o que los datos vengan bien estructurados
	//   - Agrega una Simulacion asociada al proyecto
	//   - Activa el servicio C++ de inicio de simulaciones
	//   - Si todo sale bien, agrega la simulacion a la BD
	//   - Retorna el mismo json agregando datos adicionales (id primero)
//	app.POST("/start-simulation/", StartSimulation)
	
	
	
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




