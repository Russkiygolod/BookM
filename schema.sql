
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
