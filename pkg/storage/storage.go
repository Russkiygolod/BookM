package storage

type Order struct {
	ID              int    // номер заказа
	DeliveryAddress string // адрес доставки
	Date            string // дата
	Books           []Book `json:"books"` // книги в заказе
	Price           int    //общая цена
}

type Book struct {
	Id     int    //артикул
	Title  string //название
	Author string //автор
	Price  int    //цена
}

type Interface interface {
	Orders() ([]Order, error)         // все заказы
	OrdersOne(Order) ([]Order, error) // заказ по id
	AddOrders(Order) (int, error)     // добавить заказ
	UpdateOrder(Order) error          // изменить заказ
	DeleteOrder(Order) error          // удалить заказ
	Books() ([]Book, error)           // все книги
	AddBooks(Book) error              // добавить книги
	UpdateBook(Book) error            // изменить книгу
	DeleteBook(Book) error            // удалить книгу
}
