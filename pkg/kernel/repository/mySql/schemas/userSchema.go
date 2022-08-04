package schemas

const UserSchema = `
CREATE TABLE IF NOT EXISTS users (
	id INT AUTO_INCREMENT,	
	uuid VARCHAR(36),
	alias VARCHAR(30) NOT NULL UNIQUE,
	name VARCHAR(30) NOT NULL,
	second_name VARCHAR(30) NOT NULL,
	email VARCHAR(30) NOT NULL UNIQUE,
	password BLOB,
	PRIMARY KEY(id, uuid),
	INDEX (id)
);
`
