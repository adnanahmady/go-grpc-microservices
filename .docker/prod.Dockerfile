FROM golang:1.24-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_NAME
RUN echo "Building ${SERVICE_NAME}..."

RUN CGO_ENABLED=0 go build -o /app/server ./cmd/${SERVICE_NAME}

FROM scratch

WORKDIR /app

COPY --from=builder /app/server /app

CMD ["./server"]