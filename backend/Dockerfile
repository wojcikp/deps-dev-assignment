FROM golang:1.22.1-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache gcc musl-dev libgcc libc-dev

FROM base AS build-stage
COPY . .
RUN CGO_ENABLED=1 go build -o deps-dev-assignment-backend ./cmd/app

FROM alpine:3.18
WORKDIR /app
COPY --from=build-stage /app/deps-dev-assignment-backend .
COPY --from=build-stage /app/data /app/data
CMD ["./deps-dev-assignment-backend"]
