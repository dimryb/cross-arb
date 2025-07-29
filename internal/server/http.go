package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dimryb/cross-arb/internal/service"
)

type HTTPServer struct {
	store *service.TickerStore
}

func NewHTTPServer(store *service.TickerStore) *HTTPServer {
	return &HTTPServer{store: store}
}

func (s *HTTPServer) Run(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/tickers", s.handleTickers)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server.ListenAndServe()
}

func (s *HTTPServer) handleTickers(w http.ResponseWriter, _ *http.Request) {
	tickers := s.store.GetAll()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tickers); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
