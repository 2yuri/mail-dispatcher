FROM golang:buster as builder

# add debug mode
RUN go get github.com/go-delve/delve/cmd/dlv

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /go/src/app

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

ENV MANAGER_PORT=19000

EXPOSE 19000 8081
CMD ["air", "d"]
