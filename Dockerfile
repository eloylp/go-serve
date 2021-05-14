FROM golang:1.16.3 AS build
WORKDIR /src
COPY . .
RUN useradd -u 10001 nonprivuser

RUN make build

FROM scratch
# Copy CA certificates from build image.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy users table from build image
COPY --from=build /etc/passwd /etc/passwd
# Copy binary build image
COPY --from=build /src/dist/go-serve /app/go-serve

EXPOSE 8080
USER nonprivuser
WORKDIR /app
CMD ["./go-serve"]