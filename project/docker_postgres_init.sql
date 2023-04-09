CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL UNIQUE,
  fullname VARCHAR(255) NOT NULL,
  first_order_id INT,
  password VARCHAR(100) NOT NULL,
  role VARCHAR(10) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP,
  FOREIGN KEY (first_order_id) REFERENCES orders_items(id)
);

CREATE TABLE orders_items (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  price INT NOT NULL,
  expired_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE TABLE orders_histories (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  order_item_id INT NOT NULL,
  descriptions TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (order_item_id) REFERENCES orders_items(id) ON DELETE CASCADE
);

INSERT INTO users (username,fullname, first_order_id, password, role, created_at, updated_at, deleted_at)
VALUES ('admin','admin', null, '$2a$12$ZR3sqMWXNcCEiTy.sJ1jkOC0DN75Pp2UN6oBH2ZdWHxskJcObfECi', 'admin', NOW(), null, null);
