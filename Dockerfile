FROM golang:1.19 AS build

WORKDIR /build
COPY . .

# Use CGO_ENABLED so that the binary gets built statically
RUN CGO_ENABLED=0 go build -ldflags '-s -w' -o go-dyndns


# Use distroless base image. This already contains SSL certificates and timezone data
FROM gcr.io/distroless/static

COPY --from=build /build/go-dyndns /bin/go-dyndns

ENTRYPOINT [ "/bin/go-dyndns" ]