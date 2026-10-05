FROM golang:1.26-bookworm AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /out && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/sampler .

FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install --yes --no-install-recommends ca-certificates libasound2 \
  && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /home/sampler && chown 10001:10001 /home/sampler

ENV HOME=/home/sampler
COPY --from=build /out/sampler /usr/local/bin/sampler
USER 10001:10001
ENTRYPOINT ["/usr/local/bin/sampler"]
