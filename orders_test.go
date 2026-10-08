package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderQueryAll(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rec := httptest.NewRecorder()

	orderQueryHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body OrderQueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Total != len(seedOrders) || len(body.Orders) != len(seedOrders) {
		t.Fatalf("body = %+v", body)
	}
}

func TestOrderQueryKeyword(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders?q=compiler", nil)
	rec := httptest.NewRecorder()

	orderQueryHandler(rec, req)

	var body OrderQueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Total != 1 || body.Orders[0].ID != "o3" {
		t.Fatalf("body = %+v", body)
	}
}

func TestOrderQueryFilters(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders?userId=u1&status=paid", nil)
	rec := httptest.NewRecorder()

	orderQueryHandler(rec, req)

	var body OrderQueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Total != 2 {
		t.Fatalf("body = %+v", body)
	}
	for _, order := range body.Orders {
		if order.UserID != "u1" || order.Status != "paid" {
			t.Fatalf("order = %+v", order)
		}
	}
}

func TestOrderQueryNoMatch(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders?q=nobody", nil)
	rec := httptest.NewRecorder()

	orderQueryHandler(rec, req)

	var body OrderQueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Total != 0 || len(body.Orders) != 0 {
		t.Fatalf("body = %+v", body)
	}
}

func TestOrderQueryHead(t *testing.T) {
	req := httptest.NewRequest(http.MethodHead, "/orders?status=pending", nil)
	rec := httptest.NewRecorder()

	orderQueryHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("head body = %q, want empty", rec.Body.String())
	}
}

func TestOrderQueryMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/orders", nil)
	rec := httptest.NewRecorder()

	orderQueryHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestOrderByID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/o2", nil)
	req.SetPathValue("id", "o2")
	rec := httptest.NewRecorder()

	orderByIDHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body Order
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.ID != "o2" || body.UserID != "u2" || body.Status != "pending" {
		t.Fatalf("body = %+v", body)
	}
}

func TestOrderByIDNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/missing", nil)
	req.SetPathValue("id", "missing")
	rec := httptest.NewRecorder()

	orderByIDHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestOrderRoutes(t *testing.T) {
	mux := newMux()

	listReq := httptest.NewRequest(http.MethodGet, "/orders?item=engine", nil)
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d", listRec.Code)
	}

	var list OrderQueryResponse
	if err := json.Unmarshal(listRec.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if list.Total != 2 {
		t.Fatalf("list = %+v", list)
	}

	oneReq := httptest.NewRequest(http.MethodGet, "/orders/o1", nil)
	oneRec := httptest.NewRecorder()
	mux.ServeHTTP(oneRec, oneReq)
	if oneRec.Code != http.StatusOK {
		t.Fatalf("get status = %d", oneRec.Code)
	}

	var body Order
	if err := json.Unmarshal(oneRec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Item != "Analytical Engine" {
		t.Fatalf("body = %+v", body)
	}
}
