FROM golang:1.26.4-alpine3.24 AS build

WORKDIR app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app cmd/manager/main.go

FROM alpine:3.24.0

LABEL maintainer=bormon_off
LABEL contact=tg@bormon_off

COPY --from=build /build/app app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["/app"]
