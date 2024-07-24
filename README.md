## Chat Application
This is a simple chat application built with Go, Gin, WebSockets, Kafka, and a PostgreSQL database. The application allows users to log in, select another user to chat with, and exchange messages in real-time.

## Features
* User login
* Real-time messaging with WebSockets
* Chat message persistence with PostgreSQL
* Message status indicators (single and double ticks)

## Prerequisites
* Go 1.16 or later
* MySQL
* Kafka

## Installation

`1.` Clone the repository:

```
git clone https://github.com/your-username/chat-app.git
cd chat-app
```

`2.` Install dependencies:

```
go mod download
```

`3.` Set up MySQL:

  Create a database and a table for storing messages. You can use the following SQL script as an example:
```
CREATE DATABASE chat_app;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    message_text TEXT NOT NULL,
    status VARCHAR(10) NOT NULL
);

```

`4.` Set up Kafka:
Ensure Kafka is running and create a topic for the chat messages:
```
kafka-topics --create --topic chat-messages --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```

`5.` Configure the application:
Create a `config.yml` file in the root directory with the following content:

```
kafka:
  url: "localhost:9092"
  topic: "chat-messages"

postgres:
  user: "your-db-user"
  password: "your-db-password"
  dbname: "chat_app"
  host: "localhost"
  port: 5432

```

  ## Running the Application
  `1.` Start the Kafka consumer:
  ```
  go run main.go
  ```
  `2.` Run the web server:
  ```
  go run main.go
  ```
  The application will be available at http://localhost:8080

  ## Usage
  `1.` Login:
   * Open http://localhost:8080/login in your browser.
   * Enter a username and password to log in.
     
  `2.` Chat:
   * Select a user from the user list.
   * Type a message in the input box and click "Send".
   * Messages will be displayed in the chat window with status indicators.


