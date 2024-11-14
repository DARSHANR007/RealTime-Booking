package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func Book_Ticket(client *redis.Client, ctx context.Context, ticket_id string) {
	availabe, err := client.HGet(ctx, "ticket:"+ticket_id, "available").Result()
	if err != nil {
		fmt.Println("Could not get ticket:", err)
		panic(err)
	}

	if availabe != "true" {
		fmt.Println("Ticket not available")
		return
	}

	err = client.HSet(ctx, "ticket:"+ticket_id, "available", "false").Err()
	if err != nil {
		fmt.Println("Could not book ticket:", err)
		panic(err)
	}

	updatedTicket, err := client.HGetAll(ctx, "ticket:"+ticket_id).Result()
	if err != nil {
		fmt.Println("Error retrieving updated ticket:", err)
		return
	}

	fmt.Println("Ticket booked successfully!")
	fmt.Println("Updated Ticket Details:")
	for field, value := range updatedTicket {
		fmt.Printf("%s: %s\n", field, value)
	}
}
