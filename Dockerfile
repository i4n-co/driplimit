# golang image multistage build

# Build stage
FROM golang:1.22 AS build-env
WORKDIR /go/src/driplimit

COPY . .

RUN go install honnef.co/go/tools/cmd/staticcheck@latest
RUN go mod download
RUN go vet ./...
RUN staticcheck ./...
RUN go test ./...
RUN go build -ldflags="-extldflags=-static" -o /driplimit github.com/i4n-co/driplimit/cmd/driplimit


# Final stage
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build-env /driplimit /usr/local/bin/driplimit
VOLUME /home/nonroot/
ENV DATA_DIR /home/nonroot/
ENV ADDR 0.0.0.0
EXPOSE 7131

ENTRYPOINT ["/usr/local/bin/driplimit"]