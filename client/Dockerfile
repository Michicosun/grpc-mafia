FROM golang:1.19

WORKDIR /mafia_bot

# Copy application source.
COPY . .

# Download dependencies.
RUN go mod download && go mod verify

# Build executable
RUN go build -o mafia_bot cmd/client/main.go

CMD ["./mafia_bot"]
