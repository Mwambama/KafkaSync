# KafkaSync

**KafkaSync** is a fault-tolerant, Kafka-based file synchronization and download system designed to ensure safe and complete transfers between a remote and local server. It is ideal for scenarios where data integrity is critical, such as edge computing, telemetry syncing, or CDN replication.

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

