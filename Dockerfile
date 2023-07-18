# It is built in one Dockerfile deliberately, becuase mutlti staging FROM directives can be used:
# https://docs.docker.com/build/building/multi-stage/
# Thanks to this solution it is isolated from local environment, and eventually it will produce small image.
ARG GO_VERSION=1.20.5
# First stage downloads all necesarry packages. It should be fast in 99% cases unless go.mod is changed.
FROM golang:${GO_VERSION} AS downloader

WORKDIR /app

COPY go.mod .
COPY go.sum .
# It is done only once unless go.mod has been changed.
RUN go mod download

# Second stage of a mutli-stage builds binary using packages from first stage.
# It will be fast unless source files are changed.
FROM downloader AS builder

COPY api api
COPY cmd cmd
COPY pkg pkg

RUN CGO_ENABLED=0 GO111MODULE=on go build -ldflags "-extldflags '-static'" -o /bin/ports ./cmd/main.go

# The third and last stage it is alpine 3.18 with a binary and assets files ~18 MB.
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /bin/ports /bin/ports
COPY assets /app/assets

EXPOSE 8080
ENTRYPOINT ["/bin/ports"]
# CMD directives could take path to the initial file which could be changed when container is started.

