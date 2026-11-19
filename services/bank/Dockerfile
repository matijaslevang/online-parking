FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/service .

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=build /out/service /usr/local/bin/service
EXPOSE 8080
CMD ["service"]
