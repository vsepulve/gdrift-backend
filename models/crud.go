package models

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net"
	"strings"
	"encoding/binary"
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
	
	connection, error := net.Dial("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}

//	client := &Client{socket: connection}
//	go client.receive()

	for {
//		reader := bufio.NewReader(os.Stdin)
//		message, _ := reader.ReadString('\n')
		
		// Envio el tipo de request
		request_type := []byte{7}
		connection.Write(request_type)
		
		// Envio el mensaje
		// Notar que esto, si o si, implica convertir un numero en binario
		
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, 4)
		connection.Write(bs)
		
		message := "test"
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
	

}


func StartSimulation(c *gin.Context) {
	
	fmt.Printf("StartSimulation - Inicio\n")


}
