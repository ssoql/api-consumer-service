package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"api-consumer-service/config"
	"api-consumer-service/global"
	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/infrastructure"
	"api-consumer-service/internal/use_cases"
	apiClient "api-consumer-service/pkg/client"
)

func Run(env *config.Env) {
	log.Printf("CPUs: %d", runtime.NumCPU())
	urlSeed := env.ApiSeedUrl
	concurrency := runtime.NumCPU() / 2

	postDispatch, err := infrastructure.CreateDispatcher(env.AppEnv, env.RabbitMqUrl, env.RabbitMqQueueName)
	if err != nil {
		log.Fatal(err)
	}

	client := apiClient.NewApiClient(global.RequestTimeout)
	getTotalUseCase := use_cases.NewGetTotalUseCase(client)
	retryStrategy := infrastructure.NewExponentialBackoff(global.Retry)
	getPostsUseCase := use_cases.NewGetPostsUseCase(client, retryStrategy)
	sendPostsUseCase := use_cases.NewSendPostsUseCase(postDispatch)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	total, err := getTotalUseCase.GetTotal(ctx, fmt.Sprintf(urlSeed, 1, 0))
	if err != nil {
		log.Fatal(err)
	}

	pageChan := make(chan string, total)
	postsChan := make(chan []dto.Post, total)

	createPagesQueue(total, urlSeed, pageChan)

	var wg sync.WaitGroup

	go sendPostsUseCase.Send(ctx, postsChan)
	fetchAllPosts(ctx, &wg, getPostsUseCase, concurrency, pageChan, postsChan)

	wg.Wait()
	close(postsChan)
}

func createPagesQueue(total int, urlSeed string, pageChan chan<- string) {
	totalPages := total / global.ApiReadLimit
	if total%global.ApiReadLimit > 0 {
		totalPages++
	}
	fmt.Printf("Total results: %d\n", total)
	fmt.Printf("Total Pages mod: %d\n", totalPages)

	// Enqueue all pages to be processed
	for page := 1; page <= totalPages; page++ {
		pageChan <- fmt.Sprintf(urlSeed, global.ApiReadLimit, (page-1)*global.ApiReadLimit)
	}
	close(pageChan)
}

func fetchAllPosts(ctx context.Context, wg *sync.WaitGroup, useCase *use_cases.GetPostsUseCase, jobsCounter int, pageChan <-chan string,
	postsChan chan<- []dto.Post) {
	for i := 0; i < jobsCounter; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for page := range pageChan {
				fmt.Printf("Worker %d fetching page %s\n", id, page)
				start := time.Now()

				posts, err := useCase.GetPosts(ctx, page)
				if err != nil {
					fmt.Printf("Error fetching posts: %s\n", err)
					if errors.Is(err, context.Canceled) {
						return
					}
				}
				elapsed := time.Since(start)

				postsChan <- posts
				fmt.Printf("Worker %d finished page %s in %s\n", id, page, elapsed.String())
				//time.Sleep(10 * time.Millisecond)
				time.Sleep(3 * time.Second)
			}
		}(i)
	}
}
