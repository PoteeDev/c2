FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /c2 .

FROM alpine
WORKDIR /usr/app
COPY --from=build /c2 .
ENV PORT=8080
ENTRYPOINT [ "./c2"]
