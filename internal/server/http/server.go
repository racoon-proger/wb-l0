//go:generate ~/go/bin/mockgen --build_flags=--mod=mod -destination=./server_mock.go -package=http -source=server.go
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/racoon-proger/wb-l0/internal/domain"
)

type service interface {
	GetOrderByID(ctx context.Context, id int) (order *domain.Order, err error)
}

type server struct {
	service service
	port    int
}

func (s *server) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *server) ServeHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	path := filepath.Join("public", "html", "index.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (s *server) GetOrder(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	order, err := s.service.GetOrderByID(r.Context(), id)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = json.NewEncoder(w).Encode(&order)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}

func NewServer(service service, port int) *server {
	return &server{
		service: service,
		port:    port,
	}
}
