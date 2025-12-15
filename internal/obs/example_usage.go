// +build ignore

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ironystock/agentic-obs/internal/obs"
)

// This file demonstrates how to use the OBS client wrapper.
// Build and run with: go run example_usage.go

func main() {
	// Create a new OBS client
	config := obs.ConnectionConfig{
		Host:     "localhost",
		Port:     "4455",
		Password: "", // Set if your OBS WebSocket has a password
	}

	client := obs.NewClient(config)

	// Set up event handler
	eventHandler := obs.NewEventHandler(func(eventType obs.EventType, data map[string]interface{}) {
		msg, _ := obs.FormatEventNotification(eventType, data)
		log.Printf("Event: %s", msg)
	})
	client.SetEventCallback(eventHandler)

	// Connect to OBS
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to OBS: %v", err)
	}
	defer client.Close()

	log.Println("Connected to OBS!")

	// Get connection status
	status, err := client.GetConnectionStatus()
	if err != nil {
		log.Fatalf("Failed to get connection status: %v", err)
	}
	fmt.Printf("OBS Version: %s\n", status.OBSVersion)
	fmt.Printf("WebSocket Version: %s\n", status.WebSocketVersion)

	// List all scenes
	scenes, currentScene, err := client.GetSceneList()
	if err != nil {
		log.Fatalf("Failed to get scene list: %v", err)
	}
	fmt.Printf("Current Scene: %s\n", currentScene)
	fmt.Printf("Available Scenes: %v\n", scenes)

	// Get detailed scene information
	if len(scenes) > 0 {
		scene, err := client.GetSceneByName(scenes[0])
		if err != nil {
			log.Fatalf("Failed to get scene details: %v", err)
		}
		fmt.Printf("Scene '%s' has %d sources\n", scene.Name, len(scene.Sources))
	}

	// Get recording status
	recordStatus, err := client.GetRecordingStatus()
	if err != nil {
		log.Fatalf("Failed to get recording status: %v", err)
	}
	fmt.Printf("Recording Active: %v\n", recordStatus.Active)

	// Get streaming status
	streamStatus, err := client.GetStreamingStatus()
	if err != nil {
		log.Fatalf("Failed to get streaming status: %v", err)
	}
	fmt.Printf("Streaming Active: %v\n", streamStatus.Active)

	// Get overall OBS status
	obsStatus, err := client.GetOBSStatus()
	if err != nil {
		log.Fatalf("Failed to get OBS status: %v", err)
	}
	fmt.Printf("FPS: %.2f\n", obsStatus.FPS)
	fmt.Printf("Frame Time: %.2f ms\n", obsStatus.FrameTime)

	// List all audio inputs
	sources, err := client.ListSources()
	if err != nil {
		log.Fatalf("Failed to list sources: %v", err)
	}
	fmt.Printf("Found %d audio/video sources\n", len(sources))

	// Example: Create a new scene
	// err = client.CreateScene("Test Scene")
	// if err != nil {
	// 	log.Printf("Failed to create scene: %v", err)
	// }

	// Example: Switch to a different scene
	// if len(scenes) > 1 {
	// 	err = client.SetCurrentScene(scenes[1])
	// 	if err != nil {
	// 		log.Printf("Failed to switch scene: %v", err)
	// 	}
	// }

	// Example: Start/stop recording
	// err = client.StartRecording()
	// if err != nil {
	// 	log.Printf("Failed to start recording: %v", err)
	// }
	// time.Sleep(5 * time.Second)
	// outputPath, err := client.StopRecording()
	// if err != nil {
	// 	log.Printf("Failed to stop recording: %v", err)
	// } else {
	// 	fmt.Printf("Recording saved to: %s\n", outputPath)
	// }

	// Keep the program running to receive events
	fmt.Println("\nListening for OBS events (press Ctrl+C to exit)...")
	time.Sleep(30 * time.Second)
}
