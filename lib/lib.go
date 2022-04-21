package lib

import (
       "golang.org/x/crypto/bcrypt"
       "database/sql"
       "math/rand"
       "time"
)
/* Open db connection */
func DBConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "testedb"
    dbPass := "testeDB#12"
    dbName := "testedb"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

// Open Posgres connection
func DBConn()(db *sql.DB)
{
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatalf("Error opening database: %q", err)
    }
    
    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS user_table (
      id       SERIAL,
      username VARCHAR(64) NOT NULL UNIQUE,
      CHECK (CHAR_LENGTH(TRIM(username)) > 0)
    );`)
    
      user_id    serial PRIMARY KEY,
  user_login 		VARCHAR(15) NOT NULL,
  user_pass 		VARCHAR(70) NOT NULL,
  user_name 		VARCHAR(40) NOT NULL,
  user_email 		VARCHAR(40) NOT NULL,
  user_address      VARCHAR(100),
  user_telephone    varchar(20),
  PRIMARY KEY(user_login)
    
    return db
}


/* Randon String generator */
func RandonString(length int) string {
    b := make([]byte, length)
    var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
    charset := "abcdefghijklmnopqrstuvwxyz" +
        "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    for i := range b {
            b[i] = charset[seededRand.Intn(len(charset))]
    }
    return string(b)
}



func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
