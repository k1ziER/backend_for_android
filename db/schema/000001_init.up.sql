CREATE TABLE peoples
(
    id SERIAL PRIMARY KEY,
    userName VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    email VARCHAR(255)  UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT false,
    birthday DATE,
    age INTEGER
);

