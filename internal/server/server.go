package server

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"net"
	"io"

	"github.com/cleamoon/gored-kv-database/pkg/kvstore"
)

type Server struct {
	store *kvstore.KVStore
	port string
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", ":" + s.port)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server is listening on port %s", s.port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading request: %v", err)
			}
			return
		}

		line = strings.TrimSpace(line)
		// parts[0] = command, parts[1] = key, parts[2] = value
		parts := strings.SplitN(line, " ", 3)

		if len(parts) < 1 {
			conn.Write([]byte("Invalid request\n"))
			writer.Flush()
			continue
		}

		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "SET":
			if len(parts) != 3 {
				conn.Write([]byte("Invalid request\nUsage: SET <key> <value>\n"))
				writer.Flush()
				continue
			}
			
			key := parts[1]
			value := parts[2]
			s.store.Set(key, value)
			conn.Write([]byte("OK\n"))
			continue
		case "GET":
			if len(parts) != 2 {
				conn.Write([]byte("Invalid request\nUsage: GET <key>\n"))
				continue
			}
			
			key := parts[1]
			value, exists := s.store.Get(key)
			if !exists {
				conn.Write([]byte("Key not found\n"))
				continue
			}
			
			conn.Write([]byte(value + "\n"))
			continue
		case "DEL":
			if len(parts) != 2 {
				conn.Write([]byte("Invalid request\nUsage: DEL <key>\n"))
				continue
			}
			
			key := parts[1]
			exists := s.store.Delete(key)
			if !exists {
				conn.Write([]byte("Key not found\n"))
				continue
			}
			
			conn.Write([]byte("OK\n"))
			continue
		default:
			conn.Write([]byte("Invalid command\nUsage: SET <key> <value>, GET <key>, DEL <key>\n"))
			continue
		}
	}
}

