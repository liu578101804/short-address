package app

import (
	"os"
	"strconv"
	"log"
)

type Env struct {
	S Storage
}

func GetEnv() *Env  {
	var(
		addr,passwd,dbS string
		db int
		err error
		r *RedisCli
	)
	if addr = os.Getenv("APP_REDIS_ADDR");addr == ""{
		addr = "localhost:6379"
	}
	if passwd = os.Getenv("APP_REDIS_PASSED");passwd == ""{
		passwd = ""
	}
	if dbS = os.Getenv("APP_REDIS_DB");dbS == ""{
		dbS = "0"
	}
	if db,err = strconv.Atoi(dbS);err != nil{
		log.Fatal(err)
	}
	log.Printf("connect to redis (addr: %s password: %s db: %d)", addr, passwd, db)
	r = NewRedisCli(addr, passwd, db)
	return &Env{S: r}
}