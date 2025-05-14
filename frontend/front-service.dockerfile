FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o brokerFront .


FROM scratch
WORKDIR /app

COPY --from=builder /app/brokerFront ./brokerFront
COPY --from=builder /app/cmd/web/templates ./cmd/web/templates

CMD [ "./brokerFront" ]
