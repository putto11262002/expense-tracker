# Expense Tracker

A shared expense tracking application implemennted in Go. The application allows users to create groups, add other users to the group and track their shared expenses. Users can add expenses to the group, specifying splits of the expense and settle any outstanding debts among the group members. 

## Project structure

### api

The api directory contains the code for the api. The api is written using the Gin web framework in GO.

### web

The web directory contains the code for the web application. It is written in React + Typescript + Vite. 


## Getting Started
### Prerequisites
Before you begin, ensure you have the following software installed:

- Go (v1.19 or higher)
- Node.js (v18 or higher):


### API 

**Configure Environment Variables**

In the api directory, create a .env file with the following content and adjust as needed

```
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_NAME=expense_tracker
ALLOWED_ORIGINS=http://localhost:3000 
DOMAIN=localhost
```

Note: In production environment (GO_ENV=production) environment variables are loaded straight from the shell not the .env file

**Compile and run directly (Development)**

To run compile and run api directly from source code run the following command in the root of the project.

```
make dev-api
```

**Build And Run Binary (Production)**

To compile the source code into binary and run the binary run the following command in the root of the project

```
make build-api 
./api/bin/main 
```

## Web App

**Install Dependencies**

In the web directory, install the web application's dependencies by running:

```
npm install
```

**Configure Environment Variables**

In the web directory, create a .env file with the following content and adjust as needed

```
API_BASE_URL=http://localhost:3000
```

**Running With Development Server**

Start development server on port 3000

```
make dev-web
```

**Building**

Build application into static bundle

```
make build-web
```

## Getting Started With Docker

This approach starts the api server, web application server and database server without any additional configurations. By default the api server will be running on port 3001 and the we application server will be running on port 3000.
### Prerequisites
Before you begin, ensure you have the following software installed:

- Docker 

- Docker Compose (if not already installed with docker)

### Build Docker Images

Building docker images for api and web application.
```
make docker-build
```

### Run Docker Containers

```
make docker-run
```


### Stop Docker Containers

```
make docker-stop
```





