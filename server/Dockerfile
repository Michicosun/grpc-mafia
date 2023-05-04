FROM golang:1.19

WORKDIR /mafia_server

# Copy application source.
COPY . .

# Download dependencies.
RUN go mod download && go mod verify

# Build executable
RUN go build -o mafia_server cmd/server/main.go

CMD ["./mafia_server"]
