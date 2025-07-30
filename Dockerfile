FROM golang:1.23.2-alpine AS apibuild
WORKDIR /app
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend .
RUN go build -o dist/main ./main.go
RUN go build -o dist/seed ./cmd/seed/main.go
RUN go build -o dist/seed-plans ./cmd/seed-plans/main.go

FROM golang:1.23.2-alpine AS api
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate
WORKDIR /app
COPY --from=apibuild /app/dist/main .
COPY --from=apibuild /app/dist/seed .
COPY --from=apibuild /app/dist/seed-plans .
COPY ./backend/migrations ./migrations
COPY ./docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["./docker-entrypoint.sh"]

FROM node:20.18.1-alpine AS webbuild
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build
RUN npm prune --production || true

FROM node:20.18.1-alpine AS web
WORKDIR /app
COPY --from=webbuild /app/build ./build
COPY --from=webbuild /app/node_modules ./node_modules
COPY --from=webbuild /app/package.json ./package.json
EXPOSE 3000
CMD ["node", "build"]
