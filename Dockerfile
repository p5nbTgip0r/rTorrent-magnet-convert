FROM golang:alpine as builder
RUN apk add --no-cache ca-certificates git && \
    mkdir /build
ADD . /build
RUN cd /build && \
    go mod download && \
    go build -o executable . && \
    chmod 755 executable

FROM scratch
COPY --from=builder /build/executable /rtorrent-magnet-convert

ENTRYPOINT [ "/rtorrent-magnet-convert" ]