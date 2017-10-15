# Build image
FROM golang:latest

WORKDIR /go/src/github.com/vlad-stoian/postfacto-concourse-resource/

ADD . .

RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/check cmd/check/main.go

# Run image
FROM concourse/buildroot:git

RUN ln -s /usr/bin/gpg2 /usr/bin/gpg

COPY --from=0 /go/bin/check /opt/resource/

# ADD assets/ /opt/resource/
RUN chmod +x /opt/resource/*
