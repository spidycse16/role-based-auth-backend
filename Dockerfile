FROM golang:1.23.0-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .
RUN go mod tidy

COPY .env .env
# Expose port 8080
EXPOSE 8080

#uncommnet this line and comment CMD["go", "run" , "cmd/migration/main.go"] for running migration for first time
# CMD ["go", "run", "cmd/migration/main.go"]

#after the migration is done always run this line

CMD ["go","run","cmd/main.go"]