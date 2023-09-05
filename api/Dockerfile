FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

# Set environmental variables
ENV PORT=3000
ENV GO_ENV=production

WORKDIR /

COPY --from=build-stage /api /api

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/api"]