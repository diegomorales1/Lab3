package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	pb "hextech/grpc-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Servidor Hextech
type hextechServer struct {
	pb.UnimplementedBrokerServer
	clock []int32    // Reloj vectorial
	mutex sync.Mutex // Para sincronizar el acceso al reloj
}

var ID_server = 2

// ProcessCommand maneja los comandos enviados por el broker
func (s *hextechServer) ProcessCommand(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Registrar el comando en Logs.txt
	logEntry := fmt.Sprintf("%s %s %s %s\n", req.Command, req.Region, req.ProductName, req.Value)
	err := appendToFile("Logs.txt", logEntry)
	if err != nil {
		return nil, nil
	}

	// Actualizar el archivo regiones.txt
	err = updateRegionClock("regiones.txt", req.Region, s.clock, ID_server)
	if err != nil {
		return nil, nil
	}

	fmt.Printf("Servidor Hextech 3: Comando registrado en Logs.txt: %s\n", logEntry)

	return &pb.CommandResponse{
		Message:        "Comando registrado en el log y regiones.txt actualizado",
		RelojVectorial: s.clock,
	}, nil
}

// Actualiza el reloj de la region especifica a realizar cambios
func updateRegionClock(filename, region string, clock []int32, serverID int) error {
	// Leer el archivo regiones.txt (si existe)
	data, err := readFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Inicializar el mapa de regiones y relojes
	regionClocks := make(map[string][]int32)

	// Parsear las líneas del archivo (si existe contenido)
	if data != "" {
		lines := strings.Split(data, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}

			// Parsear la región y su reloj
			regionName := parts[0]
			clockValues := parseClock(parts[1:])
			regionClocks[regionName] = clockValues
		}
	}

	// Actualizar o inicializar el reloj para la región
	if _, exists := regionClocks[region]; !exists {
		// Si la región no existe, inicializar en 0 0 0
		regionClocks[region] = make([]int32, len(clock))
	}

	// Incrementar el reloj en la posición correspondiente al servidor
	regionClocks[region][serverID]++

	// Escribir el archivo actualizado
	var lines []string
	for regionName, clockValues := range regionClocks {
		line := fmt.Sprintf("%s %s", regionName, formatClock(clockValues))
		lines = append(lines, line)
	}
	return writeFile(filename, strings.Join(lines, "\n"))
}

// funcion auxiliar para reestructurar el reloj
func parseClock(clockParts []string) []int32 {
	var clock []int32
	for _, part := range clockParts {
		val, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		clock = append(clock, int32(val))
	}
	return clock
}

// funcion auxiliar para formatear reloj
func formatClock(clock []int32) string {
	var parts []string
	for _, val := range clock {
		parts = append(parts, fmt.Sprintf("%d", val))
	}
	return strings.Join(parts, " ")
}

// Recibe las regiones combinadas de todos los servidores
func (s *hextechServer) ReceiveMergedFile2(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	content := string(req.Content) // Convertir los bytes recibidos a string

	// Sobrescribir el archivo regiones.txt con el contenido recibido
	err := writeFile("regiones.txt", content)
	if err != nil {
		log.Printf("Error al escribir regiones.txt: %v", err)
		return &pb.FileResponse{Content: "Error"}, nil
	}

	log.Printf("Contenido de regiones.txt recibido y actualizado exitosamente")
	return &pb.FileResponse{Content: "Success"}, nil
}

// Recibe el log combinado de todos los servidores
func (s *hextechServer) ReceiveMergedFile(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	// Convertir el contenido del archivo recibido de bytes a string
	mergedContent := string(req.Content)

	// Log para depuración
	fmt.Println("Archivo mergeado recibido:", mergedContent)

	// Procesar el contenido recibido
	procesarLogsFromContent(mergedContent)

	// Responder al cliente con un mensaje de éxito
	return &pb.FileResponse{Content: "Archivo procesado con éxito"}, nil
}

// Obtiene los archivos de de los logs
func (s *hextechServer) GetFile(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	filename := req.Filename

	// Abrir archivo
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo: %v", err)
	}
	defer file.Close()

	// Leer el archivo línea por línea y concatenar en un string
	var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error leyendo archivo: %v", err)
	}

	// Vaciar el log
	err = writeFile("Logs.txt", "")
	if err != nil {
		fmt.Printf("Error limpiando Logs.txt: %v\n", err)
	}

	// Retornar el contenido concatenado
	return &pb.FileResponse{
		Content: strings.Join(result, ", "),
	}, nil
}

// GetVectorClock devuelve un reloj vectorial [x, y, z]
func (s *hextechServer) GetClock(ctx context.Context, req *pb.ClockRequest) (*pb.ClockResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	clock, err := getRegionClock(req.Region)
	if err != nil {
		return nil, nil
	}
	// Devolver una copia del reloj actual
	return &pb.ClockResponse{
		Clock:    clock,
		Servidor: int32(ID_server),
	}, nil
}

