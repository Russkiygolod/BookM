package mongo

import (
	"context"
	"log"
	"sync"
	"time"

	"BookM/pkg/model/Book"
	"BookM/pkg/model/Order"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const (
	databaseName   = "bookM"
	collectionName = "orders"
	collectionBook = "books"
)

type Store struct {
	cl *mongo.Client
	m  sync.Mutex
}

func New(P string) (*Store, error) {
	var Mongo Store
	//подключение к БД
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(P))
	if err != nil {
		log.Fatal(err) // todo  Fatalf используй чтобы более информативно описать на каком этапе проищошла ошибка
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	Mongo.cl = client
	return &Mongo, err
}

// закрывает подкючение
func (s *Store) Close() {
	log.Fatal(s.cl.Disconnect(context.Background()))
}

// выдает все книги
func (s *Store) GetBooks(ctx context.Context) ([]Book.Book, error) {
	collection := s.cl.Database(databaseName).Collection(collectionBook)
	filter := bson.M{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var books []Book.Book
	for cur.Next(ctx) {
		var book Book.Book
		err := cur.Decode(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, cur.Err()
}

// добавить книги
func (s *Store) AddBooks(ctx context.Context, B Book.Book) error {
	s.m.Lock()
	defer s.m.Unlock()
	collection := s.cl.Database(databaseName).Collection(collectionBook)
	_, err := collection.InsertOne(ctx, B)
	if err != nil {
		return err
	}
	return nil
}

// изменить книгу
func (s *Store) UpdateBook(ctx context.Context, B Book.Book) error {
	s.m.Lock()
	defer s.m.Unlock()
	collection := s.cl.Database(databaseName).Collection(collectionBook)
	filter := bson.M{"id": B.Id}
	update := bson.M{"$set": bson.M{"Title": B.Title, "Author": B.Author, "Price": B.Price}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// удалить книгу
func (s *Store) DeleteBook(ctx context.Context, B Book.Book) error {
	collection := s.cl.Database(databaseName).Collection(collectionBook)
	_, err := collection.DeleteMany(ctx, bson.M{"id": B.Id})
	if err != nil {
		return err
	}
	return nil
}

// выдает все заказы
func (s *Store) GetOrders(ctx context.Context) ([]Order.Order, error) {
	collection := s.cl.Database(databaseName).Collection(collectionName)
	filter := bson.M{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var orders []Order.Order
	for cur.Next(ctx) {
		var order Order.Order
		err := cur.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, cur.Err()
}

// создает новый заказ
func (s *Store) AddOrders(ctx context.Context, O Order.Order) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var books []Book.Book
	var Price int
	collectionB := s.cl.Database(databaseName).Collection(collectionBook)
	for i := 0; i < len(O.Books); i++ {
		book := O.Books[i]
		filter := bson.M{"id": book.Id}
		cur, err := collectionB.Find(ctx, filter)
		if err != nil {
			return 0, err
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			err = cur.Decode(&book)
			if err != nil {
				return 0, err
			}
			Price = Price + book.Price
			books = append(books, book)
		}
	}
	O.Books = books
	O.Date = time.Now().Format("02/01/2006")
	O.Price = Price
	collection := s.cl.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(ctx, O)
	if err != nil {
		return 0, err
	}
	return O.ID, nil
}

// обновляет документы из БД.
func (s *Store) UpdateOrder(ctx context.Context, O Order.Order) error {
	s.m.Lock()
	defer s.m.Unlock()
	var books []Book.Book
	var Price int
	collectionB := s.cl.Database(databaseName).Collection("books")
	for i := 0; i < len(O.Books); i++ {
		b := O.Books[i]
		filter := bson.M{"id": b.Id}
		cur, err := collectionB.Find(ctx, filter)
		if err != nil {
			return err
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			err := cur.Decode(&b)
			if err != nil {
				return err
			}
			Price = Price + b.Price
			books = append(books, b)
		}
	}
	O.Books = books
	O.Price = Price
	collection := s.cl.Database(databaseName).Collection(collectionName)
	filter := bson.M{"id": O.ID}
	var update bson.M
	if O.DeliveryAddress == "" {
		update = bson.M{"$set": bson.M{"Books": O.Books, "Price": O.Price}}
	} else {
		update = bson.M{"$set": bson.M{"DeliveryAddress": O.DeliveryAddress, "Books": O.Books, "Price": O.Price}}
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// удаляет заказ
func (s *Store) DeleteOrder(ctx context.Context, O Order.Order) error {
	collection := s.cl.Database(databaseName).Collection(collectionName)
	_, err := collection.DeleteMany(ctx, bson.M{"id": O.ID})
	if err != nil {
		return err
	}
	return nil
}

// выдает заказ по ID
func (s *Store) GetOrderByID(ctx context.Context, O Order.Order) ([]Order.Order, error) {
	collection := s.cl.Database(databaseName).Collection(collectionName)
	filter := bson.M{"id": O.ID}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var orders []Order.Order
	for cur.Next(ctx) {
		var order Order.Order
		err := cur.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, cur.Err()
}
