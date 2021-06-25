# build
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git gcc
COPY . memcachedbetween
RUN cd memcachedbetween && go build -o memcachedbetween

## app
FROM alpine
WORKDIR /app
COPY --from=build-env /go/memcachedbetween/memcachedbetween /app/
CMD ./memcachedbetween