// Obtiene el reloj de la region especifica
func getRegionClock(region string) ([]int32, error) {
	// Abrir el archivo regiones.txt
	file, err := os.Open("regiones.txt")
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	// Leer el archivo línea por línea
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Obtener la línea y dividirla en partes
		line := scanner.Text()
		parts := strings.Fields(line)

		// Si la región coincide con la que buscamos, devolver el reloj
		if parts[0] == region {
			// Convertir las partes de la línea en enteros y devolver el reloj
			var clock []int32
			for _, val := range parts[1:] {
				var v int32
				fmt.Sscanf(val, "%d", &v)
				clock = append(clock, v)
			}
			return clock, nil
		}
	}

	// Si no se encontró la región, devolver el reloj inicializado en 0, 0, 0
	return []int32{0, 0, 0}, nil
}

// appendToFile agrega texto al final de un archivo, creándolo si no existe
func appendToFile(filename, text string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	return err
}

// Pocesa los comandos de Logs.txt
func procesarLogsFromContent(content string) {
	// Dividir el contenido en líneas
	logLines := strings.Split(content, "\n")
	for _, line := range logLines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Printf("Línea inválida en el contenido mergeado: %s\n", line)
			continue
		}

		command := parts[0]
		region := parts[1]
		product := parts[2]
		value := strings.Join(parts[3:], " ") // En caso de valores con espacios

		regionFile := fmt.Sprintf("%s.txt", region)

		switch command {
		case "AgregarProducto":
			entry := fmt.Sprintf("%s %s %s\n", region, product, value)
			appendToFile(regionFile, entry)
		case "RenombrarProducto":
			renameProductInRegion(regionFile, product, value)
		case "ActualizarValor":
			updateProductValueInRegion(regionFile, product, value)
		case "BorrarProducto":
			removeProductFromRegion(regionFile, product)
		default:
			fmt.Printf("Comando desconocido en el contenido mergeado: %s\n", line)
		}
	}
}

// readFile lee el contenido de un archivo y lo devuelve como un string
func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content.String(), nil
}

// writeFile escribe un string en un archivo, sobrescribiéndolo
func writeFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}
	return writer.Flush()
}

// updateFileLines actualiza las líneas de un archivo según una función de transformación
func updateFileLines(filename string, transform func(string) string) {
	content, err := readFile(filename)
	if err != nil {
		fmt.Printf("Error leyendo %s: %v\n", filename, err)
		return
	}

	var output []string
	for _, line := range strings.Split(content, "\n") {
		fmt.Printf("Procesando línea: %s\n", line) // Depuración
		newLine := transform(line)
		if strings.TrimSpace(newLine) != "" {
			output = append(output, newLine)
		}
	}

	err = writeFile(filename, strings.Join(output, "\n"))
	if err != nil {
		fmt.Printf("Error escribiendo %s: %v\n", filename, err)
	}
}

// renameProductInRegion renombra un producto en el archivo de región
func renameProductInRegion(regionFile, oldName, newName string) {
	updateFileLines(regionFile, func(line string) string {
		if strings.Contains(line, oldName) {
			return strings.Replace(line, oldName, newName, 1)
		}
		return line
	})
}

// updateProductValueInRegion actualiza el valor de un producto en el archivo de región
func updateProductValueInRegion(regionFile, product, newValue string) {
	updateFileLines(regionFile, func(line string) string {
		// Verifica que la línea contiene exactamente el producto
		parts := strings.Fields(line)
		if len(parts) >= 4 && parts[1] == product { // Asegúrate de que "parts[1]" es el producto
			parts[2] = newValue // Actualiza el valor del producto
			return strings.Join(parts, " ")
		}
		return line // Deja la línea sin cambios si no corresponde
	})
}

// removeProductFromRegion elimina un producto del archivo de región
func removeProductFromRegion(regionFile, product string) {
	updateFileLines(regionFile, func(line string) string {
		if strings.Contains(line, product) {
			return ""
		}
		return line
	})
}

// ObtenerProducto maneja la consulta de un producto en una región específica
func (s *hextechServer) ObtenerProducto(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	regionFile := fmt.Sprintf("%s.txt", req.Region)

	data, err := readFile(regionFile)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de región %s: %v", req.Region, err)
	}

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.Contains(line, req.ProductName) {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return &pb.CommandResponse{
					Message:        "Producto encontrado",
					Value:          parts[2], // Asume que el valor está en la posición 2
					RelojVectorial: s.clock,
				}, nil
			}
		}
	}

	return &pb.CommandResponse{
		Message: "Producto no encontrado",
	}, nil
}

func main() {

	file, err := os.Create("Logs.txt")
	if err != nil {
		fmt.Printf("Error creando Logs.txt: %v\n", err)
		return
	}
	defer file.Close()

	listener, err := net.Listen("tcp", ":50054")
	if err != nil {
		fmt.Printf("Error iniciando el servidor: %v\n", err)
		return
	}
	fmt.Println("Servidor Hextech 3 escuchando en el puerto 50054")

	server := grpc.NewServer()

	//Inicializar reloj vectorial de server en [0,0,0]
	hextech := &hextechServer{
		clock: make([]int32, 3), // Inicializar reloj vectorial con 3 nodos
	}

	pb.RegisterBrokerServer(server, hextech)
	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		fmt.Printf("Error sirviendo: %v\n", err)
	}
}
