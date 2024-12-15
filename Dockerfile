FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /tmp/tfbin

FROM alpine:3.19
COPY --from=builder /tmp/tfbin /usr/local/bin

CMD ["sh", "-c", "tfbin --user=${GITHUB_ACTOR} --verb=${INPUT_VERB} --tasks=\"${INPUT_TASKS}\" --workers=${INPUT_WORKERS} --version=${INPUT_VERSION}"]