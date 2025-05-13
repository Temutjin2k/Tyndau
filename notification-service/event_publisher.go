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
		{
			"tyndau.user_registered",
			map[string]interface{}{
				"event_type": "user.registered",
				"user_id":    "user-001",
				"email":      "beibarys7ergaliev@gmail.com",
				"timestamp":  timestamp,
				"data": map[string]interface{}{
					"name": "Beibars",
				},
			},
		},
		{
			"tyndau.user_registered",
			map[string]interface{}{
				"event_type": "user.registered",
				"user_id":    "user-002",
				"email":      "bakhytzhanabdilmazhit@gmail.com",
				"timestamp":  timestamp,
				"data": map[string]interface{}{
					"name": "Bakhytzhan",
				},
			},
		},
		{
			"tyndau.user_registered",
			map[string]interface{}{
				"event_type": "user.registered",
				"user_id":    "user-003",
				"email":      "mudrec6putei228@mail.ru",
				"timestamp":  timestamp,
				"data": map[string]interface{}{
					"name": "Mudrec",
				},
			},
		},
		{
			"tyndau.user_registered",
			map[string]interface{}{
				"event_type": "user.registered",
				"user_id":    "user-004",
				"email":      "230311@astanait.edu.kz",
				"timestamp":  timestamp,
				"data": map[string]interface{}{
					"name": "AITU Student",
				},
			},
		},
		{
			"tyndau.album_released",
			map[string]interface{}{
				"event_type": "music.album_released",
				"user_id":    "user-004",
				"email":      "230311@astanait.edu.kz",
				"timestamp":  timestamp,
				"data": map[string]interface{}{
					"album_name":  "Dreams of AITU",
					"artist_name": "Synthfox",
				},
			},
		},
		{
			"tyndau.user_registered",
			map[string]interface{}{
				"event_type": "test.created",
				"user_id":    "test-001",
				"email":      "test@example.com",
				"timestamp":  timestamp,
				"data": map[string]interface{}{
					"description": "This is a test event just for fun 🎉",
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
