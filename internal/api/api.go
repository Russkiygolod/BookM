package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"BookM/internal"
	"BookM/pkg/model/Book"
	"BookM/pkg/model/Order"

	"github.com/gorilla/mux"
)

type API struct {
	r *mux.Router // маршрутизатор запросов
	I internal.Inter
}

func New(I internal.Inter) *API {
	api := API{
		r: mux.NewRouter(),
		I: I,
	}
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.HandleFunc("/orders", api.GetOrdersHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/orders/{id}", api.GetOrderOneHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/orders", api.addOrderHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/orders/{id}", api.updateOrderHandler).Methods(http.MethodPatch)
	api.r.HandleFunc("/orders/{id}", api.deleteOrderHandler).Methods(http.MethodDelete)
	api.r.HandleFunc("/books", api.GetBooksHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/books", api.addBooksHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/books/{id}", api.updateBooksHandler).Methods(http.MethodPatch)
	api.r.HandleFunc("/books/{id}", api.deleteBooksHandler).Methods(http.MethodDelete)
}

// ordersHandler возвращает все заказы.
func (api *API) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из БД
	orders, err := api.I.GetOrders(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	// Отправка данных клиенту в формате JSON.
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// ordersHandler возвращает заказ по ID
func (api *API) GetOrderOneHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из БД.
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var O Order.Order
	O.ID = id
	orders, err := api.I.GetOrderByID(context.Background(), O)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	// Отправка данных клиенту в формате JSON.
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// addOrderHandler создает новый заказ.
func (api *API) addOrderHandler(w http.ResponseWriter, r *http.Request) {
	var O Order.Order
	err := json.NewDecoder(r.Body).Decode(&O)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := api.I.AddOrders(context.Background(), O)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.Write([]byte(strconv.Itoa(id)))
		w.WriteHeader(http.StatusOK)
	}
}

// updateOrderHandler обновляет данные заказа по ID
func (api *API) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра {id} из пути запроса.
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var O Order.Order
	err = json.NewDecoder(r.Body).Decode(&O)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	O.ID = id
	// Обновление данных в БД.
	err = api.I.UpdateOrder(context.Background(), O)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// deleteOrderHandler удаляет заказ по ID
func (api *API) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var O Order.Order
	O.ID = id
	err = api.I.DeleteOrder(context.Background(), O)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// ordersHandler возвращает все книги
func (api *API) GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из БД.
	book, err := api.I.GetBooks(context.Background())
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// addOrderHandler добавляет новую книгу
func (api *API) addBooksHandler(w http.ResponseWriter, r *http.Request) {
	var b Book.Book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.I.AddBooks(context.Background(), b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// updateOrderHandler обновляет данные книги по ID
func (api *API) updateBooksHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра {id} из пути запроса.
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var b Book.Book
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b.Id = id
	// Обновление данных в БД.
	err = api.I.UpdateBook(context.Background(), b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// deleteOrderHandler удаляет книгу по ID
func (api *API) deleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var b Book.Book
	b.Id = id
	err = api.I.DeleteBook(context.Background(), b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
