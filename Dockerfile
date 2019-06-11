FROM golang:latest AS build
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build -o tornimo-agent

FROM alpine:latest
#RUN adduser -Dh /tornimo tornimo 
#WORKDIR /tornimo
COPY --from=build /build/tornimo-agent /tornimo/tornimo-agent
ENTRYPOINT ["/tornimo/tornimo-agent"]
