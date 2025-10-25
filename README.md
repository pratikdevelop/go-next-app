# URL Shortener Application

This is a full-stack URL shortener application with a Go backend and a Next.js frontend.

## Features

*   Shorten long URLs into concise, manageable links.
*   Option to create custom short codes for personalized URLs.
*   View a list of all shortened URLs.
*   Edit the original long URL for existing shortened links.
*   Delete shortened URLs.
*   Track click counts for each shortened URL.

## Technologies Used

### Backend (Go)

*   **GoLang**: Programming language.
*   **Gin Gonic**: Web framework for building APIs.
*   **MongoDB**: NoSQL database for storing URL data.
*   **Go DotEnv**: For loading environment variables.*   **CORS**: Middleware for handling Cross-Origin Resource Sharing.

### Frontend (Next.js)

*   **Next.js**: React framework for building user interfaces.
*   **React**: JavaScript library for building UIs.
*   **Tailwind CSS**: Utility-first CSS framework for styling.
*   **Material UI**: React component library for faster and easier web development.

## Setup and Installation

To set up the project, follow these steps:

### 1. Clone the Repository

```bash
git clone https://github.com/pratikdevelop/go-next-app
cd go-next-app
```

### 2. Backend Setup

Navigate to the `api` directory:

```bash
cd api
```

Install Go dependencies:

```bash
go mod tidy
```

Create a `.env` file in the `api` directory and add your MongoDB URI:

```
MONGO_URI="your_mongodb_connection_string"
```

Replace `your_mongodb_connection_string` with your actual MongoDB connection string (e.g., from MongoDB Atlas).

### 3. Frontend Setup

Navigate to the `my-app` directory:

```bash
cd ../my-app
```

Install Node.js dependencies:

```bash
npm install
```

## Running the Application

### 1. Run the Backend

From the `api` directory, run the Go application with `air`:

```bash
air
```

The backend server will start on `http://localhost:8081`.

### 2. Run the Frontend

From the `my-app` directory, run the Next.js development server:

```bash
npm run dev
```

The frontend application will be accessible at `http://localhost:3000`.

## API Endpoints

*   `POST /shorten`: Shorten a long URL. Accepts `long_url` and optional `short_code`.
*   `GET /urls`: Get all shortened URLs.
*   `PUT /urls/:id`: Update the long URL for a given shortened link.
*   `DELETE /urls/:id`: Delete a shortened URL.
*   `GET /:shortCode`: Redirect to the original long URL.