package schemas

const SentEmail = `
CREATE TABLE IF NOT EXISTS sent_email (
	id INT AUTO_INCREMENT,	
	user_uuid VARCHAR(36),
	sent BOOLEAN,
	sent_on DATETIME,
	error VARCHAR(200),
	PRIMARY KEY(id),
	INDEX (id)
);
`
