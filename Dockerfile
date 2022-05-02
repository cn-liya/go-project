# Build
FROM golang:1.22 as builder
ARG srv
WORKDIR /build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.google.cn
COPY . .
RUN go mod download && go build -tags sonic -o app ./${srv}

# Deploy
FROM scratch
ARG srv
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /build/app .
COPY ${srv}/conf.yaml .
COPY ${srv}/docs docs
COPY ${srv}/log log
ENTRYPOINT ["./app"]
