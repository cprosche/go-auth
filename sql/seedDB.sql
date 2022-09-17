-- @block
DROP TABLE users;
-- @block
USE auth;
-- @block
CREATE TABLE users(
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) UNIQUE,
    pw VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- @block
SELECT *
FROM users;
-- @block
CREATE DATABASE auth;
-- @block
DELETE FROM users;