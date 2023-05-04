FROM golang:1.19

WORKDIR /chat

# Copy application source.
COPY . .

# Download dependencies.
RUN go mod download && go mod verify

# Build executable
RUN go build -o chat_coordinator cmd/chat/main.go

CMD ["./chat_coordinator"]
