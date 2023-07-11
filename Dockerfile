# build stage
FROM golang:alpine AS build
ARG VERSION
COPY . /src
WORKDIR /src
ENV CGO_ENABLED 0
ENV VERSION=${VERSION}
RUN go build --ldflags "-s -w -X github.com/DVKunion/SeaMoon/pkg/consts.Version=${VERSION}" -o /tmp/client cmd/client.go
# run stage
FROM scratch
LABEL maintainer="dvkunion@gamil.com"
WORKDIR /app
COPY --from=build /tmp/client /app/client
RUN chmod +x client
EXPOSE 7777 1080 9000
ENTRYPOINT ["/app/client"]