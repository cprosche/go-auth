-- @block
DROP TABLE users;
-- @block
USE auth;
-- @block
CREATE TABLE users(
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    pw VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- @block
INSERT INTO users (username, email, pw)
VALUES(
        'caderosche',
        'caderosche@gmail.com',
        'password'
    ),
    (
        'bri',
        'brihadley@gmail.com',
        'password2'
    ),
    (
        'dave',
        'dave@gmail.com',
        'password3'
    );
-- @block
SELECT *
FROM users;
-- @block
CREATE DATABASE auth;
-- @block
DELETE FROM users;