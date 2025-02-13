FROM golang:1.23.5 AS builder

WORKDIR /api.piquel.fr

# Setup env
RUN export PATH="$PATH:$(go env GOPATH)/bin"

# Dependencies
# Go dependencies
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/a-h/templ/cmd/templ@latest
# Tailwind
RUN apt-get install curl
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
RUN mv tailwindcss-linux-x64 tailwindcss
RUN chmod +x tailwindcss
RUN mv tailwindcss /usr/bin

# Setup go dependencies
COPY go.mod .
RUN go mod download

# Copy static files
COPY public public

# Generate sqlc files
COPY sqlc.yml .
COPY database database
RUN sqlc generate

# Templ files
COPY views views
COPY components components
RUN tailwindcss -i views/css/styles.css -o public/styles.css

# Copy everything else
COPY . .

# Generate templ related files
RUN templ generate

RUN go mod tidy

# Build the binary
RUN CGO_ENABLED=0 go build -o ./bin/main ./main.go

# Now for run env
FROM alpine:latest

WORKDIR /api.piquel.fr

# Copy static files and configuration
COPY --from=builder /api.piquel.fr/bin/main .
COPY --from=builder /api.piquel.fr/public public
COPY --from=builder /api.piquel.fr/config config

CMD [ "./main" ]
