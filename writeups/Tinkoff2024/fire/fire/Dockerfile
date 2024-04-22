FROM golang:1.21.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY cmd /app/cmd
COPY internal /app/internal

RUN CGO_ENABLED=0 GOOS=linux go build -a -o fire cmd/fire/main.go


FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteravtive
RUN apt update
RUN apt -y install ca-certificates

ARG UID=2000
ARG GID=2000
ARG USER=app
ARG DIR=/home/${USER}

RUN addgroup --gid ${GID} ${USER} && adduser --gecos "" --disabled-password --gid ${GID} --uid ${UID} ${USER}

WORKDIR ${DIR}

COPY --from=builder /app/fire .

COPY static ${DIR}/static
COPY templates ${DIR}/templates
RUN chmod 444 static/js/* templates/*

RUN chown ${USER}.${USER} fire
USER ${USER}

CMD ["./fire"]