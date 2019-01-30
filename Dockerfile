# https://www.cloudreach.com/blog/containerize-this-golang-dockerfiles/
#
# First create main in a rich build environment

FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o main .

# then copy it to a small image.

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
