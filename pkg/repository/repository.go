package repository

import (
	"database/sql"
	elpb "github.com/slntopp/nocloud-proto/events_logging"
)

type SqliteRepository struct {
	*sql.DB
}

func NewSqliteRepository(datasource string) *SqliteRepository {
	db, err := sql.Open("sqlite", datasource)
	if err != nil {
		return nil
	}

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS EVENTS (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    ENTITY TEXT,
    UUID TEXT,
    SCOPE TEXT,
    ACTION TEXT,
    RC INTEGER,
    REQUESTOR TEXT
);

CREATE TABLE IF NOT EXISTS SNAPSHOTS (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    DIFF TEXT,
    EVENT_ID INTEGER,

    FOREIGN KEY (EVENT_ID) REFERENCES EVENTS(ID) ON DELETE CASCADE
);
`)
	if err != nil {
		return nil
	}
	return &SqliteRepository{db}
}

func (r *SqliteRepository) CreateEvent(event *elpb.Event) error {
	return nil
}
