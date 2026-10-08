package main

import (
	"net/http"
	"strings"
)

// Order 是可查询的订单记录。Amount 是订单金额，单位为分。
type Order struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Item   string `json:"item"`
	Status string `json:"status"`
	Amount int    `json:"amount"`
}

// OrderQueryResponse 是 /orders 的查询结果。
type OrderQueryResponse struct {
	Total  int     `json:"total"`
	Orders []Order `json:"orders"`
}

// seedOrders 是进程内示例订单，服务重启后恢复为这组记录。
var seedOrders = []Order{
	{ID: "o1", UserID: "u1", Item: "Analytical Engine", Status: "paid", Amount: 12800},
	{ID: "o2", UserID: "u2", Item: "Enigma notes", Status: "pending", Amount: 4500},
	{ID: "o3", UserID: "u3", Item: "Compiler manual", Status: "shipped", Amount: 3200},
	{ID: "o4", UserID: "u1", Item: "Engine notebook", Status: "paid", Amount: 2600},
}

func orderQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/orders" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	orders := queryOrders(query.Get("q"), query.Get("id"), query.Get("userId"), query.Get("item"), query.Get("status"))
	writeJSON(w, r, OrderQueryResponse{Total: len(orders), Orders: orders})
}

func orderByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := orderIDFromRequest(r)
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	order, ok := findOrderByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, r, order)
}

func queryOrders(keyword, id, userID, item, status string) []Order {
	keyword = strings.TrimSpace(keyword)
	id = strings.TrimSpace(id)
	userID = strings.TrimSpace(userID)
	item = strings.TrimSpace(item)
	status = strings.TrimSpace(status)

	matched := make([]Order, 0)
	for _, order := range seedOrders {
		if id != "" && !strings.EqualFold(order.ID, id) {
			continue
		}
		if userID != "" && !strings.EqualFold(order.UserID, userID) {
			continue
		}
		if item != "" && !containsFold(order.Item, item) {
			continue
		}
		if status != "" && !strings.EqualFold(order.Status, status) {
			continue
		}
		if keyword != "" && !orderMatchesKeyword(order, keyword) {
			continue
		}
		matched = append(matched, order)
	}
	return matched
}

func findOrderByID(id string) (Order, bool) {
	for _, order := range seedOrders {
		if strings.EqualFold(order.ID, id) {
			return order, true
		}
	}
	return Order{}, false
}

func orderMatchesKeyword(order Order, keyword string) bool {
	return containsFold(order.ID, keyword) ||
		containsFold(order.UserID, keyword) ||
		containsFold(order.Item, keyword) ||
		containsFold(order.Status, keyword)
}

func orderIDFromRequest(r *http.Request) string {
	if id := strings.TrimSpace(r.PathValue("id")); id != "" {
		return id
	}
	const prefix = "/orders/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return ""
	}
	return strings.Trim(strings.TrimPrefix(r.URL.Path, prefix), "/")
}
