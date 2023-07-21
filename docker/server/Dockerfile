# build stage
FROM golang:alpine AS build
ARG VERSION
COPY ../.. /src
WORKDIR /src
ENV CGO_ENABLED 0
ENV VERSION=${VERSION}
RUN go build  -ldflags "-X github.com/DVKunion/SeaMoon/pkg/consts.Version=${VERSION}" -o /tmp/server cmd/aliyun_server.go
RUN chmod +x /tmp/server
# run stage
FROM scratch
LABEL maintainer="dvkunion@gamil.com"
WORKDIR /app
COPY --from=build /tmp/server /app/server
EXPOSE 9000
ENTRYPOINT ["/app/server"]