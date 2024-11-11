package main

import (
	"github.com/roxxers/surfe-techtest/internal/adapters/primary"
	"github.com/roxxers/surfe-techtest/internal/adapters/secondary"
	"github.com/roxxers/surfe-techtest/internal/core/services"
)

func main() {
	db := secondary.NewMemoryDatabase()
	s := services.NewService(db)
	http := primary.NewHTTPServer(db, s)
	http.Serve("127.0.0.1:8080") // Replace with configurable variable in future
}
