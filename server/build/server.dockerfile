FROM golang:1.22-alpine

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod tidy && go mod download -x

COPY . .
RUN go build -C cmd/server

ENTRYPOINT ["cmd/server/server"]