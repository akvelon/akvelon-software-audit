FROM golang:1.12.5-alpine

# Force the go compiler to use modules
ENV GO111MODULE=on

WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache git \
  && go get github.com/boyter/lc && apk add curl

# We want to populate the module cache based on the go.{mod,sum} files
COPY go.mod .
COPY go.sum .

RUN go mod download

# Import the code from the context
COPY . .

# Compile the project
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o audit-worker ./cmd/worker/main.go
RUN chmod +x scripts/run.sh

ENTRYPOINT [ "./scripts/run.sh" ]
