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
    PRIMARY KEY (id)
);

CREATE TABLE orders (
    id int NOT NULL AUTO_INCREMENT,
    product_id int,
    customer_id int,
    quantity int,
    price decimal(6,2),
    tax decimal(6,2),
    donation decimal(4,2),  -- Renamed from 'decimal' to 'donation'
    timestamp bigint,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id) REFERENCES product(id),
    FOREIGN KEY (customer_id) REFERENCES customer(id)
);

-- populate customer table
INSERT INTO customer (first_name, last_name, email)
VALUES ('Mickey', 'Mouse', 'mmouse@mines.edu'),
       ('Mayar', 'AlAnsari', 'alansari@mines.edu');

-- insert products
INSERT INTO product (product_name, image_name, price, in_stock)
VALUES ('2024 G80 M3', '2024 G80 M3.png', 78000, 0),
       ('2024 S63 AMG', '2024 S63 AMG.png', 183000, 3),
       ('2024 Audi RS7', '2024 Audi RS7.png', 128000, 10);
