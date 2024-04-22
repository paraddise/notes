FROM golang:1.20

COPY ./src /app/src

WORKDIR /app/src

RUN go build -o ./main -mod=vendor main.go

RUN mv /app/src/main /app/main

WORKDIR /app

RUN rm -rf /app/src

CMD /app/main
