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
		log.Fatal(err)
	}
	defer nc.Drain()

	timestamp := time.Now().UTC()

	// --- 1. user_registered ---
	userEvent := map[string]interface{}{
		"event_type": "user.registered",
		"user_id":    "user-123",
		"email":      "230311@astanait.edu.kz", // Обновленная почта
		"timestamp":  timestamp,
		"data": map[string]interface{}{
			"name": "Test User",
		},
	}
	userData, _ := json.Marshal(userEvent)
	err = nc.Publish("tyndau.user_registered", userData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Sent user_registered event")

	time.Sleep(1 * time.Second)

	// --- 2. album_released ---
	albumEvent := map[string]interface{}{
		"event_type": "music.album_released",
		"user_id":    "user-123",
		"email":      "230311@astanait.edu.kz", // Обновленная почта
		"timestamp":  timestamp,
		"data": map[string]interface{}{
			"album_name":  "The Beatles",
			"artist_name": "Blackbirds",
		},
	}
	albumData, _ := json.Marshal(albumEvent)
	err = nc.Publish("tyndau.album_released", albumData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Sent album_released event")
}
