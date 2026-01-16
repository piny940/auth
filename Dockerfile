ARG GO_VERSION=1.25.6
FROM golang:${GO_VERSION} AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main /app/cmd/main.go

FROM scratch AS final

WORKDIR /app

COPY --from=build /app/main .

CMD [ "/app/main" ]
