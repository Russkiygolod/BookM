package Order

import "BookM/pkg/model/Book"

type Order struct {
	ID              int         // номер заказа
	DeliveryAddress string      // адрес доставки
	Date            string      // дата
	Books           []Book.Book `json:"books"` // книги в заказе
	Price           int         //общая цена
}
