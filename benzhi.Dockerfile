FROM golang:1.23

ENV GOPROXY=off GOSUMDB=off

WORKDIR /app
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
COPY web ./web

RUN go build -mod=vendor ./... && go build -mod=vendor -o /app/turbine ./cmd/turbine

CMD ["/app/turbine"]
