package env

import (
	"os"
	"strconv"
	"log"
	"go-short/storage"
)

type Env struct {
	S storage.Storage
}

func GetEnv() *Env {

	addr := os.Getenv("APP_REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	passwd := os.Getenv("APP_REDIS_PASSWD")
	if passwd == "" {
		passwd = ""
	}
	dbs := os.Getenv("APP_REDIS_DB")

	if dbs == "" {
		dbs = "0"
	}

	db, err := strconv.Atoi(dbs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connect to redis (addr: %s password: %s db: %d)", addr, passwd, db)

	r := storage.NewRedisClient(addr, passwd, db)
	return &Env{S: r}
}
