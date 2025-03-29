package main

import (
	db2 "github.com/lipusipu44/Social/internal/db"
	"github.com/lipusipu44/Social/internal/env"
	"github.com/lipusipu44/Social/internal/store"
	"log"
)

func main() {
	addr :=
		env.GetEnv("DB_ADDR", "postgresql://admin:adminpassword@localhost/"+
			"socialnetwork?sslmode=disable")
	dbConn, err := db2.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()
	storage := store.NewStorage(dbConn)
	db2.Seed(storage, dbConn)
}
