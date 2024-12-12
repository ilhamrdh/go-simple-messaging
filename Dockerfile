FROM golang:1.22.2-alpine

WORKDIR /app

COPY go.mod go.sum  ./

RUN go mod tidy

COPY . .

RUN go build -o go-simple-messaging-app && chmod +x go-simple-messaging-app

EXPOSE 4000

EXPOSE 8080

CMD [ "./go-simple-messaging-app" ]