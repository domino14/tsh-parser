FROM golang:alpine as build-env

RUN mkdir /opt/program
WORKDIR /opt/program

RUN apk update
RUN apk add build-base ca-certificates git

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

WORKDIR /opt/program/cmd/tshparser

ARG BUILD_HASH=unknown
ARG BUILD_DATE=unknown

RUN go build -ldflags  "-X=main.BuildDate=${BUILD_DATE} -X=main.BuildHash=${BUILD_HASH}"

# Build minimal image:
FROM alpine
COPY --from=build-env /opt/program/cmd/tshparser/tshparser /opt/tshparser
COPY --from=build-env /opt/program/migrations /opt/migrations
COPY --from=build-env /opt/program/cfg /opt/cfg
RUN apk --no-cache add curl
EXPOSE 8001

WORKDIR /opt
CMD ["./tshparser"]

LABEL org.opencontainers.image.source https://github.com/domino14/tsh-parser