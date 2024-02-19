FROM golang as build-app

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o app-binary *.go

FROM scratch

COPY --from=build-app /app/app-binary ./app-binary

EXPOSE 8080

CMD ["./app-binary"]
