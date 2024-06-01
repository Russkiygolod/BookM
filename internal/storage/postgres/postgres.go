package postgres

import (
	"context"
	"log"
	"sync"
	"time"

	"BookM/pkg/model/Book"
	"BookM/pkg/model/Order"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	cl *pgxpool.Pool
	m  sync.Mutex
}

// Конструктор объекта хранилища.
func New(conf *pgxpool.Config) *Store {
	var postgres Store
	var err error

	postgres.cl, err = pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		log.Fatalf("Failed to init DB conf - %v", err)
	}

	//postgres.orders = &Order.Order{}
	return &postgres
}

// выдает все заказы
func (s *Store) GetOrders(ctx context.Context) ([]Order.Order, error) {
	s.m.Lock()
	defer s.m.Unlock()
	rows, err := s.cl.Query(
		ctx, `
			SELECT orders.id, orders.deliveryAddress, orders.date, orders.price, (
				SELECT ARRAY_TO_JSON(ARRAY_AGG(ROW_TO_JSON(b.*))) AS array_to_json FROM (
					SELECT
					id,
					title,
					author,
					price
					FROM books where id = any (orders.books_id)
					) AS b
			) AS books from orders
			ORDER BY orders.id ASC
`)
	if err != nil {
		return nil, err
	}
	var posts []Order.Order
	for rows.Next() {
		var t Order.Order
		err = rows.Scan(
			&t.ID, &t.DeliveryAddress, &t.Date, &t.Price, &t.Books,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, t)
	}
	return posts, rows.Err()
}

// выдает заказ по ID
func (s *Store) GetOrderByID(ctx context.Context, O Order.Order) ([]Order.Order, error) {
	s.m.Lock()
	defer s.m.Unlock()
	rows, err := s.cl.Query(
		ctx,
		`	SELECT orders.id, orders.deliveryAddress, orders.date, orders.price, (
			SELECT ARRAY_TO_JSON(ARRAY_AGG(ROW_TO_JSON(b.*))) AS array_to_json FROM (
				SELECT
				id,
				title,
				author,
				price
				FROM books where id = any (orders.books_id)
				) AS b
			) AS books from orders
			WHERE orders.id = $1
`, O.ID,
	)
	if err != nil {
		return nil, err
	}
	var posts []Order.Order
	for rows.Next() {
		var t Order.Order
		err = rows.Scan(
			&t.ID, &t.DeliveryAddress, &t.Date, &t.Price, &t.Books,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, t)
	}
	return posts, rows.Err()
}

// создает новый заказ
func (s *Store) AddOrders(ctx context.Context, O Order.Order) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var date string
	Books := O.Books
	var bookID []int
	var bookPrice int
	for i := 0; i < len(Books); i++ {
		b := Books[i]
		rows, err := s.cl.Query(
			ctx,
			`
			SELECT price from books 
			WHERE id = $1; 
		`,
			b.Id,
		)
		if err != nil {
			return 0, err
		}
		var price int

		for rows.Next() {
			err = rows.Scan(
				&b.Price,
			)
			if err != nil {
				return 0, err
			}
			price = b.Price
			bookID = append(bookID, b.Id)
			bookPrice = bookPrice + price
		}
	}

	date = time.Now().Format("02/01/2006")
	rows, err := s.cl.Query(
		ctx,
		`
		INSERT INTO orders (deliveryAddress, date, books_id, price) VALUES ($1, $2, $3, $4) RETURNING id; 
	`,
		O.DeliveryAddress,
		date,
		bookID,
		bookPrice,
	)
	if err != nil {
		return 0, err
	}
	var id int
	for rows.Next() {
		var t Order.Order
		err = rows.Scan(
			&t.ID,
		)
		if err != nil {
			return 0, err
		}
		id = t.ID
	}
	return id, nil
}

// обновляет заказ
func (s *Store) UpdateOrder(ctx context.Context, O Order.Order) error {
	s.m.Lock()
	defer s.m.Unlock()

	if O.DeliveryAddress != "" {
		_, err := s.cl.Exec(
			ctx, `
			UPDATE orders
			SET deliveryAddress = $2
			WHERE id = $1;
		`, O.ID, O.DeliveryAddress,
		)
		if err != nil {
			return err
		}
	}
	if O.Books != nil {
		var Books = O.Books
		var bookPrice int
		var bookID []int
		for i := 0; i < len(Books); i++ {
			b := Books[i]
			rows, err := s.cl.Query(
				ctx,
				`
				SELECT price from books 
                WHERE id = $1; 
			`,
				b.Id,
			)
			if err != nil {
				return err
			}
			var price int
			for rows.Next() {
				var t Book.Book
				err = rows.Scan(
					&t.Price,
				)
				if err != nil {
					return err
				}
				price = t.Price
				bookID = append(bookID, b.Id)
				bookPrice = bookPrice + price
			}
		}
		_, err := s.cl.Exec(
			ctx, `
			UPDATE orders
			SET books_id = $2,
			price = $3
			WHERE id = $1;
		`, O.ID, bookID, bookPrice,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// удаляет заказ
func (s *Store) DeleteOrder(ctx context.Context, O Order.Order) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.cl.Exec(
		ctx, `DELETE FROM orders WHERE id = $1 ;`, O.ID)
	if err != nil {
		return err
	}
	return nil
}

// выдает все книги
func (s *Store) GetBooks(ctx context.Context) ([]Book.Book, error) {
	s.m.Lock()
	defer s.m.Unlock()
	rows, err := s.cl.Query(
		ctx, `
		SELECT * from books
		ORDER BY id ASC;
`)
	if err != nil {
		return nil, err
	}
	var book []Book.Book
	for rows.Next() {
		var b Book.Book
		err = rows.Scan(
			&b.Id,
			&b.Title,
			&b.Author,
			&b.Price,
		)
		if err != nil {
			return nil, err
		}
		book = append(book, b)
	}
	return book, rows.Err()
}

// добавляет книгу
func (s *Store) AddBooks(ctx context.Context, B Book.Book) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.cl.Exec(
		ctx,
		`
		INSERT INTO books (title, author, price) VALUES ($1, $2, $3);
	`,
		B.Title,
		B.Author,
		B.Price,
	)
	if err != nil {
		return err
	}
	return nil
}

// обновляет данные книги
func (s *Store) UpdateBook(ctx context.Context, B Book.Book) error {
	s.m.Lock()
	defer s.m.Unlock()
	if B.Author != "" {
		_, err := s.cl.Exec(
			ctx,
			`
			UPDATE books
			SET author = $2
			WHERE id = $1;
		`, B.Id, B.Author,
		)
		if err != nil {
			return err
		}
	}
	if B.Price != 0 {
		_, err := s.cl.Exec(
			ctx,
			`
			UPDATE books
			SET price = $2
			WHERE id = $1;
		`, B.Id, B.Price,
		)
		if err != nil {
			return err
		}
	}
	if B.Title != "" {
		_, err := s.cl.Exec(
			ctx,
			`
			UPDATE books
			SET title = $2
			WHERE id = $1;
		`, B.Id, B.Title,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// удаляет книгу
func (s *Store) DeleteBook(ctx context.Context, B Book.Book) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.cl.Exec(
		ctx, `DELETE FROM books WHERE id = $1 ;`, B.Id)
	if err != nil {
		return err
	}
	return nil
}
