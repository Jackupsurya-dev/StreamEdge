
# StreamEdge

StreamEdge is a real-time data streaming system that processes user data efficiently using RabbitMQ, PostgreSQL, Redis, and Server-Sent Events (SSE). It enables seamless ingestion, storage, and live updates of user data.

## 🚀 Features
- 📤 **Producer Service**: Reads a CSV file and publishes user data to RabbitMQ.
- 📥 **Consumer Service**: Consumes messages from RabbitMQ and stores data in PostgreSQL & Redis.
- 📊 **ReactJS UI**: Displays real-time user data using SSE.
- ⚡ **Real-time Processing**: Ensures immediate data updates.
- 🏗 **Scalable & Containerized**: Dockerized for easy deployment.

## 🛠 Tech Stack
- **Backend**: Go (Golang), RabbitMQ, PostgreSQL, Redis
- **Frontend**: React.js, Server-Sent Events (SSE)
- **Infrastructure**: Docker, Docker Compose

## 🔧 Installation & Setup

### 1️⃣ Clone the Repository
```sh
git clone https://github.com/Jackupsurya-dev/StreamEdge.git
cd streamedge
```

## API Documentation 
[Postman Documentation](https://documenter.getpostman.com/view/28116588/2sAYdctYuf)

```
```

# Producer Service

This guide explains how to set up and run the Producer service. The Producer service ingests user data from a CSV file and publishes it to RabbitMQ.

---
## Prerequisites

Ensure that the following tools are installed on your system:
- **Go** ([Installation Guide](https://go.dev/doc/install))
- **Docker** & **Docker Compose** ([Installation Guide](https://docs.docker.com/get-docker/))
- **RabbitMQ** - [Download](https://www.rabbitmq.com/download.html)
- **cURL** (for API testing)

### Install Go (if not installed)

#### Windows
1. Download the Go installer from [Go Downloads](https://go.dev/dl/).
2. Run the installer and follow the setup instructions.
3. Verify the installation:
   ```sh
   go version
   ```

#### macOS (Homebrew)
```sh
brew install go
```

#### Linux (Debian/Ubuntu)
```sh
sudo apt update && sudo apt install -y golang
```

---
## Running the Producer Service Locally

1. navigate to the `producer` directory:
   ```sh
   cd producer
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Start RabbitMQ manually or ensure they are running as services.

4. Set up environment variables (create a `config.json` file in the project directory add the below variables into it):

   ```sh
   {
     "rabbitmq_url": "amqp://guest:guest@localhost:5672/",
     "encryption_key": "thisis32byteencryptionkeyexample"
   }
   ```
4. Start the service:
   ```sh
   go run main.go
   ```
---
## Running the Producer Service with Docker

Set up environment variables (create a `config.json` file in the project directory and add the below variables into it):

   ```sh
   {
      "rabbitmq_url": "amqp://guest:guest@rabbitmq:5672/",
      "encryption_key": "thisis32byteencryptionkeyexample"
   }
   ```

Navigate to the `producer` directory and start the service:

   ```sh
   cd path/to/producer

   # creating a shared network
   docker network create rabbitmq_network

   docker-compose up -d --build
   ```

Verify that the service is running:

```sh
docker ps
```

---

## Stopping and Removing the Producer Container

To stop the producer container:

```sh
docker-compose down
```

To remove all containers, networks, and volumes:

```sh
docker system prune -a --volumes
```

---
## Troubleshooting

- If you get a **port conflict error**, ensure no services are running on the required ports (5672, 15672) and stop them before running Docker.
- If RabbitMQ fails to connect, check logs using:
  
  ```sh
  docker logs rabbitmq-producer
  ```

---
## Additional Notes

- The Producer service publishes messages to RabbitMQ.
- Ensure RabbitMQ is running before testing the Producer API.


```
```
# Consumer Service

This guide explains how to set up and run the Consumer service using Docker or locally. The Consumer service listens to RabbitMQ, consumes user data, and stores it in PostgreSQL and Redis.

---
## Prerequisites

Ensure that the following tools are installed on your system:
- **Docker** & **Docker Compose** ([Installation Guide](https://docs.docker.com/get-docker/))
- **Go** (for running the service locally) - [Download](https://go.dev/doc/install)
- **PostgreSQL** - [Download](https://www.postgresql.org/download/)
- **Redis** - [Download](https://redis.io/docs/getting-started/)
- **RabbitMQ** - [Download](https://www.rabbitmq.com/download.html)
- **cURL** (for API testing)

### Check for Running Containers

Before starting, check if Redis, RabbitMQ, or PostgreSQL are already running:

```sh
docker ps | grep -E "redis|rabbitmq|postgres"
```

If any container is running on ports **6379**, **5672**, **15672**, or **5432**, stop them using:

```sh
docker stop <container_id> && docker rm <container_id>
```

---
## Running the Consumer Service Locally

1. navigate to the consumer directory:

```sh
cd your-repo/consumer
```

2. Install dependencies:

```sh
go mod tidy
```

3. Start PostgreSQL, Redis, and RabbitMQ manually or ensure they are running as services.

4. Set up environment variables (create a `config.json` file in the project directory add the below variables into it):

```sh
  {
    "rabbitmq_url": "amqp://guest:guest@localhost:5672/",
    "redis_address": "localhost:6379",
    "postgresql_host": "localhost",
    "postgresql_port": "5432",
    "postgresql_user": "postgres",
    "postgresql_password": "password",
    "postgresql_dbname": "testdb",
    "encryption_key": "thisis32byteencryptionkeyexample"
  }
```

5. Run the service:

```sh
go run main.go
```

---
## Running the Consumer Service with Docker

Set up environment variables (create a `config.json` file in the project directory add the below variables into it):

```sh
  {
    "rabbitmq_url": "amqp://guest:guest@rabbitmq:5672/",
    "redis_address": "redis:6379",
    "postgresql_host": "postgres",
    "postgresql_port": "5432",
    "postgresql_user": "postgres",
    "postgresql_password": "password",
    "postgresql_dbname": "testdb",
    "encryption_key": "thisis32byteencryptionkeyexample"
  }
```

Navigate to the `consumer` directory and start the service:

```sh
cd path/to/consumer

docker-compose up -d --build
```

Verify that the service is running:

```sh
docker ps
```

---
## Creating the Users Table

After the services are running, create the `users` table in the **PostgreSQL** database.

1. Connect to PostgreSQL inside the container:

```sh
docker exec -it postgres-consumer psql -U postgres -d testdb
```

2. Run the following SQL command:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email_address VARCHAR(255) NOT NULL,
    created_at VARCHAR(255),
    merged_at VARCHAR(255),
    deleted_at VARCHAR(255),
    parent_user_id FLOAT
);
```

3. Exit the PostgreSQL shell:

```sh
\q
```

---

## Stopping and Removing the Consumer Container

To stop the consumer container:

```sh
docker-compose down
```

To remove all containers, networks, and volumes:

```sh
docker system prune -a --volumes
```

---
## Troubleshooting

- If RabbitMQ fails to connect, check logs using:
  
  ```sh
  docker logs rabbitmq-consumer
  ```
- If PostgreSQL fails to connect, verify if it's running:

  ```sh
  docker logs postgres-consumer
  ```

```
```

# React App Dockerization Guide

## Prerequisites
Ensure you have the following installed on your system:
- [Node.js](https://nodejs.org/) (for local development)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Installation and Running Locally

1. **Navigate to the `user-stream-ui` directory and start the service:**
   ```sh
   cd user-stream-ui
   ```

2. **Install Dependencies:**
   ```sh
   npm install
   ```

3. **Start the Development Server:**
   ```sh
   npm start
   ```

The React app will be available at `http://localhost:3000`.


## Building and Running the Docker Container

1. **Build the Docker Image:**
   ```sh
   cd path/to/user-stream-ui

   docker-compose build
   ```

2. **Run the Container:**
   ```sh
   cd path/to/user-stream-ui

   docker-compose up -d
   ```

3. Open `http://localhost:3000` in your browser to access the app.

## Stopping and Removing the Container

To stop the running container:
```sh
docker-compose down
```

To remove all unused images, containers, and volumes:
```sh
docker system prune -a
```

## Troubleshooting
- If port `3000` is already in use, stop the process using it:
  ```sh
  sudo lsof -i :3000
  sudo kill -9 <PID>
  ```
- Check logs of the running container:
  ```sh
  docker-compose logs -f
  ```

Now your React app is successfully containerized and running in Docker! 🚀

---

## Using the Producer API

### Upload a CSV File

Send a CSV file to the producer service to ingest user data into RabbitMQ:

```sh
curl --location 'http://localhost:8080/upload-csv' \
--form 'file=@"/path/to/Users1.csv"'
```

---

## Using the Consumer API

### Retrieve User Data

Get users from the database by first name:

```sh
curl --location 'http://localhost:8081/users?first_name=Grace&last_name=Taylor'
```

---

Happy coding! 🚀

