package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/segmentio/kafka-go"
)

// 📦 DownloadNotification defines the JSON structure expected from Kafka
type DownloadNotification struct {
	Hash     string `json:"info_hash"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// 📁 Config structs
type tomlConfig struct {
	KafkaUrl      string `toml:"kafka_url"`
	NumThreads    int    `toml:"num_threads"`
	DebugLevel    string `toml:"debug_level"`
	RemoteDetails remoteDetails
	Locations     locations
}

type remoteDetails struct {
	Host     string
	Username string
	Password string
}

type locations struct {
	Incompletes string
	Completes   string
}

var conf tomlConfig

func init() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
}

// 📁 ensureDir creates the directory if it doesn't exist
func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func main() {
	// 🔧 Setup logging
	logFile, err := os.OpenFile("consumer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)

	// 🔧 Ensure both incompletes and completes directories exist at startup
	if err := ensureDir(conf.Locations.Incompletes); err != nil {
		log.Fatalf("❌ Failed to create incompletes directory: %v", err)
	}
	if err := ensureDir(conf.Locations.Completes); err != nil {
		log.Fatalf("❌ Failed to create completes directory: %v", err)
	}

	// ⚙️ Kafka configuration
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{conf.KafkaUrl},
		Topic:          "kafkasync-files",
		GroupID:        "file-consumer-group",
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})
	defer reader.Close()

	log.Println("✅ Kafka consumer is now listening for messages...")

	for {
		ctx := context.TODO()
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Error reading message: %v", err)
			continue
		}

		var notification DownloadNotification
		if err := json.Unmarshal(message.Value, &notification); err != nil {
			log.Printf("❌ Failed to parse JSON message: %v\n", err)
			continue
		}

		log.Printf("⬇️  Preparing to download: %+v\n", notification)

		// ⚙️ Build remote command
		remoteCommand := genRemoteCommand(notification.Location, notification.Name)

		// Use "-e" instead of "-c" to allow command chaining
		cmd := exec.Command("wsl.exe", "lftp", "-e", remoteCommand)

		cmd.Dir = conf.Locations.Incompletes

		// 🔧 Capture and show output
		if conf.DebugLevel == "debug" {
			log.Printf("🛠 Executing command: %s", cmd.String())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		log.Printf("🚀 Running download...")

		if err := cmd.Run(); err != nil {
			log.Printf("❌ Download failed for %s: %v", notification.Name, err)
			continue
		}

		// Move from incompletes to completes using safe path handling
		from := filepath.Join(conf.Locations.Incompletes, notification.Name)

		// Confirm file exists before moving
		if _, err := os.Stat(from); os.IsNotExist(err) {
			log.Printf("❌ File not found after download: %s", from)
			continue
		}

		to := filepath.Join(conf.Locations.Completes, notification.Name)

		if err = os.Rename(from, to); err != nil {
			log.Printf("❌ Failed to move %s to completes: %v", notification.Name, err)
			continue
		}
		log.Printf("✅ File moved to completed: %s", to)

		log.Printf("📨 Message committed for %s", notification.Name)
	}
}

// 🔧 Generate LFTP command
func genRemoteCommand(location, name string) string {
	safeName := strings.ReplaceAll(name, " ", "\\ ")
	safeName = strings.ReplaceAll(safeName, "'", "\\'")

	fullRemote := fmt.Sprintf("sftp://%s:%s@%s%s/%s",
		conf.RemoteDetails.Username,
		conf.RemoteDetails.Password,
		conf.RemoteDetails.Host,
		location,
		safeName,
	)

	var lftpCommand string
	if strings.HasSuffix(name, ".mkv") {
		lftpCommand = fmt.Sprintf("pget -n %d -c %s", conf.NumThreads, fullRemote)
	} else {
		lftpCommand = fmt.Sprintf("mirror --use-pget-n=%d -p -c %s", conf.NumThreads, fullRemote)
	}

	// Prepend 'set sftp:auto-confirm yes;' to fix host key errors
	return fmt.Sprintf("set sftp:auto-confirm yes; %s; bye", lftpCommand)
}
