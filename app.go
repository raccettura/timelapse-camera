package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron/v3"
)

type Config struct {
	StreamURL string `json:"stream_url"`
	Interval  string `json:"interval"`
	OutputDir string `json:"output_dir"`
}

func main() {
	configFile := flag.String("config", "config.json", "Path to the config file")
	flag.Parse()

	log.Printf("Starting timelapse camera application")
	log.Printf("Using config file: %s", *configFile)

	// Read the config file
	file, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	log.Printf("Successfully read config file")

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}

	log.Printf("Configuration loaded:")
	log.Printf("  Stream URL: %s", config.StreamURL)
	log.Printf("  Interval: %s", config.Interval)
	log.Printf("  Output Dir: %s", config.OutputDir)

	// Create a cron job
	c := cron.New()
	_, err = c.AddFunc(config.Interval, func() {
		log.Printf("Starting capture at %s", time.Now().Format(time.RFC3339))

		// Build the ffmpeg command
		timestamp := time.Now().Format("20060102150405")
		tempFile := fmt.Sprintf("%s/%s.%s", config.OutputDir, timestamp, "png")
		cmd := exec.Command("ffmpeg", "-y", "-i", config.StreamURL, "-frames:v", "1", tempFile)
		log.Printf("Running ffmpeg command: %s", cmd.String())

		// Run the ffmpeg command
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to run ffmpeg: %v\nOutput: %s", err, string(output))
			return
		}
		log.Printf("ffmpeg completed successfully")
		log.Printf("ffmpeg output: %s", string(output))

		log.Printf("Capture completed at %s", time.Now().Format(time.RFC3339))
	})

	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	log.Printf("Cron job added successfully with interval: %s", config.Interval)

	// Start the cron job
	c.Start()
	log.Printf("Cron scheduler started")

	// Keep the program running
	log.Printf("Application running. Press Ctrl+C to stop.")
	for {
		time.Sleep(time.Hour)
	}
}
