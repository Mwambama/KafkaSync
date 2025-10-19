# KafkaSync

**KafkaSync** A reliable, event-driven file synchronization system built with Go and Kafka. This project automates the transfer of files from a remote SFTP server by using Kafka as a message broker to trigger downloads managed by a resilient consumer.
## 🚀 Key Features

- 🔔 Event-driven file transfer using **Kafka** (Producer → Consumer architecture)
- ✅ Transfers triggered only when remote file is ready and local server is online
- 📂 Uses **LFTP** for high-performance, resumable downloads
- 🛡️ Ensures files are only committed in Kafka after full and successful download
- 🧪 Fault-tolerant: handles power failures, network interruptions, and incomplete transfers
- 🧰 Shell + system tools integration for simplicity and robustness

## 🛠️ Tech Stack

- **Kafka** – event queue/message broker
- **LFTP** – for efficient file transfers (resumable, scriptable)
- **Shell (Bash)** – for scripting and automation
- **Linux utilities** – inotify, mv, mkdir, etc.
- (Optional: Python version coming soon)

## 📁 Folder Structure

