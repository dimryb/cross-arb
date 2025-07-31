package grpc

import (
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TickerService — реализация gRPC-сервиса.
type TickerService struct {
	proto.UnimplementedTickerServiceServer
	app i.Application
}

// NewTickerService — создаём сервис, внедряя Application.
func NewTickerService(app i.Application) *TickerService {
	return &TickerService{app: app}
}

// Subscribe — серверный стрим.
func (s *TickerService) Subscribe(_ *proto.SubscribeRequest, stream proto.TickerService_SubscribeServer) error {
	// Получаем хранилище через интерфейс
	store := s.app.TickerStore()

	// Подписываемся — получаем интерфейс
	subscriber := store.AddSubscriber()
	defer subscriber.Close()

	for {
		event, ok := subscriber.Recv()
		if !ok {
			return status.Error(codes.Canceled, "subscriber closed")
		}

		// Отправляем в gRPC-стрим
		if err := stream.Send(&proto.TickerUpdate{
			Data: &proto.TickerData{
				Symbol:   event.Ticker.Symbol,
				Exchange: event.Ticker.Exchange,
				BidPrice: event.Ticker.BidPrice,
				AskPrice: event.Ticker.AskPrice,
				BidQty:   event.Ticker.BidQty,
				AskQty:   event.Ticker.AskQty,
			},
		}); err != nil {
			return err // stream.Send может вернуть io.EOF или context.Canceled
		}
	}
}
