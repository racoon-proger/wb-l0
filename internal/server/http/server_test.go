package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/racoon-proger/wb-l0/internal/domain"
)

func TestGetOrder_Ok(t *testing.T) {
	svc := NewMockservice(gomock.NewController(t))
	order := domain.Order{
		ID:          1,
		TrackNumber: "10",
	}
	svc.EXPECT().GetOrderByID(gomock.Any(), 1).Return(&order, nil)

	server := NewServer(svc, 8080)

	req, err := http.NewRequest(http.MethodGet, "/get-order", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "1")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetOrder)
	handler.ServeHTTP(rr, req)

	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, rr.Code)

}
