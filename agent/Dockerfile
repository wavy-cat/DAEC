FROM dockerhub.timeweb.cloud/library/golang:1.22.3
LABEL authors="wavycat"

WORKDIR /app

COPY . .

RUN go mod download

ENTRYPOINT ["go", "run", "agent"]