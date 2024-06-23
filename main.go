package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var ClusterHolder []ClusterMemoryDB

var DefaultCacheDuration = 30 * time.Second

func main() {
	listener, err := net.Listen("tcp", ":57")
	if err != nil {
		log.Fatalf("Error listening: %s", err)
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Error closing listener: %s", err)
		}
	}()

	fmt.Println("Server started on tcp port :57")
	clientAccept(listener)
}

func clientAccept(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}

		connectionAddress := connection.RemoteAddr().String()

		fmt.Println("Client connected from:", connectionAddress)

		go handleClient(connection, ClusterHolder)
	}
}

func handleClient(connection net.Conn, cluster []ClusterMemoryDB) {
	defer func() {
		err := connection.Close()
		if err != nil {
			log.Printf("Error closing connection: %s", err)
		}
	}()

	scanner := bufio.NewScanner(connection)

	var cache *KeyValueCache

	byteValue, err := os.ReadFile("guide.txt")
	_, err = connection.Write(byteValue)

	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
	if err != nil {
		fmt.Printf("Error writing to client: %v\n", err)
	}
	fmt.Println("\n")
	for scanner.Scan() {
		request := scanner.Text()
		fmt.Printf("Received request : %v from Client Address: %v\n", request, connection.RemoteAddr().String())

		parts := strings.Fields(request)
		if len(parts) < 2 {
			_, err := connection.Write([]byte("Invalid command format\n"))
			if err != nil {
				log.Printf("Error writing response: %s", err)
			}
			continue
		}

		command := parts[0]
		key := parts[1]

		connectionInt, _ := strconv.Atoi(connection.RemoteAddr().String())

		switch command {
		case "Config":
			if len(parts) < 2 {
				_, err := connection.Write([]byte("Usage: Config <duration_seconds>\n"))
				if err != nil {
					log.Printf("Error writing response: %s", err)
				}
				continue
			}

			timeConfig, err := strconv.Atoi(key)
			if err != nil {
				log.Printf("Error parsing duration: %s", err)
				_, err := connection.Write([]byte("Invalid duration format\n"))
				if err != nil {
					log.Printf("Error writing response: %s", err)
				}
				continue
			}

			duration := time.Duration(timeConfig) * time.Second
			Cluster := NewCluster(connectionInt, NewCache(3, duration))
			ClusterHolder = append(ClusterHolder, *Cluster)
			cache = findFromCluster(ClusterHolder, connection.RemoteAddr().String())

		case "GET":
			if cache == nil {
				_, err := connection.Write([]byte("Cache not configured. Please use 'Config' command first.\n"))
				if err != nil {
					log.Printf("Error writing response: %s", err)
				}
				continue
			}
			cache = findFromCluster(ClusterHolder, connection.RemoteAddr().String())
			value, found := cache.get(key)
			if found {
				_, err := connection.Write([]byte(fmt.Sprintf("Value for %s: %s\n", key, value)))
				if err != nil {
					log.Printf("Error writing response: %s", err)
				}
			} else {
				_, err := connection.Write([]byte(fmt.Sprintf("Response with %d means Key must be deleted with the expiration Machinism or haven't set \n", -1)))
				if err != nil {
					log.Printf("Error writing response: %s", err)
				}
			}

		case "SET":
			if cache == nil {
				_, err := connection.Write([]byte("Cache not configured. Please use 'Config' command first.\n"))
				if err != nil {
					log.Printf("Error writing response: %s", err)
				}
				continue
			}

			if len(parts) < 3 {
				_, err := connection.Write([]byte("Usage: SET <key> <value>\n"))
				if err != nil {
					log.Printf("Error writing response: %s", err)
				}
				continue
			}
			value := parts[2]
			cache.put(key, value)
			_, err := connection.Write([]byte(fmt.Sprintf("Set %s = %s\n", key, value)))
			if err != nil {
				log.Printf("Error writing response: %s", err)
			}

		default:
			_, err := connection.Write([]byte("Invalid Command\n"))
			if err != nil {
				log.Printf("Error writing response: %s", err)
			}
		}
	}
}

func findFromCluster(cluster []ClusterMemoryDB, address string) *KeyValueCache {
	connectionAddInt, _ := strconv.Atoi(address)
	clusterSort(cluster, 0, len(cluster)-1)
	return BinarySearch(cluster, connectionAddInt)
}

func BinarySearch(cluster []ClusterMemoryDB, ID int) *KeyValueCache {
	low := 0
	high := len(cluster) - 1

	for low <= high {
		mid := (low + high) / 2
		if cluster[mid].ClientID == ID {
			return cluster[mid].KeyValueCache
		} else if cluster[mid].ClientID < ID {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return &KeyValueCache{}
}
