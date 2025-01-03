USE ctorosuarez;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS customer;
DROP TABLE IF EXISTS product;

CREATE TABLE customer (
    id int NOT NULL AUTO_INCREMENT,
    first_name varchar(255),
    last_name varchar(255),
    email varchar(255),
    PRIMARY KEY (id)
);

CREATE TABLE product (
    id int NOT NULL AUTO_INCREMENT,
    product_name varchar(255),
    image_name varchar(255),
    price decimal(10,2),
    in_stock int,
    inactive TINYINT DEFAULT 0,
    PRIMARY KEY (id)
);

CREATE TABLE orders (
    id int NOT NULL AUTO_INCREMENT,
    customer_first VARCHAR(100),
    customer_last VARCHAR(100),
    product_name VARCHAR(255),
    quantity int,
    price decimal(10,2),
    tax decimal(10,2),
    donation decimal(4,2),
    timestamp bigint,
    PRIMARY KEY (id)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    password VARCHAR(255),
    email VARCHAR(100) UNIQUE NOT NULL,
    role INT NOT NULL
);

-- populate customer table
INSERT INTO customer (first_name, last_name, email)
VALUES ('Mickey', 'Mouse', 'mmouse@mines.edu'),
       ('Mayar', 'AlAnsari', 'alansari@mines.edu');

-- insert products
INSERT INTO product (product_name, image_name, price, in_stock, inactive)
VALUES ('2024 G80 M3', '2024 G80 M3.png', 78000, 10, 0),
       ('2024 S63 AMG', '2024 S63 AMG.png', 183000, 10, 0),
       ('2024 Audi RS7', '2024 Audi RS7.png', 128000, 0, 1);

-- insert users
INSERT INTO users (first_name, last_name, password, email, role)
VALUES 
    ('Frodo', 'Baggins', 'fb', 'fb@mines.edu', 1),
    ('Harry', 'Potter', 'hp', 'hp@mines.edu', 2);
