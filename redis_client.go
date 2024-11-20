package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

func get_Ticket() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Could not connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

	auth(client, ctx)

	numofTickets := 100000
	numofWorkers := 100

	jobs := make(chan map[string]string, 1000)

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for ticket := range jobs {
			if err := createTicket(client, ctx, ticket); err != nil {
				fmt.Printf("Failed to create ticket: %s\n", err)
			}
		}

	}

	for x := 0; x < numofWorkers; x++ {
		wg.Add(1)
		go worker()
	}
	start := time.Now()

	for i := 0; i < numofTickets; i++ {
		jobs <- generateTicketData(i)
	}
	close(jobs)
	wg.Wait()

	fmt.Printf("Created %d tickets in %s\n", numofTickets, time.Since(start))

}
