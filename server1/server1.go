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
	"time"

	pb "hextech/grpc-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Servidor Hextech
type hextechServer struct {
	pb.UnimplementedBrokerServer
	clock        []int32    // Reloj vectorial
	mutex        sync.Mutex // Para sincronizar el acceso al reloj
	regionClocks map[string][]int32
}

var ID_server = 0

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

	fmt.Printf("Servidor Hextech 1: Comando registrado en Logs.txt: %s\n", logEntry)

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

// GetVectorClock devuelve un reloj vectorial [x, y, z]
func (s *hextechServer) GetClock(ctx context.Context, req *pb.ClockRequest) (*pb.ClockResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	clock, err := getRegionClock(req.Region)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el reloj de la región: %v", err)
	}
	// Devolver una copia del reloj actual
	return &pb.ClockResponse{
		Clock:    clock,
		Servidor: int32(ID_server),
	}, nil
}

// Obtiene un reloj de region especifica leyendo el archivo regiones que guarda los relojes
func getRegionClock(region string) ([]int32, error) {
	// Abrir el archivo regiones.txt
	file, err := os.Open("regiones.txt")
	if err != nil {
		return nil, fmt.Errorf("error abriendo el archivo regiones.txt: %v", err)
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

// Funcion para mandar los logs.txt actualizados a los demás servidores
func distributeMergedLogs(mergedContent string) {
	servers := []string{"server2:50053", "server3:50054"} // Direcciones de otros Fulcrum

	for _, server := range servers {
		conn, err := grpc.Dial(server, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("Servidor Hextech 1: Error al conectar con el servidor %s: %v\n", server, err)
			continue
		}
		defer conn.Close()

		client := pb.NewBrokerClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		req := &pb.FileRequest{
			Filename: "MergedLogs",
			Content:  []byte(mergedContent), // Convertir el contenido a bytes
		}

		_, err = client.ReceiveMergedFile(ctx, req)
		if err != nil {
			log.Printf("Servidor Hextech 1: Error al enviar el archivo mergeado al servidor %s: %v\n", server, err)
		} else {
			log.Printf("Servidor Hextech 1: Archivo mergeado enviado exitosamente al servidor %s", server)
		}
	}
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

// Obtiene los archivos logs.txt de cada servidor para realizar posteriormente el merge
func getFileFromServer(address, filename string) (string, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return "", nil
	}
	defer conn.Close()

	client := pb.NewBrokerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.FileRequest{Filename: filename}
	res, err := client.GetFile(ctx, req)
	if err != nil {
		return "", nil
	}

	// Retornar el contenido como string
	return res.Content, nil
}

// funcion auxiliar para realizar conexiones con los servidores
func distributeMergedFile(content, filename string) {
	servers := []string{"server2:50053", "server3:50054"}
	for _, server := range servers {
		err := sendFileToServer(server, content)
		if err != nil {
			log.Printf("Error al enviar %s a %s: %v\n", filename, server, err)
		}
	}
}

// Envia los archivos de logs correspondientes al merge a los servidores
func sendFileToServer(server, content string) error {
	conn, err := grpc.Dial(server, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil
	}
	defer conn.Close()

	client := pb.NewBrokerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.FileRequest{
		Content: []byte(content), // Enviar el contenido como bytes
	}
	_, err = client.ReceiveMergedFile2(ctx, req)
	if err != nil {
		return nil
	}

	return nil
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

// Actualiza los relojes de las regiones correspondientes al merge
func mergeClocks(clock1, clock2 []int32) []int32 {
	maxLen := max(len(clock1), len(clock2))
	merged := make([]int32, maxLen)

	for i := 0; i < maxLen; i++ {
		val1 := int32(0)
		if i < len(clock1) {
			val1 = clock1[i]
		}

		val2 := int32(0)
		if i < len(clock2) {
			val2 = clock2[i]
		}

		merged[i] = maxInt32(val1, val2)
	}

	return merged
}

// funcion auxiliar para determinar que variable es mayor XD
func maxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// Procesa los logs de cada servidor, funcion main para realizar el merge
func procesarLogs() {
	for {
		time.Sleep(30 * time.Second)
		servers := []string{"server2:50053", "server3:50054"} // Direcciones de otros Fulcrum
		filename := "Logs.txt"
		filename2 := "regiones.txt"
		var allLines [][]string
		var allRegions [][]string
		mergedClocks := make(map[string][]int32)

		// Extraer la información
		for _, server := range servers {
			content, _ := getFileFromServer(server, filename)
			content2, err := getFileFromServer(server, filename2)
			if err != nil {
				log.Printf("Servidor Hextech 1: Error al obtener el archivo de log desde %s: %v\n", server, err)
				continue
			}
			if content != " " {
				lines := strings.Split(content, ", ")
				allLines = append(allLines, lines)
			}
			if content2 != " " {
				lines2 := strings.Split(content2, ", ")
				allRegions = append(allRegions, lines2)
			} else {
				// Crear archivo vacío si no existe
				err := func() error {
					file, err := os.Create(filename) // Crear el archivo
					if err != nil {
						return err
					}
					defer file.Close()
					writer := bufio.NewWriter(file) // Inicializar un escritor bufferizado
					_, err = writer.WriteString("") // Escribir una cadena vacía
					if err != nil {
						return err
					}
					return writer.Flush() // Asegurar que todo se escriba al archivo
				}()
				if err != nil {
					log.Printf("Servidor Hextech 1: Error al crear archivo vacío en %s: %v\n", server, err)
					continue
				}
			}
		}

		logData, err := readFile("Logs.txt")
		if err != nil {
			fmt.Printf("Servidor Hextech 1: Error leyendo Logs.txt: %v\n", err)
			continue
		}

		// Leer el archivo local `regiones.txt`
		localData, err := readFile(filename2)
		if err == nil {
			lines := strings.Split(localData, "\n")
			allRegions = append(allRegions, lines)
		}

		// Fusionar los relojes por región
		for _, regionLines := range allRegions {
			for _, line := range regionLines {
				if strings.TrimSpace(line) == "" {
					continue
				}

				parts := strings.Fields(line)
				if len(parts) < 2 {
					log.Printf("Línea inválida en regiones.txt: %s\n", line)
					continue
				}

				region := parts[0]
				clock := parseClock(parts[1:])

				// Fusionar el reloj para la región
				if localClock, exists := mergedClocks[region]; exists {
					mergedClocks[region] = mergeClocks(localClock, clock)
				} else {
					mergedClocks[region] = clock
				}
			}
		}

		// Escribir el archivo mergeado `regiones.txt`
		var mergedLines []string
		for region, clock := range mergedClocks {
			line := fmt.Sprintf("%s %s", region, formatClock(clock))
			mergedLines = append(mergedLines, line)
		}

		err = writeFile(filename2, strings.Join(mergedLines, "\n"))
		if err != nil {
			log.Printf("Error escribiendo %s: %v\n", filename2, err)
		}

		// Distribuir el archivo mergeado a los otros servidores
		distributeMergedFile(strings.Join(mergedLines, "\n"), filename2)

		lines := strings.Split(logData, "\n")
		allLines = append(allLines, lines)

		// Juntar toda la información
		var merged []string
		processedCommands := make(map[string]bool)
		for _, array := range allLines {
			for _, line := range array {
				if !processedCommands[line] {
					processedCommands[line] = true
					merged = append(merged, line)
				}
			}
		}

		// Distribuir el archivo mergeado a los otros servidores
		distributeMergedLogs(strings.Join(merged, "\n"))

		// Procesar cada línea de `merged`
		for _, line := range merged {
			if strings.TrimSpace(line) == "" {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) < 3 {
				fmt.Printf("Servidor Hextech 1: Línea inválida en los logs mergeados: %s\n", line)
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
				fmt.Printf("Servidor Hextech 1: Comando desconocido en los logs mergeados: %s\n", line)
			}
		}
		// Vaciar el log
		err = writeFile("Logs.txt", "")
		if err != nil {
			fmt.Printf("Servidor Hextech 1: Error limpiando Logs.txt: %v\n", err)
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
		if len(parts) >= 3 && parts[1] == product { // Asegúrate de que "parts[1]" es el producto
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
		return nil, nil
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

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		fmt.Printf("Servidor Hextech 1: Error iniciando el servidor: %v\n", err)
		return
	}
	fmt.Println("Servidor Hextech 1 escuchando en el puerto 50052")

	// Inicializar mapa de relojes vectoriales por región
	regionClocks := make(map[string][]int32)

	// Inicializar servidor GRPC
	server := grpc.NewServer()
	hextech := &hextechServer{
		regionClocks: regionClocks, // Mapa de relojes
	}
	pb.RegisterBrokerServer(server, hextech)
	reflection.Register(server)

	// Iniciar la rutina para procesar logs
	go procesarLogs()

	if err := server.Serve(listener); err != nil {
		fmt.Printf("Error sirviendo: %v\n", err)
	}
}
