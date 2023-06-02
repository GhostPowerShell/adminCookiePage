FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod init github.com/GhostPowerShell/adminCookiePage
RUN go mod tidy

RUN go build -o main .

EXPOSE 1337

CMD ["./main"]
