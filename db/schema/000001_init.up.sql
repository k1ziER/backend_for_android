CREATE TABLE peoples
(
    id SERIAL PRIMARY KEY,
    userName VARCHAR(255) NOT NULL,
    loginn VARCHAR(255) UNIQUE NOT NULL,
    surname VARCHAR(255),
    email VARCHAR(255)  UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

