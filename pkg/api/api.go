package api

import (
	"BookM/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API приложения.
type API struct {
	r  *mux.Router       // маршрутизатор запросов
	db storage.Interface // база данных
}

// маршрутизатор запросов создается var router = mux.NewRouter()

// Конструктор API.
func New(db storage.Interface) *API {
	api := API{
		r:  mux.NewRouter(),
		db: db,
	}
	api.endpoints() //func main(){x := db.New() v:= New(x) v.endpoints() }
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
// api.r - mux.NewRouter()

// HandleFunc - зарезервируемая функция
func (api *API) endpoints() {
	api.r.HandleFunc("/orders", api.ordersHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/orders/{id}", api.orderOneHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/orders", api.addOrderHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/orders/{id}", api.updateOrderHandler).Methods(http.MethodPatch)
	api.r.HandleFunc("/orders/{id}", api.deleteOrderHandler).Methods(http.MethodDelete)
	api.r.HandleFunc("/books", api.booksHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/books", api.addBooksHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/books/{id}", api.updateBooksHandler).Methods(http.MethodPatch)
	api.r.HandleFunc("/books/{id}", api.deleteBooksHandler).Methods(http.MethodDelete)

}

// ordersHandler возвращает все заказы.
func (api *API) ordersHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из БД.
	orders, err := api.db.Orders()
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// ordersHandler возвращает заказ по ID
func (api *API) orderOneHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из БД.
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var o storage.Order
	o.ID = id
	orders, err := api.db.OrdersOne(o)
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// addOrderHandler создает новый заказ.
func (api *API) addOrderHandler(w http.ResponseWriter, r *http.Request) {
	var o storage.Order
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := api.db.AddOrders(o)
	w.Write([]byte(strconv.Itoa(id)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// updateOrderHandler обновляет данные заказа по ID
func (api *API) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра {id} из пути запроса.
	// Например, /orders/45.
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Декодирование в переменную тела запроса,
	// которое должно содержать JSON-представление
	// обновляемого объекта.
	var o storage.Order
	err = json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	o.ID = id
	// Обновление данных в БД.
	api.db.UpdateOrder(o)
	// Отправка клиенту статуса успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

// deleteOrderHandler удаляет заказ по ID
func (api *API) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var o storage.Order
	o.ID = id
	api.db.DeleteOrder(o)
	w.WriteHeader(http.StatusOK)
}

// ordersHandler возвращает все книги
func (api *API) booksHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из БД.
	book, err := api.db.Books()
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
	var b storage.Book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.AddBooks(b)
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
	// Например, /orders/45.
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Декодирование в переменную тела запроса,
	// которое должно содержать JSON-представление
	// обновляемого объекта.
	var b storage.Book
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b.Id = id
	// Обновление данных в БД.
	api.db.UpdateBook(b)
	// Отправка клиенту статуса успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

// deleteOrderHandler удаляет книгу по ID
func (api *API) deleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var b storage.Book
	b.Id = id
	api.db.DeleteBook(b)
	w.WriteHeader(http.StatusOK)
}
