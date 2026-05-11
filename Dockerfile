FROM golang:1.25.3-bookworm AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
RUN go install github.com/air-verse/air@latest
RUN go install github.com/aarondl/sqlboiler/v4@v4.19.7
RUN go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@v4.19.7
COPY . .
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot AS prod
WORKDIR /app
COPY --from=builder /out/server /app/server
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]