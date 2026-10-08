package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// User 是可查询的用户记录。
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserQueryResponse 是 /users 的查询结果。
type UserQueryResponse struct {
	Total int    `json:"total"`
	Users []User `json:"users"`
}

// seedUsers 是进程内示例数据，服务重启后恢复为这组记录。
var seedUsers = []User{
	{ID: "u1", Name: "Ada Lovelace", Email: "ada@example.com"},
	{ID: "u2", Name: "Alan Turing", Email: "alan@example.com"},
	{ID: "u3", Name: "Grace Hopper", Email: "grace@example.com"},
}

func userQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/users" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users := queryUsers(r.URL.Query().Get("q"), r.URL.Query().Get("id"), r.URL.Query().Get("name"), r.URL.Query().Get("email"))
	writeJSON(w, r, UserQueryResponse{Total: len(users), Users: users})
}

func userByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := userIDFromRequest(r)
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	user, ok := findUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, r, user)
}

func queryUsers(keyword, id, name, email string) []User {
	keyword = strings.TrimSpace(keyword)
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	matched := make([]User, 0)
	for _, user := range seedUsers {
		if id != "" && !strings.EqualFold(user.ID, id) {
			continue
		}
		if name != "" && !containsFold(user.Name, name) {
			continue
		}
		if email != "" && !containsFold(user.Email, email) {
			continue
		}
		if keyword != "" && !userMatchesKeyword(user, keyword) {
			continue
		}
		matched = append(matched, user)
	}
	return matched
}

func findUserByID(id string) (User, bool) {
	for _, user := range seedUsers {
		if strings.EqualFold(user.ID, id) {
			return user, true
		}
	}
	return User{}, false
}

func userMatchesKeyword(user User, keyword string) bool {
	return containsFold(user.ID, keyword) || containsFold(user.Name, keyword) || containsFold(user.Email, keyword)
}

func containsFold(value, substr string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(substr))
}

func userIDFromRequest(r *http.Request) string {
	if id := strings.TrimSpace(r.PathValue("id")); id != "" {
		return id
	}
	const prefix = "/users/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return ""
	}
	return strings.Trim(strings.TrimPrefix(r.URL.Path, prefix), "/")
}

func writeJSON(w http.ResponseWriter, r *http.Request, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
