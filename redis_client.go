package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
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

}
