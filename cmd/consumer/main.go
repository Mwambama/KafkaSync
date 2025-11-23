package main

import (
	"context"
	"database/sql"
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
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/segmentio/kafka-go"
)

type DownloadNotification struct {
	Hash     string `json:"info_hash"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type tomlConfig struct {
	KafkaUrl      string         `toml:"kafka_url"`
	NumThreads    int            `toml:"num_threads"`
	DebugLevel    string         `toml:"debug_level"`
	RemoteDetails remoteDetails  `toml:"remoteDetails"`
	Locations     locations      `toml:"locations"`
	Database      databaseConfig `toml:"database"`
	ObjectStorage s3Config       `toml:"objectStorage"` // ✅ New S3 config
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

type databaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

type s3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	Region    string
}

var conf tomlConfig
var db *sql.DB
var minioClient *minio.Client // ✅ Global S3 Client

func init() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
}

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func initDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.DbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("❌ Database unreachable: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL database")

	query := `
	CREATE TABLE IF NOT EXISTS downloads (
		id SERIAL PRIMARY KEY,
		filename TEXT NOT NULL,
		remote_location TEXT,
		hash TEXT,
		status TEXT,
		downloaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("❌ Failed to create table: %v", err)
	}
}

// ✅ Initialize MinIO/S3
func initS3() {
	var err error
	minioClient, err = minio.New(conf.ObjectStorage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.ObjectStorage.AccessKey, conf.ObjectStorage.SecretKey, ""),
		Secure: conf.ObjectStorage.UseSSL,
	})
	if err != nil {
		log.Fatalf("❌ Failed to create S3 client: %v", err)
	}

	// Check connection by checking/creating bucket
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, conf.ObjectStorage.Bucket)
	if err != nil {
		log.Fatalf("❌ Failed to connect to S3/MinIO: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, conf.ObjectStorage.Bucket, minio.MakeBucketOptions{Region: conf.ObjectStorage.Region})
		if err != nil {
			log.Fatalf("❌ Failed to create bucket: %v", err)
		}
		log.Printf("✅ Created new bucket: %s", conf.ObjectStorage.Bucket)
	} else {
		log.Printf("✅ Connected to S3 Bucket: %s", conf.ObjectStorage.Bucket)
	}
}

// ✅ Upload file to S3
func uploadToStorage(filePath string, filename string) error {
	ctx := context.Background()
	contentType := "application/octet-stream"

	// Upload the file
	info, err := minioClient.FPutObject(ctx, conf.ObjectStorage.Bucket, filename, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	log.Printf("☁️  Successfully uploaded %s to cloud (Size: %d bytes)", filename, info.Size)
	return nil
}

func recordDownload(notification DownloadNotification, status string) {
	query := `INSERT INTO downloads (filename, remote_location, hash, status) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, notification.Name, notification.Location, notification.Hash, status)
	if err != nil {
		log.Printf("⚠️ Failed to log to DB: %v", err)
	} else {
		log.Println("🗂️  Download recorded in database")
	}
}

func main() {
	logFile, err := os.OpenFile("consumer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)

	if err := ensureDir(conf.Locations.Incompletes); err != nil {
		log.Fatalf("❌ Failed to create incompletes directory: %v", err)
	}
	if err := ensureDir(conf.Locations.Completes); err != nil {
		log.Fatalf("❌ Failed to create completes directory: %v", err)
	}

	initDB()
	initS3() // Connect to Cloud

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

		remoteCommand := genRemoteCommand(notification.Location, notification.Name)
		cmd := exec.Command("wsl.exe", "lftp", "-e", remoteCommand)
		cmd.Dir = conf.Locations.Incompletes

		if conf.DebugLevel == "debug" {
			log.Printf("🛠 Executing command: %s", cmd.String())
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		log.Printf("🚀 Running download...")

		if err := cmd.Run(); err != nil {
			log.Printf("❌ Download failed for %s: %v", notification.Name, err)
			recordDownload(notification, "FAILED")
			continue
		}

		from := filepath.Join(conf.Locations.Incompletes, notification.Name)
		if _, err := os.Stat(from); os.IsNotExist(err) {
			log.Printf("❌ File not found after download: %s", from)
			recordDownload(notification, "MISSING")
			continue
		}

		to := filepath.Join(conf.Locations.Completes, notification.Name)
		if err = os.Rename(from, to); err != nil {
			log.Printf("❌ Failed to move %s to completes: %v", notification.Name, err)
			recordDownload(notification, "MOVE_FAILED")
			continue
		}
		log.Printf("✅ File moved to completed: %s", to)

		// Upload to Cloud
		err = uploadToStorage(to, notification.Name)
		if err != nil {
			log.Printf("❌ Failed to upload to S3: %v", err)
			recordDownload(notification, "UPLOAD_FAILED")
		} else {
			recordDownload(notification, "COMPLETED_AND_UPLOADED")
		}

		log.Printf("📨 Message committed for %s", notification.Name)
	}
}

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

	lftpCommand := fmt.Sprintf("pget -n %d -c %s", conf.NumThreads, fullRemote)
	return fmt.Sprintf("set sftp:auto-confirm yes; %s; bye", lftpCommand)
}
