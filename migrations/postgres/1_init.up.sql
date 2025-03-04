CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	telegram_id BIGINT UNIQUE NOT NULL CHECK (telegram_id > 0)
);

CREATE TABLE IF NOT EXISTS favorites (
	product_id SERIAL PRIMARY KEY,
	product_name VARCHAR(300) NOT NULL,
	product_link VARCHAR(600) UNIQUE NOT NULL,
	base_price INT NOT NULL,
	product_brand VARCHAR(100) NOT NULL,
	supplier VARCHAR(100) NOT NULL,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tracked_products (
	tracked_id SERIAL PRIMARY KEY,
	product_name VARCHAR(300) NOT NULL,
	price_down_bound INT NOT NULL,
	price_up_bound INT NOT NULL,
	market_filter TEXT[] NOT NULL,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
