# Kyooar

## Prerequisites

Before running this application, ensure you have the following installed:

### Required Software

1. **Docker** - For running PostgreSQL database
   - Install from: <https://docker.com/get-started>
   - Verify installation: `docker --version`

2. **Go 1.24.5** - Backend API server
   - Install from: <https://golang.org/dl/>
   - Verify installation: `go version`

3. **Node.js e24.4.1** - Frontend development
   - Install from: <https://nodejs.org/>
   - Verify installation: `node --version`

4. **npm e11.4.2** - Package manager
   - Comes with Node.js or install separately
   - Verify installation: `npm --version`

## Quick Start

### Setup Backend

- Make sure docker is running.
- Run `cd backend`.
- Run `make setup`.

### Setup Frontend

- Run `cd frontend`.
- Run `npm i`.

### Run Application

- On frontend folder run `npm run dev`.
- On backend folder run `make dev`.
- See terminal folder to see where application will be accessible.

## Reset DB

Make sure you are on backend folder, run the command `docker compose down -v`, after that, reconfigure backend using `make setup` in the backend folder.

## Build Docker Images

At the root folder run the command `make build`, this will generate the docker images. Remember the application needs
some environment variables to be set, see this in the .env.example inside the backend folder.
