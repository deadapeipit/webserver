package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"webserver/entity"

	"github.com/gorilla/mux"
)

type OrderHandlerInterface interface {
	OrdersHandler(w http.ResponseWriter, r *http.Request)
}

type OrderHandler struct {
	//postgrespool *pgxpool.Pool
}

func NewOrderHandler() OrderHandlerInterface {
	//return &UserHandler{postgrespool: postgrespool}
	return &OrderHandler{}
}

func (h *OrderHandler) OrdersHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	switch r.Method {
	case http.MethodGet:
		if id != "" { // get by id
			getOrdersByIDHandler(w, r, id)
		} else { // get all
			h.getOrdersHandler(w, r)
		}
	case http.MethodPost:
		if id != "" {
			updateOrderHandler(w, r, id)
		} else {
			createOrdersHandler(w, r)
		}
	case http.MethodPut:
		updateOrderHandler(w, r, id)
	case http.MethodDelete:
		deleteOrderHandler(w, r, id)
	}
}

func (h *OrderHandler) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	orders, err := SqlConnect.GetOrders(ctx)
	if err != nil {
		writeJsonResp(w, statusError, err.Error())
		return
	}
	writeJsonResp(w, statusSuccess, orders)
}

func getOrdersByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if idInt, err := strconv.Atoi(id); err == nil {
		orders, err := SqlConnect.GetOrderByID(ctx, idInt)
		if err != nil {
			writeJsonResp(w, statusError, err.Error())
			return
		}
		if idInt != orders.OrderId {
			writeJsonResp(w, statusError, "Data not exists")
			return
		}
		writeJsonResp(w, statusSuccess, orders)
	}
}

func createOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	decoder := json.NewDecoder(r.Body)
	var order entity.OrderWithItems
	if err := decoder.Decode(&order); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}

	orders, err := SqlConnect.CreateOrder(ctx, order)
	if err != nil {
		writeJsonResp(w, statusError, err.Error())
		return
	}
	writeJsonResp(w, statusSuccess, orders)
}

func updateOrderHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()

	if id != "" { // get by id
		if idInt, err := strconv.Atoi(id); err == nil {
			if orders, err := SqlConnect.GetOrderByID(ctx, idInt); err != nil {
				writeJsonResp(w, statusError, err.Error())
				return
			} else if idInt != orders.OrderId {
				writeJsonResp(w, statusError, "Data not exists")
				return
			} else {
				decoder := json.NewDecoder(r.Body)
				var order entity.OrderWithItems
				if err := decoder.Decode(&order); err != nil {
					w.Write([]byte("error decoding json body"))
					return
				}

				orders, err := SqlConnect.UpdateOrder(ctx, idInt, order)
				if err != nil {
					writeJsonResp(w, statusError, err.Error())
					return
				}
				writeJsonResp(w, statusSuccess, orders)
			}
		}
	}
}

func deleteOrderHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if id != "" { // get by id
		if idInt, err := strconv.Atoi(id); err == nil {
			if orders, err := SqlConnect.GetOrderByID(ctx, idInt); err != nil {
				writeJsonResp(w, statusError, err.Error())
				return
			} else if idInt != orders.OrderId {
				writeJsonResp(w, statusError, "Data not exists")
				return
			} else {
				order, err := SqlConnect.DeleteOrder(ctx, idInt)
				if err != nil {
					writeJsonResp(w, statusError, err.Error())
					return
				}
				writeJsonResp(w, statusSuccess, order)
			}
		}
	}
}
