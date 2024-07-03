FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/rarimo/forms-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/forms-svc /go/src/github.com/rarimo/forms-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/forms-svc /usr/local/bin/forms-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["forms-svc"]
