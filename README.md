
# Product Microservice

This microservice was developed during my 2nd year of master's at [ISEP](https://www.isep.ipp.pt/) and focuses exclusively on the creation and updating of product data. It was built with Go and follows a microservice architecture to ensure modularity and scalability.

## Table of Contents
- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Setup](#setup)

## Overview

This service provides RESTful endpoints for managing product information within a larger microservices ecosystem. It was designed to handle create and update operations with a focus on simplicity and performance, leveraging Go's strong concurrency model.

## Tech Stack

- **Language**: Go
- **Framework**: Fiber
- **Database**: MongoDB
- **Architecture**: RESTful API, Microservice

## Setup

### Prerequisites

- Go 1.19+
- MongoDB
- Docker

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/joao-silva1007/product-command.git
   cd product-command
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up environment variables (create a `.env` file):
   ```dotenv
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_user
   DB_PASS=your_password
   DB_NAME=your_database
   ```

4. Run the application:
   ```bash
   go run main.go
   ```