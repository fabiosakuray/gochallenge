 #  database: testdb
 #  user: testdb
 #  pass: testeDB#12


DROP TABLE user_table;
DROP TABLE reset_table;

# Posgres
    CREATE TABLE IF NOT EXISTS user_table (
	user_id    serial PRIMARY KEY,
	user_login 		VARCHAR(15) UNIQUE  NOT NULL,
	user_pass 		VARCHAR(70) NOT NULL,
	user_name 		VARCHAR(40) NOT NULL,
	user_email 		VARCHAR(40) NOT NULL,
	user_address      VARCHAR(100),
	user_telephone    varchar(20),
    );

    CREATE TABLE IF NOT EXISTS reset_table (
	user_id 	int NOT NULL,
	res_key 	int NOT NULL,
	res_exp_date varchar(100) NOT NULL,
	PRIMARY KEY(res_key )
    );
# ------------------------------------------
# mysql

CREATE TABLE user_table(
  user_id    INT unsigned NOT NULL AUTO_INCREMENT,
  user_login 		VARCHAR(15) NOT NULL,
  user_pass 		VARCHAR(70) NOT NULL,
  user_name 		VARCHAR(40) NOT NULL,
  user_email 		VARCHAR(40) NOT NULL,
  user_address      VARCHAR(100),
  user_telephone    varchar(20),
  PRIMARY KEY(user_id,user_login)
);

    
CREATE TABLE reset_table (
  user_id 	int NOT NULL,
  res_key 	int NOT NULL,
  res_exp_date varchar(100) NOT NULL,
PRIMARY KEY(res_key )
) ;
