package main

import (
	"database/sql"
	"log/slog"
	"os"
)

type svc struct {
}

func (s *svc) GetConfig(info map[string]string) ([]byte, error) {
	return nil, nil
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("Server started")
	repo := struct {
		db *sql.DB
	}{
		db: &sql.DB{},
	}
	db := *sql.DB

}

// http.ListenAndServe(":8080", router.Handler())
// fmt.Println("Server listen on port :8080")
