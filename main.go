package main

import (
	"database/sql"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var db *sql.DB

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("connectDB:", err)
		return
	}
	defer db.Close()

	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	r := newRoom()
	http.Handle("/chat-ws/room", r)

	go r.run()

	log.Println("Starting Chat Server. Port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func connectDB() (*sql.DB, error) {
	return sql.Open("postgres", "host=/var/run/postgresql/ dbname=chat sslmode=disable")
}
