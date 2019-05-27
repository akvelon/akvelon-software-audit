FROM golang:1.12.5-alpine

# Force the go compiler to use modules
ENV GO111MODULE=on

WORKDIR $GOPATH/src/github.com/akvelon/akvelon-software-audit

RUN apk update && apk upgrade && apk add --no-cache git make \
        && go get golang.org/x/tools/go/vcs 

# We want to populate the module cache based on the go.{mod,sum} files
COPY go.mod .
COPY go.sum .

RUN go mod download

# Import the code from the context
COPY . .

# Compile the project
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 777

CMD ["./akvelon-software-audit"]