FROM golang:1.20 as build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o registry -a -gcflags=all="-l -B -wb=false" -ldflags="-w -s" *.go

FROM scratch
COPY --from=build /build/registry ./registry
COPY --from=build /build/openapi.yml ./registry.yml
ENV AWS_ACCESS_KEY_ID=mortymorty
ENV AWS_SECRET_ACCESS_KEY=mortymorty
ENTRYPOINT ["./registry"]