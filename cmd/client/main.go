package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/dimryb/cross-arb/internal/report"
	"github.com/dimryb/cross-arb/internal/types"
	"github.com/dimryb/cross-arb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverAddr = flag.String("addr", "localhost:9090", "gRPC server address")

func main() {
	flag.Parse()

	// Подключаемся к gRPC-серверу
	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Создаём клиент
	client := proto.NewTickerServiceClient(conn)

	// Вызываем стрим
	stream, err := client.Subscribe(context.Background(), &proto.SubscribeRequest{})
	if err != nil {
		log.Fatalf("Subscribe failed: %v", err) //nolint:gocritic
	}

	fmt.Printf(" Connected to %s\n", *serverAddr)
	fmt.Println(" Waiting for ticker updates...")

	// Читаем обновления
	for {
		update, err := stream.Recv()
		if err != nil {
			log.Printf(" Стрим завершён: %v", err)
			return
		}

		// Конвертируем proto → BookTicker → Result
		book := types.ToBookTicker(update.GetData())
		result := types.Result{
			Symbol: book.Symbol,
			Data:   book,
		}

		// Печатаем как одноэлементный отчёт
		report.PrintTickersReport([]types.Result{result})
	}
}
