

INSERT INTO managers (name, phone, password, is_admin)
VALUES ('vasya', '+992000000001', '$2a$10$OaUtjCNv2DT5x/dXcV.P3eYkIPIRtBr/v8Nluwifz6brSkfyXOh6m', true);



ALTER SEQUENCE customers_id_seq RESTART WITH 1;
ALTER SEQUENCE sale_positions_id_seq RESTART WITH 1;
ALTER SEQUENCE managers_id_seq RESTART WITH 1;
DELETE FROM products;
DELETE FROM managers;
DELETE FROM customers;
DELETE FROM sales;
DELETE FROM sale_positions;