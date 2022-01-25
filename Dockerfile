FROM golang:1.17 AS build-env

WORKDIR /build
COPY . .

RUN make build

FROM gcr.io/distroless/base
COPY --from=build-env /build/bin ./
COPY --from=build-env /build/config ./config
COPY --from=build-env /build/docs ./docs
COPY --from=build-env /build/db/migrations ./db/migrations

ENTRYPOINT ["./rpc-runtime"]
CMD ["server"]