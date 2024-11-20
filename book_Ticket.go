package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func Book_Ticket(client *redis.Client, ctx context.Context, ticket_id string) error {
	availabe, err := client.HGet(ctx, "ticket:"+ticket_id, "available").Result()
	if err != nil {
		fmt.Println("Could not get ticket:", err)
		return err
	}

	if availabe != "true" {
		fmt.Printf("Ticket %s not available\n", ticket_id)
		return fmt.Errorf("ticket %s is not available", ticket_id)
	}

	// Attempt to book the ticket by setting it to unavailable
	err = client.HSet(ctx, "ticket:"+ticket_id, "available", "false").Err()
	if err != nil {
		fmt.Println("Could not book ticket:", err)
		return err
	}

	// Retrieve updated ticket details
	updatedTicket, err := client.HGetAll(ctx, "ticket:"+ticket_id).Result()
	if err != nil {
		fmt.Println("Error retrieving updated ticket:", err)
		return err
	}

	fmt.Printf("Ticket %s booked successfully!\n", ticket_id)
	fmt.Println("Updated Ticket Details:")
	for field, value := range updatedTicket {
		fmt.Printf("%s: %s\n", field, value)
	}

	return nil
}
