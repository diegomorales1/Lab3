package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	pb "hextech/grpc-server/proto"

	"google.golang.org/grpc"
)

type brokerServer struct {
	pb.UnimplementedBrokerServer
	clients []pb.BrokerClient
}

// getRandomClient selecciona un cliente al azar de la lista.
func (b *brokerServer) getRandomClient() pb.BrokerClient {
	index := rand.Intn(3)
	return b.clients[index]
}

func (b *brokerServer) NotificarInconsistencia(ctx context.Context, req *pb.InconsistenciaRequest) (*pb.InconsistenciaResponse, error) {
	// Registrar la inconsistencia en los logs del broker
	logMessage := fmt.Sprintf("Inconsistencia detectada en la región: %s. Detalles: %s\n", req.Region, req.Mensaje)
	fmt.Print(logMessage)

	return &pb.InconsistenciaResponse{
		Status: "Inconsistencia registrada correctamente",
	}, nil
}

func (b *brokerServer) GetRandServer(ctx context.Context, req *pb.EmptyRequest) (*pb.RandServerResponse, error) {
	index := rand.Intn(3)

	return &pb.RandServerResponse{
		Server: int32(index),
	}, nil
}

// ObtenerProducto reenvía la solicitud al servidor seleccionado aleatoriamente
func (b *brokerServer) ObtenerProducto(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	fmt.Printf("Consulta recibida: Región: %s, Producto: %s\n", req.Region, req.ProductName)

	// Seleccionar el servidor aleatoriamente
	client := b.getRandomClient()

	// Reenviar la solicitud al servidor seleccionado
	res, err := client.ObtenerProducto(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error al obtener producto del servidor seleccionado: %v", err)
	}

	fmt.Printf("Respuesta del servidor: %s\n", res.Message)
	return res, nil
}

func (b *brokerServer) ObtenerReloj(ctx context.Context, req *pb.ClockRequest) (*pb.ClockResponse, error) {
	servers := []string{"server1:50052", "server2:50053", "server3:50054"}
	fmt.Printf("Consulta recibida: Región: %s\n", req.Region)

	// Seleccionar el servidor aleatoriamente
	index := rand.Intn(3)
	println("AQUI EL INDEX:", index)
	client := b.clients[index]

	if index == 0 {
		conn, err := grpc.Dial(servers[0], grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("Error al conectar con el servidor %s: %v\n", servers[0], err)
		}
		defer conn.Close() // Asegurarse de cerrar todas las conexiones al finalizar
		client = pb.NewBrokerClient(conn)
	}

	if index == 1 {
		conn, err := grpc.Dial(servers[1], grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("Error al conectar con el servidor %s: %v\n", servers[1], err)
		}
		defer conn.Close() // Asegurarse de cerrar todas las conexiones al finalizar
		client = pb.NewBrokerClient(conn)
	}

	if index == 2 {
		conn, err := grpc.Dial(servers[2], grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Printf("Error al conectar con el servidor %s: %v\n", servers[2], err)
		}
		defer conn.Close() // Asegurarse de cerrar todas las conexiones al finalizar
		client = pb.NewBrokerClient(conn)
	}

	response, err := client.GetClock(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error al obtener reloj del servidor seleccionado: %v", err)
	}

	fmt.Printf("Enviando reloj de servidor...")

	fmt.Printf("Servidor seleccionado: %d, Reloj: %v\n", response.Servidor, response.Clock)

	return &pb.ClockResponse{
		Clock:    response.Clock,
		Servidor: response.Servidor,
	}, nil
}

// ProcessCommand maneja los comandos recibidos y los reenvía al servidor seleccionado aleatoriamente
func (b *brokerServer) ProcessCommand(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	fmt.Printf("Comando recibido: %s, Región: %s, Producto: %s, Valor: %s\n",
		req.Command, req.Region, req.ProductName, req.Value)

	client := b.clients[req.Server]

	// Reenviar el comando al servidor seleccionado
	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.ProcessCommand(ctx2, req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar el comando al servidor seleccionado: %v", err)
	}

	fmt.Printf("Respuesta del servidor: %s\n", response.Message)
	return &pb.CommandResponse{
		RandomValue: response.RandomValue,
		Message:     "Comando reenviado al servidor exitosamente",
	}, nil
}

func main() {

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Error al iniciar el broker: %v\n", err)
		return
	}
	fmt.Println("Broker escuchando en el puerto 50051")

	// Puertos de los servidores
	rand.Seed(time.Now().UnixNano())
	servers := []string{"server1:50052", "server2:50053", "server3:50054"}
	var clients []pb.BrokerClient

	// Crear conexiones a los servidores
	for _, server := range servers {
		conn, err := grpc.Dial(server, grpc.WithInsecure())
		if err != nil {
			fmt.Printf("Error al conectar con el servidor %s: %v\n", server, err)
			return
		}
		defer conn.Close() // Asegurarse de cerrar todas las conexiones al finalizar
		clients = append(clients, pb.NewBrokerClient(conn))
	}

	fmt.Println("Iniciando servidor de broker")
	// Iniciar el servidor gRPC del broker
	server := grpc.NewServer()
	pb.RegisterBrokerServer(server, &brokerServer{clients: clients})
	if err := server.Serve(listener); err != nil {
		fmt.Printf("Error al servir: %v\n", err)
	}
	fmt.Println("Terminé en broker")
}
