package main

import (
	"github.com/kwa0x2/Settle-Backend/internal/app"
	"log"
)

func main() {
	application := app.NewApp()
	application.SetupRoutes()

	if err := application.Run(); err != nil {
		log.Fatal("failed to run app: ", err)
	}
}
