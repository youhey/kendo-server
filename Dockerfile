FROM golang:1.26-bookworm AS build

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /out/kendo-server ./cmd/kendo-server

FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates wget \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /out/kendo-server /usr/local/bin/kendo-server

EXPOSE 8080

CMD ["kendo-server"]
