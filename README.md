# 💬 Chat Application

This is a chat application where the backend is written in Golang, the frontend is created in React, and MongoDB is used for data storage. Both the backend and frontend are dockerized into a Docker container for easy deployment.

## ✨ Features

- Real-time Chat: Utilizing Websockets for instant messaging.
- User Authentication: Secured via JWT for user access.
- Chat Room Management: Users can create and manage chat rooms.

## 🛠️ Technologies Used

- Backend: Golang
- Frontend: React
- Database: MongoDB
- Docker

## 🚀 Getting Started

### 📋 Prerequisites

#### For Docker
- Docker installed on your machine.
- Create a `.env` file in root directory based on `.env.example` and fill in the environment variable.

#### For Local
- Golang installed on your machine.
- Node.js and npm installed on your machine.
- MongoDB installed and running on your machine.
- Create a .env file in the root directories of both the backend and frontend based on their respective .env.example files and fill in the environment variables.

### 🏃‍♂️ Running the Application

1. Clone this repository:

    ```bash
    git clone https://github.com/your/repository.git
    ```
2. Navigate to the project directory:

    ```bash 
    cd repository
    ```
#### Using Docker

1. Create a .env file in the root directories of both the backend and frontend based on their respective .env.example files and fill in the environment variables.

2. Start the application using Docker Compose:
    ```bash 
    docker-compose up
    ```

#### Running Locally

- Backend

1. Navigate to the backend directory:

    ```bash
    cd backend
    ```
2. Create a .env file in the root directory of the backend based on .env.example and fill in the environment variables.

3. Install dependencies and run the backend:

    ```bash
    go mod tidy
    go run main.go
    ```

- Frontend

1. Navigate to the frontend directory:

    ```bash
    cd frontend
    ```
2. Create a .env file in the root directory of the frontend based on .env.example and fill in the environment variables.

3. Install dependencies and run the frontend:

    ```bash
    npm install
    npm start
    ```