FROM golang:1.23.0

WORKDIR /Lab3

COPY . .

RUN go mod tidy
RUN go build -o supervisor2 .

CMD ["./supervisor2"]
