FROM golang:1.22.3-alpine
LABEL authors="wavycat"

WORKDIR /app

COPY . .

RUN go mod download

ENTRYPOINT ["go", "run", "backend/cmd/backend"]