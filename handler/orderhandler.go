package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"webserver/entity"

	"github.com/gorilla/mux"
)

type OrderHandler struct{}

func InstallOrderAPI(r *mux.Router) {
	api := OrderHandler{}
	r.HandleFunc("/orders", api.OrdersHandler)
	r.HandleFunc("/orders/{id}", api.OrdersHandler)
}

type OrderHandlerInterface interface {
	OrdersHandler(w http.ResponseWriter, r *http.Request)
}

func NewOrderHandler() OrderHandlerInterface {
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

// getOrdersHandler
// Method: GET
// Example: localhost/orders
func (h *OrderHandler) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	orders, err := Helper.Tesdb.GetOrders(ctx)
	if err != nil {
		WriteJsonResp(w, StatusError, err.Error())
		return
	}
	WriteJsonResp(w, StatusSuccess, orders)
}

// getOrdersByIDHandler
// Method: GET
// Example: localhost/orders/1
func getOrdersByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if idInt, err := strconv.Atoi(id); err == nil {
		orders, err := Helper.Tesdb.GetOrderByID(ctx, idInt)
		if err != nil {
			WriteJsonResp(w, StatusError, err.Error())
			return
		}
		if idInt != orders.OrderId {
			WriteJsonResp(w, StatusError, "Data not exists")
			return
		}
		WriteJsonResp(w, StatusSuccess, orders)
	}
}

// createOrdersHandler
// Method: POST
// Example: localhost/orders
// JSON Body:
// {
// 		"customer_name": "customer",
// 		"ordered_at": "2022-07-07 15:50:13.793654",
// 		"items":
// 		[
// 			{
// 				"item_code": "itemcode1",
// 				"description": "itemdescription1",
// 				"quantity": 1
// 			},
// 			{
// 				"item_code": "itemcode2",
// 				"description": "itemdescription2",
// 				"quantity": 1
// 			}
// 		]
// }
func createOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	decoder := json.NewDecoder(r.Body)
	var order entity.OrderWithItems

	if err := decoder.Decode(&order); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}

	orders, err := Helper.Tesdb.CreateOrder(ctx, order)
	if err != nil {
		WriteJsonResp(w, StatusError, err.Error())
		return
	}
	WriteJsonResp(w, StatusSuccess, orders)
}

// updateOrderHandler
// Method: POST / PUT
// Example: localhost/orders/1
// JSON Body:
// {
//		"customer_name": "customer",
//		"ordered_at": "2022-07-27T12:42:31Z",
//		"items":
//		[
//			{
//				"item_code": "itemcode1",
//				"description": "itemdescription1",
//				"quantity": 1
//			},
//			{
//				"item_code": "itemcode2",
//				"description": "itemdescription2",
//				"quantity": 1
//			}
//		]
// }
func updateOrderHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()

	if id != "" { // get by id
		if idInt, err := strconv.Atoi(id); err == nil {
			if orders, err := Helper.Tesdb.GetOrderByID(ctx, idInt); err != nil {
				WriteJsonResp(w, StatusError, err.Error())
				return
			} else if idInt != orders.OrderId {
				WriteJsonResp(w, StatusError, "Data not exists")
				return
			} else {
				decoder := json.NewDecoder(r.Body)
				var order entity.OrderWithItems
				if err := decoder.Decode(&order); err != nil {
					w.Write([]byte("error decoding json body"))
					return
				}

				orders, err := Helper.Tesdb.UpdateOrder(ctx, idInt, order)
				if err != nil {
					WriteJsonResp(w, StatusError, err.Error())
					return
				}
				WriteJsonResp(w, StatusSuccess, orders)
			}
		}
	}
}

// deleteOrderHandler
// Method: DELETE
// Example: localhost/orders/1
func deleteOrderHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if id != "" { // get by id
		if idInt, err := strconv.Atoi(id); err == nil {
			if orders, err := Helper.Tesdb.GetOrderByID(ctx, idInt); err != nil {
				WriteJsonResp(w, StatusError, err.Error())
				return
			} else if idInt != orders.OrderId {
				WriteJsonResp(w, StatusError, "Data not exists")
				return
			} else {
				order, err := Helper.Tesdb.DeleteOrder(ctx, idInt)
				if err != nil {
					WriteJsonResp(w, StatusError, err.Error())
					return
				}
				WriteJsonResp(w, StatusSuccess, order)
			}
		}
	}
}
