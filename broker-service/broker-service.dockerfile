FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o brokerApp ./cmd/api


FROM scratch
WORKDIR /app

COPY --from=builder /app/brokerApp ./brokerApp
CMD [ "./brokerApp" ]
