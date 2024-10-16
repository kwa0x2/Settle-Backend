FROM golang:1.23.1

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
COPY .env ../

RUN go build -o main ./cmd/

CMD [ "./main" ]
