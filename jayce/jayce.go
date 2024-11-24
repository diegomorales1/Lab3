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

// Estructura para almacenar información de consistencia
type ProductInfo struct {
	Value       string
	VectorClock []int32
	LastServer  string
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

func traducirId(ind int) string {
	if ind == 0 {
		return "dist038:50052"
	}
	if ind == 1 {
		return "dist039:50053"
	}

	return "dist040:50054"
}

func main() {
	fmt.Println("Cliente Jayce iniciado. Usa el comando: ObtenerProducto <Region> <Producto>")
	conn, err := grpc.Dial("dist037:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error conectando al Broker: %v\n", err)
		return
	}
	defer conn.Close()

	client := pb.NewBrokerClient(conn)

	localCache := make(map[string][]int32)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Cliente Jayce iniciado. Usa el comando: ObtenerProducto <Region> <Producto>")

	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if strings.HasPrefix(command, "ObtenerProducto") {
			parts := strings.Fields(command)
			if len(parts) != 3 {
				fmt.Println("Formato incorrecto. Usa: ObtenerProducto <Region> <Producto>")
				continue
			}

			region := parts[1]
			product := parts[2]

			// Enviar el comando al broker y conectar a servidor aleatorio
			// Solicitar el reloj vectorial al broker

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			res, err := client.GetRandServer(ctx, &pb.EmptyRequest{})
			if err != nil {
				fmt.Printf("Error al obtener servidor: %s\n", err)
				continue
			}

			fmt.Printf("Servidor recibido por Jayce: servidor %d\n", res.Server)

			conexion := traducirId(int(res.Server))

			conn2, err := grpc.Dial(conexion, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				fmt.Printf("Error conectando al Broker: %v\n", err)
				return
			}
			defer conn.Close()

			client2 := pb.NewBrokerClient(conn2)

			// Enviar solicitud al Broker
			ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			req2 := &pb.ClockRequest{
				Region: region,
			}

			res2, err := client2.GetClock(ctx2, req2)
			if err != nil {
				fmt.Printf("Error al obtener el reloj del servidor: %v\n", err)
				continue
			}

			newClock := res2.Clock

			// Acceder al reloj vectorial de la region "region"
			clock, exists := localCache[region]
			if !exists {
				fmt.Printf("La región %s no existe en el mapa. \n", region)
				localCache[region] = initializeClock(3) // Inicializar si no existe

			} else {
				fmt.Printf("Reloj actual de %s para Jayce: %v\n", region, clock)
			}

			//Verificar consistencia en los relojes
			consistencia := compareClocks(localCache[region], newClock)
			if consistencia != -1 {
				//Inconsistencia, mandar al ultimo servidor que leyó

				fmt.Printf("Relojes inconsistentes! tomando valor del cache...\n")

				// Reportar inconsistencia al broker
				err := reportarInconsistencia(conn, region)
				if err != nil {
					fmt.Printf("Error reportando inconsistencia: %v\n", err)
				}

			} else {
				//Consistente, conexion y lectura normal

				// Crear solicitud para el broker
				req := &pb.CommandRequest{
					Command:     "ObtenerProducto",
					Region:      region,
					ProductName: product,
					Value:       "0",
				}

				res3, err := client2.ObtenerProducto(ctx2, req)
				if err != nil {
					fmt.Printf("El producto %s no existe en la region %s o esta ultima no existe\n", product, region)
					continue
				}

				quantity := res3.Value
				if quantity == "" {
					quantity = "0"
				}

				fmt.Printf("Jayce: Producto: %s, Región: %s, Cantidad: %s\n", product, region, quantity)

				localCache[region] = res2.Clock
			}

		} else {
			fmt.Println("Comando no reconocido. Usa: ObtenerProducto <Region> <Producto>")

		}

	}
}
