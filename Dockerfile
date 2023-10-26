FROM golang:1.20-alpine
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /main

EXPOSE 5000

CMD [ "/main" ]