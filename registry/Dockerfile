FROM golang:1.19

WORKDIR /mafia_registry

# Copy application source.
COPY . .

# Download ubuntu dependencies.
RUN apt update && apt -y install wkhtmltopdf

# Download golang dependencies.
RUN go mod download && go mod verify

# Build executable
RUN go build -o mafia_registry cmd/registry/main.go

CMD ["./mafia_registry"]
