FROM golang:1.26-alpine AS builder

WORKDIR /app

# needed if any modules use git
RUN apk add --no-cache git

# cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy the source
COPY . .

# run go to build the static binary for linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o mressay ./cmd/mressay


# runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# cp the binary only (skip the source code)
COPY --from=builder /app/mressay .

# create the downloads directory (you need to volume mount it)
WORKDIR /app
USER nonroot:nonroot

ENTRYPOINT ["/app/mressay"]

