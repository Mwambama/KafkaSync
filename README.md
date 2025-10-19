# KafkaSync

**KafkaSync** A reliable, event-driven file synchronization system built with Go and Kafka. This project automates the transfer of files from a remote SFTP server by using Kafka as a message broker to trigger downloads managed by a resilient consumer.

## 📖 Overview
In distributed systems, synchronizing large files reliably can be complex. 
Manual SFTP or SCP scripts lack automation, retry logic, and scalability. 
KafkaSync solves this by creating a robust pipeline where file transfer jobs are published as messages to a Kafka topic. 
A Go-based consumer listens for these jobs, executes high-performance downloads using LFTP, and ensures files are moved to their final destination only after a successful transfer.

This project copies the architecture used in modern data engineering and DevOps for data ingestion and automation, showing proficiency in event-driven architecture, distributed systems, and cross-platform integration (Windows/Go + Linux/WSL).

## 🚀 Key Features
📦 Event-Driven Architecture:
Uses a Producer/Consumer model with Kafka to separate file detection from the download process.

📨 Structured Messaging: 
Sends download jobs as structured JSON messages, carrying metadata like filename, remote path, and hash.

⚙️ High-Performance Transfers: 
Leverages LFTP (via WSL) for robust and parallelized SFTP downloads, supporting commands like pget and mirror.

🛡️ Fault-Tolerant Consumer: 
Includes retry logic for failed downloads and ensures files are moved from an incompletes to a completes directory only upon success.

🔧 Configurable: 
All settings, including Kafka brokers, SFTP credentials, and file paths, are managed externally in a config.toml file.

📝 Detailed Logging: 
Logs all actions (job reception, download attempts, successes, and failures) to both the console and a consumer.log file for traceability.

## 🛠️ System Architecture & Workflow
The system is made of two main applications that communicate -> Kafka topic.

<img width="493" height="275" alt="image" src="https://github.com/user-attachments/assets/b40e857e-c822-4051-846f-6ba07554a8f3" />

The Producer sends a JSON message to the kafkasync-files topic.

The Consumer receives the message and parses the JSON.

It constructs an lftp command and executes it via WSL.

LFTP downloads the file from the remote SFTP server into the local ./incompletes directory.

If the download succeeds, the consumer moves the file to the ./completes directory.


## 🔧 Getting Started
Follow these instructions for KafkaSync.

Go: Version 1.21 or later.

Docker Desktop: To run the Kafka container.

WSL: An Ubuntu distribution(recommended).

LFTP: installed inside your WSL distribution.


# Run this inside your Ubuntu/WSL terminal
sudo apt update && sudo apt install lftp
Installation & Setup
Clone the Repository

git clone https://github.com/your-username/kafkasync.git
cd kafkasync
Configure the System Create a config.toml file in the root directory and populate it with your settings.
Example:

Ini, TOML

# config.toml

kafka_url = "localhost:9094"
num_threads = 4
debug_level = "debug"

[remoteDetails]
host = "your-sftp-server.com"
username = "your-sftp-user"
password = "your-sftp-password"

[locations]
incompletes = "./incompletes/"
completes = "./completes/"

Create Download Directories Manually create the directories specified in your config.toml. 
The consumer will also create them if they are missing.

PowerShell

mkdir incompletes
mkdir completes
Running the Application
You will need three separate terminals open to run the full system.

Terminal 1: Start Kafka Use Docker Compose to launch the Kafka broker in KRaft mode.

docker-compose up -d
Verify it's running with docker ps.

Terminal 2: Run the Consumer The consumer will connect to Kafka and wait for download jobs.


go run ./cmd/consumer/main.go
You should see the output: ✅ Kafka consumer is now listening for messages...

Terminal 3: Run the Producer The producer will prompt you to enter file details to create a download job.

go run ./cmd/producer/main.go
Follow the prompts to send a message.

🔢 Enter hash: 12345abcdef
📄 Enter file name: my-test-file.mkv
🌐 Enter remote location: /path/to/remote/files/



## 🗺️ Future Roadmap
I will update kafkaSync v2 features:

[ ] 📊 Web Dashboard: A React + Tailwind CSS interface to visualize the queue of pending downloads, see transfer history, and monitor successes or failures in real-time.

[ ] 🗂️ Metadata Database: Replace file-based logging with a SQLite or PostgreSQL database to store structured metadata for every transfer (timestamp, file size, status, retry count).

[ ] ☁️ Cloud Storage Integration: Add a post-download step to automatically upload completed files to a cloud storage provider like Amazon S3.

[ ] 🔄 Failed Job Queue: Implement a dead-letter queue in Kafka to automatically retry failed downloads after a certain delay.

[ ] 📈 Metrics and Monitoring: Integrate Prometheus to export metrics on download speeds, failure rates, and queue depth, with a Grafana dashboard for visualization.

## 📜 License
This project is licensed under the MIT License. See the LICENSE file for details.


