package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	pb "hextech/grpc-server/proto"

	"google.golang.org/grpc"
)

// Compara el reloj del supervisor con el de un servidor
func compareClocks(supervisorClock, serverClock []int32) int {

	inx := -1 //Si es -1 los relojes son consistentes

	for i := 0; i < len(supervisorClock); i++ {
		if supervisorClock[i] > serverClock[i] { //Servidor esta atrasado
			inx = i
			break
		}
	}
	return inx
}

// Inicializa un reloj vectorial de tamaño fijo
func initializeClock(size int) []int32 {
	return make([]int32, size)
}

// Le avisa al broker una inconsistencia en el servidor de una region
func reportarInconsistencia(conn *grpc.ClientConn, region string) error {
	client := pb.NewBrokerClient(conn)

	// Crear la solicitud
	req := &pb.InconsistenciaRequest{
		Region:  region,
		Mensaje: "Los relojes vectoriales son inconsistentes",
	}

	// Llamar al broker
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.NotificarInconsistencia(ctx, req)
	if err != nil {
		return fmt.Errorf("error notificando inconsistencia al broker: %v", err)
	}

	fmt.Printf("Inconsistencia en la región %s reportada al broker\n", region)
	return nil
}

func main() {

	// Conectar al broker en localhost:50051
	conn, err := grpc.Dial("broker:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Error conectando al Broker: %v\n", err)
		return
	}
	defer conn.Close()

	client := pb.NewBrokerClient(conn)

	reader := bufio.NewReader(os.Stdin)

	// Mapa para almacenar relojes vectoriales por región
	regionClocks := make(map[string][]int32)

	//ultimo servidor donde se escribio
	last_server := -1

	for {
		fmt.Println("\nIngrese un comando:")
		fmt.Println("(1) AgregarProducto <nombre region> <nombre producto> [valor]")
		fmt.Println("(2) RenombrarProducto <nombre region> <nombre producto> <nuevo nombre>")
		fmt.Println("(3) ActualizarValor <nombre region> <nombre producto> <nuevo valor>")
		fmt.Println("(4) BorrarProducto <nombre region> <nombre producto>")
		fmt.Println("(5) Salir")

		// Leer el comando del usuario
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Validar el comando (no es lo mejor que sirve)
		parts := strings.Fields(input)

		if parts[0] == "Salir" {
			break
		}

		if len(parts) < 3 {
			fmt.Println("Error: Formato inválido. Intente nuevamente.")
			continue
		}

		command := parts[0]
		region := parts[1]
		productName := parts[2]
		value := ""
		if len(parts) > 3 {
			value = strings.Join(parts[3:], " ")
		}

		// Validar el nombre del comando
		if command != "AgregarProducto" && command != "RenombrarProducto" &&
			command != "ActualizarValor" && command != "BorrarProducto" {
			fmt.Println("Error: Comando desconocido. Intente nuevamente.")
			continue
		}

		// Acceder al reloj vectorial de la region "region"
		clock, exists := regionClocks[region]
		if !exists {
			fmt.Println("La región 'Noxus' no existe en el mapa.")
			regionClocks[region] = initializeClock(3) // Inicializar si no existe
			//clock = regionClocks[region]

		} else {
			fmt.Printf("Reloj actual de %s para supervisor 1: %v\n", region, clock)
		}

		// Enviar el comando al broker y conectar a servidor aleatorio
		// Solicitar el reloj vectorial al broker
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req := &pb.ClockRequest{
			Region: region,
		}

		res, err := client.ObtenerReloj(ctx, req)
		if err != nil {
			fmt.Printf("Error al obtener producto: %v\n", err)
			continue
		}

		fmt.Printf("Reloj vectorial recibido del broker: %v\n", res.Clock)
		fmt.Printf("Servidor a conectar en supervisor 1: servidor_%d\n", res.Servidor)

		//Verificar consistencia en los relojes
		consistencia := compareClocks(regionClocks[region], res.Clock)
		if consistencia != -1 {
			//Inconsistencia, mandar al ultimo servidor que escribió
			fmt.Printf("Relojes inconsistentes! redireccionando a servidor anterior...\n")

			// Reportar inconsistencia al broker
			err := reportarInconsistencia(conn, region)
			if err != nil {
				fmt.Printf("Error reportando inconsistencia: %v\n", err)
			}

		} else {
			//Consistente, escribir directamente
			last_server = int(res.Servidor)
			regionClocks[region] = res.Clock

		}

		// Construir el request
		req2 := &pb.CommandRequest{
			Command:     command,
			Region:      region,
			ProductName: productName,
			Value:       value,
			Server:      int32(last_server),
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		res2, err := client.ProcessCommand(ctx2, req2)
		if err != nil {
			fmt.Printf("Error al enviar el comando: %v\n", err)
			continue
		}

		// Imprimir la respuesta del broker
		fmt.Printf("Comando enviado correctamente! : %s\n", res2.Message)

		// Incrementar el primer valor del reloj de la región "Noxus"
		regionClocks[region][last_server]++ //last_server es id del server Ej: Noxus -> [0,0,+1]

		// Imprimir el reloj actualizado
		fmt.Printf("Reloj actualizado de servidor %d: %v\n", last_server, regionClocks[region])
	}
}
