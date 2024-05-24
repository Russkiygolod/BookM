package mongo

import (
	"BookM/pkg/storage"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const (
	databaseName   = "bookM"  // имя  БД
	collectionName = "orders" // имя коллекции заказов в БД
	collectionBook = "books"
)

type Store struct {
	cl  *mongo.Client
	db  *storage.Order
	m   sync.Mutex
	err error
}

func New(P string) (*Store, error) {
	var Mongo Store

	//подключение к БД
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(P))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	Mongo.cl = client
	Mongo.db = &storage.Order{}
	Mongo.err = err

	return &Mongo, err
}

// закрывает подкючение
func (s *Store) Close() {
	s.cl.Disconnect(context.Background())
}

// выдает все книги
func (s *Store) Books() ([]storage.Book, error) {
	s.m.Lock()
	defer s.m.Unlock()
	collection := s.cl.Database(databaseName).Collection(collectionBook)
	filter := bson.M{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var books []storage.Book
	for cur.Next(context.Background()) {
		var book storage.Book
		err := cur.Decode(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, cur.Err()
}

// добавить книги
func (s *Store) AddBooks(B storage.Book) error {
	s.m.Lock()
	defer s.m.Unlock()
	collection := s.cl.Database(databaseName).Collection(collectionBook)
	_, err := collection.InsertOne(context.Background(), B)
	if err != nil {
		return err
	}
	return nil
}

// изменить книгу
func (s *Store) UpdateBook(B storage.Book) error {
	collection := s.cl.Database(databaseName).Collection(collectionBook)
	filter := bson.M{"id": B.Id}
	update := bson.M{"$set": bson.M{"Title": B.Title, "Author": B.Author, "Price": B.Price}}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	if err != nil {
		return err
	}
	return nil
}

// удалить книгу
func (s *Store) DeleteBook(B storage.Book) error {
	collection := s.cl.Database(databaseName).Collection(collectionBook)
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.M{"id": B.Id})
	// вывод о успешности операции //
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	if err != nil {
		return err
	}
	return nil
}

// выдает все заказы
func (s *Store) Orders() ([]storage.Order, error) {
	s.m.Lock()
	defer s.m.Unlock()
	collection := s.cl.Database(databaseName).Collection(collectionName)
	filter := bson.M{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var orders []storage.Order
	for cur.Next(context.Background()) {
		var order storage.Order
		err := cur.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, cur.Err()
}

// метод создает новый заказ
func (s *Store) AddOrders(O storage.Order) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var books []storage.Book
	var Price int
	collectionB := s.cl.Database(databaseName).Collection("books")
	for i := 0; i < len(O.Books); i++ {
		book := O.Books[i]
		filter := bson.M{"id": book.Id}
		cur, err := collectionB.Find(context.Background(), filter)
		if err != nil {
			return 0, err
		}
		defer cur.Close(context.Background())
		for cur.Next(context.Background()) {
			err := cur.Decode(&book)
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
	_, err := collection.InsertOne(context.Background(), O)
	if err != nil {
		return 0, err
	}
	return O.ID, nil
}

// Update обновляет документы из БД.
func (s *Store) UpdateOrder(O storage.Order) error {
	s.m.Lock()
	defer s.m.Unlock()
	var books []storage.Book
	var Price int
	collectionB := s.cl.Database(databaseName).Collection("books")
	for i := 0; i < len(O.Books); i++ {
		b := O.Books[i]
		filter := bson.M{"id": b.Id}
		cur, err := collectionB.Find(context.Background(), filter)
		if err != nil {
			return err
		}
		defer cur.Close(context.Background())
		for cur.Next(context.Background()) {
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
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	if err != nil {
		return err
	}
	return nil
}

// удаляет заказ
func (s *Store) DeleteOrder(O storage.Order) error {
	collection := s.cl.Database(databaseName).Collection(collectionName)
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.M{"id": O.ID})
	// вывод о успешности операции //
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	if err != nil {
		return err
	}
	return nil
}

// выдает заказ по ID
func (s *Store) OrdersOne(O storage.Order) ([]storage.Order, error) {
	s.m.Lock()
	defer s.m.Unlock()
	collection := s.cl.Database(databaseName).Collection(collectionName)
	filter := bson.M{"id": O.ID}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var orders []storage.Order
	for cur.Next(context.Background()) {
		var order storage.Order
		err := cur.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, cur.Err()
}
