FROM golang:1.10.3
WORKDIR /go/src/github.com/satococoa/prbot
COPY lib lib
COPY main.go main.go
RUN go get -u gopkg.in/src-d/go-git.v4 && \
  go get -u github.com/google/go-github/github && \
  go get -u golang.org/x/oauth2
RUN go install github.com/satococoa/prbot
ENTRYPOINT [ "prbot" ]
