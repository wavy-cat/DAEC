FROM golang:1.22.3-alpine
LABEL authors="wavycat"

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /daec-agent github.com/wavy-cat/DAEC/agent

ENTRYPOINT ["/daec-agent"]