FROM golang:1.16 as build
WORKDIR /app
COPY go.mod go.sum ./
COPY Makefile /app
COPY . ./
RUN go mod download
RUN go run github.com/prisma/prisma-client-go prefetch
RUN go generate ./...
RUN make build

FROM alpine:latest
COPY --from=build /app/build/tweet-extractor .
CMD [ "/tweet-extractor" ]