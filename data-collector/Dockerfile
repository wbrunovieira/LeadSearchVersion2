FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod  ./
RUN go mod download

COPY . .

RUN CGO_ENABLE=0 go build -o main .

CMD ["./main"]