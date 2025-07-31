package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dimryb/cross-arb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:9090", "gRPC server address")
)

func main() {
	flag.Parse()

	// Подключаемся к gRPC-серверу
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Создаём клиент
	client := proto.NewTickerServiceClient(conn)

	// Вызываем стрим
	stream, err := client.Subscribe(context.Background(), &proto.SubscribeRequest{})
	if err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	fmt.Printf(" Connected to %s\n", *serverAddr)
	fmt.Println(" Waiting for ticker updates...")

	// Читаем обновления
	for {
		update, err := stream.Recv()
		if err != nil {
			log.Printf("Stream ended: %v", err)
			break
		}

		ticker := update.GetData()
		fmt.Printf("[%s] %s @ %s\n",
			time.Now().Format("15:04:05"),
			ticker.GetSymbol(),
			ticker.GetExchange(),
		)
		fmt.Printf("    Bid: %.2f (%.4f)\n", ticker.GetBidPrice(), ticker.GetBidQty())
		fmt.Printf("    Ask: %.2f (%.4f)\n", ticker.GetAskPrice(), ticker.GetAskQty())
		fmt.Println("    " + "-")
	}
}
