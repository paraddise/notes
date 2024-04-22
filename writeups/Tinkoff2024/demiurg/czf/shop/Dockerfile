FROM golang:1.20

COPY ./src /app/src

WORKDIR /app/src

RUN go build -o ./main -mod=vendor

CMD /app/src/main
