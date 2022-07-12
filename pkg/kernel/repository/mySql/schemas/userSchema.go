package schemas

const UserSchema = `
CREATE TABLE IF NOT EXISTS users (
	uuid VARCHAR(36) PRIMARY KEY,
	alias VARCHAR(30) NOT NULL UNIQUE,
	name VARCHAR(30) NOT NULL,
	second_name VARCHAR(30) NOT NULL,
	email VARCHAR(30) NOT NULL UNIQUE,
	password BLOB
);
`
