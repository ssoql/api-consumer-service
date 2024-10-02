package main

import (
	"fmt"
	"runtime"
	"time"

	"api-consumer-service/config"
	"api-consumer-service/internal/app"
)

func main() {
	env := config.NewEnv()
	start := time.Now()
	PrintMemUsage()

	app.Run(env)

	elapsed := time.Since(start)

	PrintMemUsage()
	fmt.Printf("DONE %d posts in %s\n", 0, elapsed.String())
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
