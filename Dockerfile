# builder image
FROM golang:1.15-alpine as builder

WORKDIR /build
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app ./cmd/shazam


FROM alpine:3.12.0
COPY --from=builder /build/app ./app
COPY --from=builder /build/conf/docker.yaml ./conf.yaml

CMD [ "./app", "start", "./conf.yaml" ]
