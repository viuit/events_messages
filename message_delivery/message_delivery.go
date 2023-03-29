package main

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// func main() {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
// 	}

// 	defer func() {
// 		_ = conn.Close() // Закрываем подключение в случае удачной попытки
// 	}()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Fatalf("failed to open channel. Error: %s", err)
// 	}

// 	defer func() {
// 		_ = ch.Close() // Закрываем канал в случае удачной попытки открытия
// 	}()
// 	host, exists := os.LookupEnv("HOST")

// 	if exists {
// 		fmt.Println(host)
// 	}

// 	fmt.Println("hello")
// }
