FROM golang:1.22 AS compiling_stage

WORKDIR /apphttp

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o apphttp ./main.go

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache libc6-compat

COPY --from=compiling_stage /apphttp .

RUN chmod +x /root/apphttp

ENTRYPOINT ["./apphttp"]

EXPOSE 8080
