package storage

import (
	"BookM/pkg/model/Book"
	"BookM/pkg/model/Order"
	"context"
)

type Interface interface {
	GetOrders(context.Context) ([]Order.Order, error)                 // все заказы
	GetOrderByID(context.Context, Order.Order) ([]Order.Order, error) // заказ по id
	AddOrders(context.Context, Order.Order) (int, error)              // добавить заказ
	UpdateOrder(context.Context, Order.Order) error                   // изменить заказ
	DeleteOrder(context.Context, Order.Order) error                   // удалить заказ

	GetBooks(context.Context) ([]Book.Book, error) // все книги
	AddBooks(context.Context, Book.Book) error     // добавить книги
	UpdateBook(context.Context, Book.Book) error   // изменить книгу
	DeleteBook(context.Context, Book.Book) error   // удалить книгу
}
