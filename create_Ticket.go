package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func create_Ticket(client *redis.Client, ctx context.Context) {
	tickets := []map[string]string{
		{
			"ticket_id":  "TICKET123",
			"event_name": "Concert A",
			"date":       "2024-11-25",
			"price":      "1500",
			"available":  "true",
		},
		{
			"ticket_id":  "TICKET124",
			"event_name": "Concert B",
			"date":       "2024-12-05",
			"price":      "2000",
			"available":  "false",
		},
		{
			"ticket_id":  "TICKET125",
			"event_name": "Concert C",
			"date":       "2024-12-15",
			"price":      "1800",
			"available":  "true",
		},
	}

	for _, ticket := range tickets {
		hashFields := []string{
			"ticket_id", ticket["ticket_id"],
			"event_name", ticket["event_name"],
			"date", ticket["date"],
			"price", ticket["price"],
			"available", ticket["available"],
		}
		err := client.HSet(ctx, "ticket:"+ticket["ticket_id"], hashFields).Err()
		if err != nil {
			fmt.Println("count not create ticket:", err)
			panic(err)
			return
		}

	}

}
