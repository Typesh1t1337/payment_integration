package main

import (
	"context"
	"log"
	"payment_integration/internal/config"
	"payment_integration/internal/uow"
)

func main() {
	cfg := config.MustLoad()
	pool, err := config.NewDB(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	uow := uow.NewSQLUoW(pool)
	
	
}