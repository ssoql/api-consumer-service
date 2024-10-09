package main

import (
	"fmt"
	"time"

	"api-consumer-service/config"
	"api-consumer-service/internal/app"
	"api-consumer-service/pkg/runtime_helper"
)

func main() {
	env := config.NewEnv()
	start := time.Now()

	runtime_helper.PrintMemUsage()

	app.Run(env)

	elapsed := time.Since(start)

	runtime_helper.PrintMemUsage()
	fmt.Printf("DONE %d posts in %s\n", 0, elapsed.String())
}
