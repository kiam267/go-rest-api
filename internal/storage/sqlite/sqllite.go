package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/kiam267/student-api/internal/config"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
    Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
  db, err := sql.Open("sqlite", cfg.StoragePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open SQLite DB: %w", err)
    }

    // ensure DB is reachable
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping SQLite: %w", err)
    }

    // SQLite recommended settings
    db.SetMaxOpenConns(1)
    db.SetMaxIdleConns(1)

    // enable foreign keys
    if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
        return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
    }

    // create table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS student (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            email TEXT,
            age INTEGER
        )
    `)

    if err != nil {
        return nil, fmt.Errorf("failed to create student table: %w", err)
    }

    return &Sqlite{Db: db}, nil
}

func(s *Sqlite) CreateStudent(name string, email string, age int)(int64, error)  {
	stmt , err := s.Db.Prepare(`
	 INSERT INTO student(name, email, age) VALUES (?, ?, ?)
	
	`)
	if  err != nil{
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name,email,age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil{
		return 0, err
	}


	return lastId , nil
 
}

func (s *Sqlite) Close() error {
    return s.Db.Close()
}
