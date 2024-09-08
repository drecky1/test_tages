FROM golang:1.22-bullseye as w
WORKDIR /app
COPY .. ./
RUN env CGO_ENABLED=0 go build -v -ldflags=" \
                -X 'build_data.User=$(id -u -n)'  \
                -X 'build_data.Time=$(date)'  \
                -X 'build_data.Version=$(cat version.txt)'" -o build/test_tages_server server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=w /app/build/test_tages_server ./
CMD ["./test_tages_server"]