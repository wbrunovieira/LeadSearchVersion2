FROM golang:1.23-alpine


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLE=0 go build -o main .




EXPOSE 8085

CMD ["./main"]
