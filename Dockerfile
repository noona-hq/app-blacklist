FROM golang:1.20 AS build

WORKDIR /build
ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make install


FROM alpine:latest
COPY --from=build /go/bin/app-blacklist ./
COPY --from=build /build/server/templates/html ./server/templates/html

ENTRYPOINT ["./blacklist"]
