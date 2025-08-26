package handler

import (
	"encoding/json"
	"kafkapractisel0/apiError"
	"kafkapractisel0/repo/cache"
	"kafkapractisel0/services"
	"net/http"
	"strconv"
)

type OrderHandler interface {
	GetOrderById(w http.ResponseWriter, req *http.Request)
}

type orderHandler struct {
	orderService services.OrderService
	cache        cache.CacheRepo
}

func NewOrderHandler(orderService services.OrderService, cacheRepo cache.CacheRepo) OrderHandler {
	return &orderHandler{orderService, cacheRepo}
}

func (h *orderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case http.MethodGet:
		uid := r.PathValue("order_uid")
		uidParsed, err := strconv.Atoi(uid)
		if err != nil {
			apiError.BackendErrorWrite(w, apiError.BadRequest)
			return
		}

		res, cacheExist := h.cache.Get(uidParsed, true)
		if cacheExist {
			_, _ = w.Write(res)
			return
		}

		order, err := h.orderService.GetOrderById(int64(uidParsed))
		if err != nil {
			apiError.BackendErrorWrite(w, apiError.BadRequest)
			return
		}
		jData, errM := json.Marshal(order)
		if errM != nil {
			apiError.BackendErrorWrite(w, apiError.MarshalError)
			return
		}
		_, _ = w.Write(jData)
		h.cache.Set(uidParsed, jData)
	default:
		apiError.BackendErrorWrite(w, apiError.MethodNotAllowed)
	}
}
