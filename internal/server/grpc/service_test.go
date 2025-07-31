package grpc

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/dimryb/cross-arb/internal/types"
	"github.com/dimryb/cross-arb/mocks"
	"github.com/dimryb/cross-arb/proto"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	protocmp "google.golang.org/protobuf/proto"
)

func TestTickerService_Subscribe_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockTickerStore(ctrl)
	mockSub := mocks.NewMockTickerSubscriber(ctrl)
	mockApp := mocks.NewMockApplication(ctrl)
	mockLog := mocks.NewMockLogger(ctrl)

	testTicker := types.TickerData{
		Symbol:   "BTC_USDT",
		Exchange: "mexc",
		BidPrice: 60000.0,
		BidQty:   0.5,
		AskPrice: 60100.0,
		AskQty:   0.4,
	}

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("1. Setting up expectations")

	mockApp.EXPECT().Context().Return(appCtx).AnyTimes()
	mockApp.EXPECT().Logger().Return(mockLog).AnyTimes()
	mockApp.EXPECT().TickerStore().Return(mockStore).AnyTimes()
	mockStore.EXPECT().AddSubscriber().Return(mockSub).Times(1)

	eventCh := make(chan types.TickerEvent, 1)

	mockSub.EXPECT().Recv().DoAndReturn(func() (types.TickerEvent, bool) {
		t.Log("7,9. Recv() called, waiting for event or context cancel...")
		select {
		case event, ok := <-eventCh:
			if ok {
				t.Log("8. Recv() got event:", event.Ticker.Symbol)
			} else {
				t.Log("5. Recv() channel closed")
			}
			return event, ok
		case <-appCtx.Done():
			t.Log("13. Recv() context cancelled, returning (zero, false)")
			var zero types.TickerEvent
			return zero, false
		}
	}).AnyTimes()

	mockSub.EXPECT().Done().Return(appCtx.Done()).AnyTimes()
	mockSub.EXPECT().Close().Do(func() {
		t.Log("14. Close() called on subscriber")
	}).Times(1)

	t.Log("2. Creating service")
	service := NewTickerService(mockApp)

	// Мок-стрим с каналом для ожидания отправки
	mockStream := &MockTickerServiceSubscribeServer{
		SentUpdates: make(chan *proto.TickerUpdate, 1), // буфер на 1 сообщение
		StreamCtx:   context.Background(),
	}

	t.Log("3. Starting Subscribe in goroutine")
	errChan := make(chan error, 1)
	go func() {
		t.Log("6. Subscribe() started")
		err := service.Subscribe(nil, mockStream)
		t.Log("15. Subscribe() returned with error:", err)
		errChan <- err
	}()

	t.Log("4. Sending event to eventCh")
	eventCh <- types.TickerEvent{Ticker: testTicker}

	t.Log("5. Waiting for update to be sent via channel...")
	var update *proto.TickerUpdate
	select {
	case update = <-mockStream.SentUpdates:
		t.Log("10. Event was successfully sent to stream")
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout waiting for update to be sent")
	}

	// Проверка содержимого
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
		t.Errorf("mismatched update\ngot:  %v\nwant: %v", update, want)
	}

	t.Log("11. Canceling context to trigger shutdown")
	cancel()

	t.Log("12. Waiting for Subscribe to return...")
	select {
	case err := <-errChan:
		t.Log("16. Received error from Subscribe():", err)
		if err != nil {
			code := status.Code(err)
			if code != codes.Canceled && code != codes.OK {
				t.Errorf("unexpected error: %v", err)
			}
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for Subscribe to return")
	}
}

// Мок gRPC-стрима
type MockTickerServiceSubscribeServer struct {
	SentUpdates chan *proto.TickerUpdate // ← теперь канал!
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

func (m *MockTickerServiceSubscribeServer) RecvMsg(msg interface{}) error {
	return io.EOF
}

func (m *MockTickerServiceSubscribeServer) Context() context.Context {
	return m.StreamCtx
}

func (m *MockTickerServiceSubscribeServer) SetHeader(md metadata.MD) error {
	return nil
}

func (m *MockTickerServiceSubscribeServer) SendHeader(md metadata.MD) error {
	return nil
}

func (m *MockTickerServiceSubscribeServer) SetTrailer(md metadata.MD) {
}
