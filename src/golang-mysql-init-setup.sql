CREATE DATABASE IF NOT EXISTS golang_db;

SELECT 
    *
FROM
    mysql.user
WHERE
    User = 'golang';

DROP USER IF EXISTS 'golang'@'%';
DROP USER IF EXISTS 'golang'@'localhost';
DROP USER IF EXISTS 'golang'@'127.0.0.1';

CREATE USER 'golang'@'%' IDENTIFIED BY 'golang';
CREATE USER 'golang'@'localhost' IDENTIFIED BY 'golang';
CREATE USER 'golang'@'127.0.0.1' IDENTIFIED BY 'golang';

GRANT ALL ON golang_db.* to 'golang'@'%';
GRANT ALL ON golang_db.* to 'golang'@'localhost';
GRANT ALL ON golang_db.* to 'golang'@'127.0.0.1';

CREATE TABLE golang_db.dnslookup (dn VARCHAR(100) PRIMARY KEY, ip VARCHAR(15) NOT NULL);

SELECT * FROM golang_db.dnslookup;