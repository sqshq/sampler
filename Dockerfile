FROM golang:1.22-bookworm AS builder

WORKDIR /src

RUN apt-get update && apt-get install --yes libasound2-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /sampler .

FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install --yes --no-install-recommends \
    libasound2 \
    curl \
    bc \
    procps \
    ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /sampler /usr/local/bin/sampler
COPY example.yml /app/example.yml

ENTRYPOINT ["sampler"]
CMD ["--config", "example.yml"]
