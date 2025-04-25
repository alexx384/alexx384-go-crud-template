package main

import (
	"crud/internal"
	"crud/internal/config"
	"crud/internal/config/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	dbPool, err := database.NewPool(appConfig.DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	app := gin.Default()
	internal.SetupRouter(dbPool, app)

	err = app.Run(":8080")
	if err != nil {
		fmt.Println("Something went wrong")
		return
	}
}
