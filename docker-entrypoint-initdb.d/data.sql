
-- Insert Examples.

INSERT INTO managers (name, phone, password, is_admin)
VALUES ('vasya', '+992000000001', '$2a$10$OaUtjCNv2DT5x/dXcV.P3eYkIPIRtBr/v8Nluwifz6brSkfyXOh6m', true);

INSERT INTO products (name, price, qty)
VALUES ('Pizza', 200, 10);

INSERT INTO sales (manager_id, customer_id)
VALUES (1, 1);

INSERT INTO sale_positions (sale_id, product_id, name, qty, price)
VALUES (1, 1, 'Pizza', 5, 200);


ALTER SEQUENCE customers_id_seq RESTART WITH 1;
ALTER SEQUENCE products_id_seq RESTART WITH 1;
ALTER SEQUENCE sales_id_seq RESTART WITH 1;
ALTER SEQUENCE sale_positions_id_seq RESTART WITH 1;
ALTER SEQUENCE managers_id_seq RESTART WITH 1;
DELETE FROM managers;
DELETE FROM customers;
DELETE FROM customers_tokens;
DELETE FROM products;
DELETE FROM sales;
DELETE FROM sale_positions;
