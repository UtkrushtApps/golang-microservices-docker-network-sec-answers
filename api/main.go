package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/go-redis/redis/v8"
	"context"
)

var db *sql.DB
var rdb *redis.Client

func main() {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := "appdb"

	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	var err error
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Could not reach Postgres: %v", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Could not reach Redis: %v", err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/cache", cacheHandler)

	log.Println("API server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	row := db.QueryRow("SELECT NOW()")
	var ts string
	if err := row.Scan(&ts); err != nil {
		http.Error(w, "db error", 500)
		return
	}
	fmt.Fprintf(w, "Database time: %s\n", ts)
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	key := "test:key"
	ctx := context.Background()
	err := rdb.Set(ctx, key, "hello-cache", 0).Err()
	if err != nil {
		http.Error(w, "redis set error", 500)
		return
	}
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		http.Error(w, "redis get error", 500)
		return
	}
	fmt.Fprintf(w, "Redis value: %s\n", val)
}
