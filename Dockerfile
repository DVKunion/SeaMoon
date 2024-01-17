# build stage
FROM golang:alpine AS build
ARG VERSION
COPY .. /src
WORKDIR /src
ENV CGO_ENABLED 0
ENV VERSION=${VERSION}
RUN go build  -ldflags "-X github.com/DVKunion/SeaMoon/server/consts.Version=${VERSION}" -o /tmp/seamoon cmd/main.go
RUN chmod +x /tmp/seamoon
# run stage
FROM alpine:3.19
LABEL maintainer="dvkunion@gamil.com"
WORKDIR /app
COPY --from=build /tmp/seamoon /app/seamoon
EXPOSE 1080 8080 7777 9000
ENTRYPOINT ["/app/seamoon"]