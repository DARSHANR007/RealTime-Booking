package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
)

func generateTicketData(id int) map[string]string {
	return map[string]string{
		"ticket_id":  fmt.Sprintf("TICKET%d", id),
		"event_name": fmt.Sprintf("Event-%d", rand.Intn(1000)),
		"date":       fmt.Sprintf("2024-%02d-%02d", rand.Intn(12)+1, rand.Intn(28)+1),
		"price":      fmt.Sprintf("%d", rand.Intn(5000)+500),
		"available":  fmt.Sprintf("%t", rand.Intn(2) == 0),
	}
}

func createTicket(client *redis.Client, ctx context.Context, ticket map[string]string) error {
	hashFields := []string{
		"ticket_id", ticket["ticket_id"],
		"event_name", ticket["event_name"],
		"date", ticket["date"],
		"price", ticket["price"],
		"available", ticket["available"],
	}

	return client.HSet(ctx, "ticket:"+ticket["ticket_id"], hashFields).Err()
}
