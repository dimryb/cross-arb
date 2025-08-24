package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	grpcAdapter "github.com/dimryb/cross-arb/internal/adapter/grpc"
	"github.com/dimryb/cross-arb/internal/entity"
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
		book := grpcAdapter.ToBookTicker(update.GetData())
		result := entity.Result{
			Symbol: book.Symbol,
			Data:   book,
		}

		// Печатаем как одноэлементный отчёт
		PrintTickersReport([]entity.Result{result})
	}
}

func PrintTickersReport(results []entity.Result) {
	fmt.Printf("=== Обновление цен (%s) ===\n", time.Now().Format("15:04:05.000"))
	for _, r := range results {
		PrintTicker(r.Data)
	}
	fmt.Println()
}

func PrintTicker(t entity.BookTicker) {
	fmt.Printf(
		"  [%s] -> покупка: %s (%s) | продажа: %s (%s)\n",
		t.Symbol,
		t.BidPrice, t.BidQty,
		t.AskPrice, t.AskQty,
	)
}
