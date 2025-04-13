FROM golang:1.22 AS compiling_stage

WORKDIR /appgprc

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o appgprc ./main.go

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache libc6-compat

COPY --from=compiling_stage /appgprc .

RUN chmod +x /root/appgprc

ENTRYPOINT ["./appgprc"]

EXPOSE 3000
