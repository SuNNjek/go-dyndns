FROM golang:1.19 AS build

WORKDIR /build
COPY . .

# Use CGO_ENABLED so that the binary gets built statically, otherwise it won't run because of "FROM scratch"
RUN CGO_ENABLED=0 go build -ldflags '-s -w' -o go-dyndns


FROM scratch

# Copy root certificates from build container so that HTTPS requests work
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/go-dyndns /bin/go-dyndns

ENTRYPOINT [ "/bin/go-dyndns" ]