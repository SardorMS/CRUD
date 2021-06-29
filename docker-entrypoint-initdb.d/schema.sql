-- Table of registred customers.
CREATE TABLE IF NOT EXISTS customers
(
    id       BIGSERIAL PRIMARY KEY,
    name     TEXT      NOT NULL,
    phone    TEXT      NOT NULL UNIQUE,
    password TEXT      NOT NULL,
    active   BOOLEAN   NOT NULL DEFAULT TRUE, 
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table of managers.
CREATE TABLE IF NOT EXISTS managers 
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT      NOT NULL,
    phone      TEXT      NOT NULL UNIQUE,
    password   TEXT,
    salary     INTEGER   NOT NULL DEFAULT 0,
    plan       INTEGER   NOT NULL DEFAULT 0 ,
    boss_id    BIGINT    REFERENCES managers,
    department TEXT,
    is_admin   BOOLEAN   NOT NULL DEFAULT TRUE,
    active     BOOLEAN   NOT NULL DEFAULT TRUE, 
    created    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table of customers tokens.
CREATE TABLE IF NOT EXISTS customers_tokens 
(
    token        TEXT      NOT NULL UNIQUE,
    customer_id  BIGINT    NOT NULL REFERENCES customers,
    expire       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

-- Table of managers tokens.
CREATE TABLE IF NOT EXISTS managers_tokens 
(
    token        TEXT      NOT NULL UNIQUE,
    manager_id   BIGINT    NOT NULL REFERENCES managers,
    expire       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

-- Table of products.
CREATE TABLE IF NOT EXISTS products 
(
    id      BIGSERIAL PRIMARY KEY,
    name    TEXT      NOT NULL,
    price   INTEGER   NOT NULL CHECK (price > 0),
    qty     INTEGER   NOT NULL DEFAULT 0 CHECK (qty >= 0),
    active  BOOLEAN   NOT NULL DEFAULT TRUE, 
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Table of sales.
CREATE TABLE IF NOT EXISTS sales
(
    id          BIGSERIAL PRIMARY KEY,
    manager_id  BIGINT    NOT NULL REFERENCES managers,
    customer_id BIGINT    NOT NULL,
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Table about sales positions (cheque).
CREATE TABLE IF NOT EXISTS sale_positions
(
    id          BIGSERIAL PRIMARY KEY,
    sale_id     BIGINT    NOT NULL REFERENCES sales,
    product_id  BIGINT    NOT NULL REFERENCES products,
    price       INTEGER   NOT NULL CHECK (price >= 0),
    qty         INTEGER   NOT NULL DEFAULT 0 CHECK (qty >= 0),
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Table of users (when a single table is used for storage).
CREATE TABLE IF NOT EXISTS users
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
--DROP TABLE managers_tokens;
--DROP TABLE customers;
--DROP TABLE customers_tokens;
--DROP TABLE sales;
--DROP TABLE sale_positions;
