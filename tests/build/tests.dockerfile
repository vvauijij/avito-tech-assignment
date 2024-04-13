FROM golang:1.22-alpine

WORKDIR /tests

COPY go.mod ./
RUN go mod tidy && go mod download -x

COPY . .

ENTRYPOINT ["go", "test"]