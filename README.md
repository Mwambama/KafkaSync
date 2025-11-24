KafkaSync

KafkaSync is a fault-tolerant, distributed file synchronization system designed to automate the reliable transfer of files from a remote SFTP server. It uses Apache Kafka as a high-throughput message broker to decouple file detection from processing, ensuring data integrity and scalability.

 Overview

In distributed systems, synchronizing large files reliably is complex. Traditional manual scripts lack automation, retry logic, and auditability.

KafkaSync solves this by creating a robust pipeline:

Ingest: A Producer publishes file transfer jobs to Kafka.

Process: A resilient Go Consumer listens for jobs and executes high-performance parallel downloads using LFTP.

Archive: Files are verified, moved to a completed directory, and automatically uploaded to Cloud Object Storage (S3/MinIO).

Monitor: All activities are logged to a PostgreSQL database and visualized on a real-time React Dashboard.

This project mimics the architecture used in modern data engineering pipelines, demonstrating proficiency in event-driven design, distributed systems, and cross-platform integration (Windows/Go + Linux/WSL).

 Key Features

 Event-Driven Architecture: Decouples the source (Producer) from the destination (Consumer) using Kafka.

 High-Performance Transfers: Leverages LFTP (via WSL) for robust, multi-threaded SFTP downloads (pget/mirror).

 Cloud Archiving: Automatically uploads completed files to S3-compatible storage (AWS/MinIO).

 Structured Metadata: Records every transaction status (Success/Fail) in a PostgreSQL database.

 Real-Time Dashboard: A React + Tailwind CSS interface to monitor the queue, transfer history, and system status.

 Reliability: Includes retry logic for failed downloads and ensures atomic file movement (staging → completed).

🔧 Configurable: Centralized config.toml for managing brokers, credentials, and paths without code changes.

 System Architecture

<!-- <img width="493" height="275" alt="KafkaSync Architecture" src="https://github.com/user-attachments/assets/b40e857e-c822-4051-846f-6ba07554a8f3" /> -->

Producer sends a JSON message ({name, location, hash}) to the kafkasync-files topic.

Consumer parses the message and triggers an LFTP subprocess via WSL.

LFTP downloads the file from the SFTP server to the local ./incompletes staging area.

Upon success, the Consumer:

Moves the file to ./completes.

Uploads it to S3/MinIO.

Logs the result to PostgreSQL.

API Server reads the database and feeds the React Dashboard.

🔧 Getting Started

Prerequisites

Go: v1.21+

Node.js: v18+ (for the Dashboard)

Docker Desktop: For Kafka, Postgres, MinIO, and SFTP containers.

WSL (Ubuntu): Required for LFTP execution on Windows.

# Install LFTP inside your Ubuntu/WSL terminal
sudo apt update && sudo apt install lftp


Installation

Clone the Repository

git clone [https://github.com/your-username/kafkasync.git](https://github.com/your-username/kafkasync.git)
cd kafkasync


Configure the System
Create a config.toml file in the root directory.

kafka_url = "localhost:9094"
num_threads = 4
debug_level = "debug"

[remoteDetails]
host = "localhost:2222"
username = "testuser"
password = "password"

[locations]
incompletes = "./incompletes/"
completes = "./completes/"

[database]
host = "localhost"
port = 5433
user = "ks_user"
password = "ks_password"
dbname = "kafkasync_db"

[objectStorage]
endpoint = "localhost:9000"
access_key = "minioadmin"
secret_key = "minioadmin"
bucket = "kafkasync-archive"
use_ssl = false
region = "us-east-1"


Create Local Directories

mkdir incompletes
mkdir completes


🏃 Running the Application

You will need 4 terminal windows to run the full stack.

1. Infrastructure (Docker)
Start Kafka, Postgres, MinIO, and the test SFTP server.

docker-compose up -d


2. Backend Services
Start the Consumer (Worker) and the API Server.

# Terminal A: Consumer
go run ./cmd/consumer/main.go

# Terminal B: API Server
go run ./cmd/server/main.go


3. Dashboard (Frontend)
Launch the React UI.

cd kafkasync-dashboard
npm run dev


Open http://localhost:5173 in your browser.

4. Producer (Trigger)
Generate a dummy file and trigger a download job.

# Generate a test file in the source folder
.\generate_data.ps1 "test-data.txt"

# Send the job to Kafka
go run ./cmd/producer/main.go


Follow the prompts: Name: test-data.txt, Location: /uploads.

🗺️ Future Roadmap

[ ] Metrics & Monitoring: Integrate Prometheus to export download speeds and queue lag metrics to Grafana.

[ ] Dead Letter Queue: Automatically route permanently failed jobs to a separate Kafka topic for manual inspection.

[ ] Containerization: Dockerize the Producer, Consumer, and API Server for single-command deployment.

📜 License

This project is licensed under the MIT License. See the LICENSE file for details.