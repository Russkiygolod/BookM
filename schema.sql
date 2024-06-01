-- этот файл думаю вообще не нужно в гитхаб сувать
DROP TABLE IF EXISTS books, orders;

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    price INTEGER
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    deliveryAddress TEXT NOT NULL,
    date TEXT NOT NULL,
    books_id INTEGER[],
    price INTEGER
);


INSERT INTO books (title, author, price) VALUES ('Старик и Море', 'Хемингуэй', 300);
INSERT INTO books (title, author, price) VALUES ('Герой нашего времени', 'Лермонтов', 400);
INSERT INTO books (title, author, price) VALUES ('Теория всего', 'Стивен Хоккинг', 200);

INSERT INTO orders (deliveryAddress, date, books_id, price) VALUES ('Касимов', '28/04/2024', '{1,2}', 700); 

SELECT orders.id, orders.deliveryAddress, orders.date, orders.price, (
SELECT ARRAY_TO_JSON(ARRAY_AGG(ROW_TO_JSON(b.*))) AS array_to_json FROM (
    SELECT
    title,
    author,
    price
    FROM books where id = any (orders.books_id)
    ) AS b
) AS books from orders
ORDER BY orders.id ASC
