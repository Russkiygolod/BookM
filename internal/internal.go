package internal

import (
	"BookM/internal/storage"
	"BookM/pkg/model/Book"
	"BookM/pkg/model/Order"
	"context"
)

type Inter struct {
	s storage.Interface // база данных
}

func New(s storage.Interface) Inter {
	var Inters Inter
	Inters.s = s
	return Inters
}

func (I Inter) GetOrders(ctx context.Context) ([]Order.Order, error) {
	return I.s.GetOrders(ctx)
}
func (I Inter) GetOrderByID(ctx context.Context, O Order.Order) ([]Order.Order, error) {
	return I.s.GetOrderByID(ctx, O)
}
func (I Inter) AddOrders(ctx context.Context, O Order.Order) (int, error) {
	return I.s.AddOrders(ctx, O)
}
func (I Inter) UpdateOrder(ctx context.Context, O Order.Order) error {
	return I.s.UpdateOrder(ctx, O)
}
func (I Inter) DeleteOrder(ctx context.Context, O Order.Order) error {
	return I.s.DeleteOrder(ctx, O)
}
func (I Inter) GetBooks(ctx context.Context) ([]Book.Book, error) {
	return I.s.GetBooks(ctx)
}
func (I Inter) AddBooks(ctx context.Context, B Book.Book) error {
	return I.s.AddBooks(ctx, B)
}
func (I Inter) UpdateBook(ctx context.Context, B Book.Book) error {
	return I.s.UpdateBook(ctx, B)
}
func (I Inter) DeleteBook(ctx context.Context, B Book.Book) error {
	return I.s.DeleteBook(ctx, B)
}
