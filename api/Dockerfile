FROM golang:1.21-alpine AS builder

WORKDIR /src
COPY . .
RUN apk add --no-cache git
RUN go mod tidy && go build -o /bin/api main.go

FROM alpine:3.18
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /bin/api /bin/api
USER appuser:appgroup
EXPOSE 8080
ENTRYPOINT ["/bin/api"]
