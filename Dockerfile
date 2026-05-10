FROM golang:1.22-bookworm AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app .

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=build /out/app /app/app

ENV PORT=3000
EXPOSE 3000

USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
