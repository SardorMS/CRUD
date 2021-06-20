--таблица зарегестрированных покупателей
CREATE TABLE IF NOT EXIST customers
(
    id       BIGSERIAL PRIMARY KEY,
    name     TEXT      NOT NULL,
    phone    TEXT      NOT NULL UNIQUE,
    password TEXT      NOT NULL,
    active   BOOLEAN   NOT NULL DEFAULT TRUE, 
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--таблица сотрудников
CREATE TABLE IF NOT EXIST managers 
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT      NOT NULL,
    phone      TEXT      NOT NULL UNIQUE,
    password   TEXT      NOT NULL,
    salary     INTEGER   NOT NULL CHECK (salary > 0),
    plan       INTEGER   NOT NULL CHECK (plan > 0),
    boss_id    BIGINT    REFERENCES managers,
    department TEXT,
    is_admin   BOOLEAN   NOT NULL DEFAULT TRUE,
    active     BOOLEAN   NOT NULL DEFAULT TRUE, 
    created    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--таблица токенов покупателей
CREATE TABLE IF NOT EXIST customers_tokens 
(
    token        TEXT      NOT NULL,
    customer_id  BIGINT    NOT NULL REFERENCES customers,
    expire       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

--таблица токенов сотрудников
CREATE TABLE IF NOT EXIST managers_tokens 
(
    token        TEXT      NOT NULL,
    manager_id   BIGINT    NOT NULL REFERENCES managers,
    expire       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

-- таблица товаров
CREATE TABLE IF NOT EXIST products 
(
    id      BIGSERIAL PRIMARY KEY,
    name    TEXT      NOT NULL,
    price   INTEGER   NOT NULL CHECK (price > 0),
    qty     INTEGER   NOT NULL DEFAULT 0 CHECK (qty >= 0),
    active  BOOLEAN   NOT NULL DEFAULT TRUE, 
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


--таблица продаж
CREATE TABLE IF NOT EXIST sales
(
    id          BIGSERIAL PRIMARY KEY,
    manager_id  BIGINT    NOT NULL REFERENCES managers,
    customer_id BIGINT    REFERENCES customers,
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--конкретные позиции в продаже (чек)
CREATE TABLE IF NOT EXIST sale_positions
(
    id          BIGSERIAL PRIMARY KEY,
    sale_id     BIGINT    NOT NULL REFERENCES sales,
    product_id  BIGINT    NOT NULL REFERENCES products,
    name        TEXT      NOT NULL,
    price       INTEGER   NOT NULL CHECK (price >= 0),
    qty         INTEGER   NOT NULL DEFAULT 0 CHECK (qty >= 0),
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- таблица юзеров (когда для хранения используется одна таблица)
CREATE TABLE IF NOT EXIST users
(
    id       BIGSERIAL PRIMARY KEY,
    name     TEXT      NOT NULL,
    phone    TEXT      NOT NULL,
    password TEXT      NOT NULL,
    roles    TEXT[]    NOT NULL DEFAULT '{}',
    active   BOOLEAN   NOT NULL DEFAULT TRUE, 
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



--DROP TABLE <name_of_table> CASCADE

--DROP TABLE products;
--DROP TABLE managers;
--DROP TABLE customers;
--DROP TABLE sales;
--DROP TABLE sale_positions;
