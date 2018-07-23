FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
COPY . $GOPATH/src/
WORKDIR $GOPATH/src/

# get dependencies
RUN go get -d -v
RUN go get k8s.io/client-go/...

# build the binary
RUN go build -o /go/bin/create_job.go

# STEP 2 build a small image
# start from alpine
FROM alpine:3.8

# copy our static executable
COPY --from=builder /go/bin/create_job.go /bin/create_job.go
ENTRYPOINT ["/bin/create_job.go"]
