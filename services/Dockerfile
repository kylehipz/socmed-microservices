# -----------------------------------------------
# STAGE 1 — Build the Go binary
# -----------------------------------------------
FROM golang:1.25.3 AS builder

ARG SERVICE=follow
WORKDIR /src

# Copy go.work and modules
COPY common/go.mod common/go.sum ./common/
COPY ${SERVICE}/go.mod ${SERVICE}/go.sum ./${SERVICE}/

# Sync workspace modules
RUN go work init && go work use common ${SERVICE}
# Pre-download all dependencies
RUN go work sync && go mod download

# Copy actual source code
COPY common ./common
COPY ${SERVICE} ./${SERVICE}

# Build statically
WORKDIR /src/${SERVICE}
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/service cmd/server/main.go

# -----------------------------------------------
# STAGE 2 — Distroless image 
# -----------------------------------------------
FROM gcr.io/distroless/static

# Copy binary
COPY --from=builder /app/service /service

# Run as non-root (distroless already does this)
USER nonroot:nonroot

ENTRYPOINT ["/service"]
