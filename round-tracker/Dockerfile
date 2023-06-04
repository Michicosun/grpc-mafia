FROM golang:1.19

WORKDIR /mafia_tracker

# Copy application source.
COPY . .

# Download golang dependencies.
RUN go mod download && go mod verify

# Build executable
RUN go build -o mafia_tracker cmd/tracker/main.go

CMD ["./mafia_tracker"]
