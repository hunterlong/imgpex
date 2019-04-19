FROM golang:1.12-alpine as base
RUN apk add --no-cache libstdc++ gcc g++ make git ca-certificates linux-headers
WORKDIR /go/src/github.com/hunterlong/imgpex
ADD . /go/src/github.com/hunterlong/imgpex/
RUN go get
RUN go install

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
RUN mkdir /app/images
VOLUME /app
COPY --from=base /go/bin/imgpex /usr/local/bin/imgpex
COPY --from=base /go/src/github.com/hunterlong/imgpex/images.txt /app/images.txt
CMD imgpex
