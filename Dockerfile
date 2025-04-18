FROM golang:1.24.2 AS builder

WORKDIR /api.piquel.fr

# Setup env
RUN export PATH="$PATH:$(go env GOPATH)/bin"

# Dependencies
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Setup go dependencies
COPY go.mod .
RUN go mod download

# Generate sqlc files
COPY sqlc.yml .
COPY database database
RUN sqlc generate

# Copy everything else
COPY . .

RUN go mod tidy

# Build the binary
RUN CGO_ENABLED=0 go build -o ./bin/main ./main.go

# Now for run env
FROM alpine:latest

WORKDIR /api.piquel.fr

# Copy static files and configuration
COPY --from=builder /api.piquel.fr/bin/main .

CMD [ "./main" ]
