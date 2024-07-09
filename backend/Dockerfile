FROM golang:1.22.3-alpine
LABEL authors="wavycat"

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

RUN CGO_ENABLED=1 go mod download

RUN CGO_ENABLED=1 go build -o /daec-backend github.com/wavy-cat/DAEC/backend/cmd/backend

ENTRYPOINT ["/daec-backend"]