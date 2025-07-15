FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-go-app

EXPOSE 8080

CMD [ "go", "run", "main.go" ]

