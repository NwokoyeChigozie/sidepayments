# Build stage
FROM golang:1.20.1-alpine3.17 as build

# Install wkhtmltopdf dependencies
RUN apk add --no-cache wkhtmltopdf=0.12.6-r0

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN if test -e app.env; then echo 'found app.env'; else mv app-sample.env app.env; fi; \
    go build -v -o /dist/vesicash-payment-ms

# Deployment stage
FROM alpine:3.17

WORKDIR /usr/src/app

COPY --from=build /usr/src/app ./

COPY --from=build /dist/vesicash-payment-ms /usr/local/bin/vesicash-payment-ms

CMD ["vesicash-payment-ms"]
