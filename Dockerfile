## Build
FROM golang:1.21-alpine AS build
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
WORKDIR /
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build  -o /app

## Deploy
FROM alpine
WORKDIR /
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]