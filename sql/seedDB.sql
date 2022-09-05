-- @block
DROP TABLE users;
-- @block
USE auth;
-- @block
CREATE TABLE users(
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    pw VARCHAR(255)
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
    );
-- @block
SELECT *
FROM users;
-- @block
CREATE DATABASE auth;