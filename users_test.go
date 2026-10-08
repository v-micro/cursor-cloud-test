package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserQueryAll(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	userQueryHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body UserQueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Total != len(seedUsers) || len(body.Users) != len(seedUsers) {
		t.Fatalf("body = %+v", body)
	}
}

func TestUserQueryKeyword(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users?q=hopper", nil)
	rec := httptest.NewRecorder()

	userQueryHandler(rec, req)

	var body UserQueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Total != 1 || body.Users[0].ID != "u3" {
		t.Fatalf("body = %+v", body)
	}
}

func TestUserQueryFilters(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users?name=alan&email=example.com", nil)
	rec := httptest.NewRecorder()

	userQueryHandler(rec, req)

	var body UserQueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Total != 1 || body.Users[0].Email != "alan@example.com" {
		t.Fatalf("body = %+v", body)
	}
}

func TestUserQueryNoMatch(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users?q=nobody", nil)
	rec := httptest.NewRecorder()

	userQueryHandler(rec, req)

	var body UserQueryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Total != 0 || len(body.Users) != 0 {
		t.Fatalf("body = %+v", body)
	}
}

func TestUserQueryHead(t *testing.T) {
	req := httptest.NewRequest(http.MethodHead, "/users?q=ada", nil)
	rec := httptest.NewRecorder()

	userQueryHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("head body = %q, want empty", rec.Body.String())
	}
}

func TestUserQueryMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rec := httptest.NewRecorder()

	userQueryHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestUserByID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/u2", nil)
	req.SetPathValue("id", "u2")
	rec := httptest.NewRecorder()

	userByIDHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body User
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.ID != "u2" || body.Name != "Alan Turing" {
		t.Fatalf("body = %+v", body)
	}
}

func TestUserByIDNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/missing", nil)
	req.SetPathValue("id", "missing")
	rec := httptest.NewRecorder()

	userByIDHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestUserRoutes(t *testing.T) {
	mux := newMux()

	listReq := httptest.NewRequest(http.MethodGet, "/users?id=u1", nil)
	listRec := httptest.NewRecorder()
	mux.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d", listRec.Code)
	}

	oneReq := httptest.NewRequest(http.MethodGet, "/users/u1", nil)
	oneRec := httptest.NewRecorder()
	mux.ServeHTTP(oneRec, oneReq)
	if oneRec.Code != http.StatusOK {
		t.Fatalf("get status = %d", oneRec.Code)
	}

	var body User
	if err := json.Unmarshal(oneRec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Email != "ada@example.com" {
		t.Fatalf("body = %+v", body)
	}
}
