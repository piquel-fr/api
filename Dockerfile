FROM golang:1.24.4 AS builder

WORKDIR /api.piquel.fr

# Setup env
RUN export PATH="$PATH:$(go env GOPATH)/bin"

# Setup go dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy everything else
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -o ./bin/main ./main.go

# Now for run env
FROM alpine:latest

WORKDIR /api.piquel.fr

# Copy static files and configuration
COPY --from=builder /api.piquel.fr/bin/main .

CMD [ "./main" ]
