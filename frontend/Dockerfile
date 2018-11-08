FROM golang:1.8
ENV ENV_NAME=development
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...


CMD ["app"]