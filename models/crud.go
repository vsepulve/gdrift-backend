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
)

func Setup(app *gin.Engine) {
//	UsuarioCRUD(app)
	
	// Servicios especificos para gdrift
	
	// Servicio de entrega inicial de datos
	
	// Servicio de inicio de simulacion
	app.GET("/simulate/", TestSimulation)
	app.POST("/simulate/", StartSimulation)
	
}

type Client struct {
	socket net.Conn
	data   chan []byte
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


}
