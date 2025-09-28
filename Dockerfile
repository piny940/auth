ARG GO_VERSION=1.25.1
FROM golang:${GO_VERSION} AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main /app/cmd/main.go

FROM scratch AS final

WORKDIR /app

COPY --from=build /app/main .

CMD [ "/app/main" ]
