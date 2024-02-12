ARG GOLANG_VERSION=1.22
FROM golang:${GOLANG_VERSION}-alpine as builder

WORKDIR /go/src/tern

ARG TERN_VERSION=2.1.1
RUN wget https://github.com/jackc/tern/releases/download/v${TERN_VERSION}/tern_${TERN_VERSION}_linux_amd64.tar.gz && \
    tar xf tern_${TERN_VERSION}_linux_amd64.tar.gz

WORKDIR /go/src/app

RUN apk add --no-cache tzdata
ENV TZ=UTC

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.* .
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-s -w" -o dist/app ./main.go

FROM scratch

WORKDIR /opt/bin

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /go/src/tern/tern /opt/bin/tern
COPY --from=builder /go/src/app/migrations /opt/bin/migrations

COPY --from=builder /go/src/app/dist/app /opt/bin/app

ENV TZ=UTC
ENV USER=1000
ENV GIN_MODE=release
EXPOSE 8080

CMD [ "/opt/bin/app" ]
