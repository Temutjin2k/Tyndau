package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("❌ Failed to connect to NATS:", err)
	}
	defer nc.Drain()

	timestamp := time.Now().UTC()

	events := []struct {
		Subject string
		Event   map[string]interface{}
	}{
		{
			"tyndau.user_registered",
			map[string]interface{}{
				"event_type": "user.registered",
				"user_id":    "user-001",
				"email":      "tamutdzhin2006@mail.ru",
				"timestamp":  timestamp,
				"data": map[string]interface{}{
					"name": "Beibars",
				},
			},
		},
	}

	for _, e := range events {
		payload, err := json.Marshal(e.Event)
		if err != nil {
			log.Println("❌ Failed to marshal event:", e.Event["event_type"])
			continue
		}
		err = nc.Publish(e.Subject, payload)
		if err != nil {
			log.Println("❌ Failed to publish:", e.Subject)
		} else {
			fmt.Printf("✅ Sent event: %s → %s\n", e.Event["event_type"], e.Subject)
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("🏁 All events sent.")
}
