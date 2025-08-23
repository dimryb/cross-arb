package grpc

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/dimryb/cross-arb/internal/entity"
	"github.com/dimryb/cross-arb/mocks"
	"github.com/dimryb/cross-arb/proto"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	protocmp "google.golang.org/protobuf/proto"
)

func TestTickerService_Subscribe_Success(t *testing.T) {
	// === 1. ARRANGE: Подготовка окружения ===
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Моки
	mockLog := mocks.NewMockLogger(ctrl)
	mockStore := mocks.NewMockTickerStore(ctrl)
	mockSub := mocks.NewMockTickerSubscriber(ctrl)
	mockApp := mocks.NewMockApplication(ctrl)

	// Тестовые данные
	testTicker := entity.TickerData{
		Symbol:   "BTC_USDT",
		Exchange: "mexc",
		BidPrice: 60000.0,
		BidQty:   0.5,
		AskPrice: 60100.0,
		AskQty:   0.4,
	}

	// Контекст с возможностью отмены
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем ожидания
	mockApp.EXPECT().Context().Return(appCtx).AnyTimes()
	mockApp.EXPECT().Logger().Return(mockLog).AnyTimes()
	mockApp.EXPECT().TickerStore().Return(mockStore).AnyTimes()
	mockStore.EXPECT().AddSubscriber().Return(mockSub).Times(1)

	// Канал для эмуляции потока событий
	eventCh := make(chan entity.TickerEvent, 1)

	// Мок Recv(): возвращает событие, затем реагирует на отмену
	mockSub.EXPECT().Recv().DoAndReturn(func() (entity.TickerEvent, bool) {
		select {
		case event, ok := <-eventCh:
			if ok {
				t.Log("Received ticker event:", event.Ticker.Symbol)
			}
			return event, ok
		case <-appCtx.Done():
			t.Log("Subscriber context cancelled, exiting Recv")
			var zero entity.TickerEvent
			return zero, false
		}
	}).AnyTimes()

	mockSub.EXPECT().Done().Return(appCtx.Done()).AnyTimes()
	mockSub.EXPECT().Close().Do(func() {
		t.Log("Subscriber.Close() called")
	}).Times(1)

	// Создаём сервис
	service := NewTickerService(mockApp)

	// Мок gRPC-стрима
	mockStream := &MockTickerServiceSubscribeServer{
		SentUpdates: make(chan *proto.TickerUpdate, 1),
		StreamCtx:   context.Background(),
	}

	// Канал для получения результата Subscribe
	errChan := make(chan error, 1)

	// === 2. ACT: Запуск и взаимодействие ===
	t.Log("Starting Subscribe in goroutine")
	go func() {
		err := service.Subscribe(nil, mockStream)
		errChan <- err
	}()

	t.Log("Publishing ticker event")
	eventCh <- entity.TickerEvent{Ticker: testTicker}

	t.Log("Waiting for gRPC update to be sent...")
	var update *proto.TickerUpdate
	select {
	case update = <-mockStream.SentUpdates:
		t.Log("Successfully sent update to client")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for update to be sent")
	}

	t.Log("Shutting down by cancelling context")
	cancel()

	t.Log("Waiting for Subscribe to complete...")
	var finalErr error
	select {
	case finalErr = <-errChan:
		t.Log("Subscribe() returned:", finalErr)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for Subscribe to return")
	}

	// === 3. ASSERT: Проверка результатов ===
	t.Log("Running assertions...")

	want := &proto.TickerUpdate{
		Data: &proto.TickerData{
			Symbol:   testTicker.Symbol,
			Exchange: testTicker.Exchange,
			BidPrice: testTicker.BidPrice,
			AskPrice: testTicker.AskPrice,
			BidQty:   testTicker.BidQty,
			AskQty:   testTicker.AskQty,
		},
	}

	if !protocmp.Equal(update, want) {
		t.Errorf("Update mismatch\ngot:  %v\nwant: %v", update, want)
	}

	if finalErr != nil {
		code := status.Code(finalErr)
		if code != codes.Canceled && code != codes.OK {
			t.Errorf("Unexpected error: %v", finalErr)
		}
	}
}

// Мок gRPC-стрима.
type MockTickerServiceSubscribeServer struct {
	SentUpdates chan *proto.TickerUpdate
	StreamCtx   context.Context
}

func (m *MockTickerServiceSubscribeServer) Send(update *proto.TickerUpdate) error {
	m.SentUpdates <- update
	return nil
}

func (m *MockTickerServiceSubscribeServer) SendMsg(msg interface{}) error {
	if update, ok := msg.(*proto.TickerUpdate); ok {
		return m.Send(update)
	}
	return fmt.Errorf("unexpected message type: %T", msg)
}

func (m *MockTickerServiceSubscribeServer) RecvMsg(_ interface{}) error {
	return io.EOF
}

func (m *MockTickerServiceSubscribeServer) Context() context.Context {
	return m.StreamCtx
}

func (m *MockTickerServiceSubscribeServer) SetHeader(_ metadata.MD) error {
	return nil
}

func (m *MockTickerServiceSubscribeServer) SendHeader(_ metadata.MD) error {
	return nil
}

func (m *MockTickerServiceSubscribeServer) SetTrailer(_ metadata.MD) {
}
