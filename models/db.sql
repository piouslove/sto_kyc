CREATE TABLE users
(
	id int NOT NULL AUTO_INCREMENT, 
	name varchar(64) NOT NULL, 
	email varchar(64) NOT NULL, 
	address char(42) NOT NULL,
	selector varchar(500) NOT NULL,
	passport varchar(500) NOT NULL,
	status int NOT NULL,
	PRIMARY KEY (id)
) default charset = utf8, ENGINE=InnoDB;