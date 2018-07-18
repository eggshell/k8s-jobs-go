FROM golang:1.10.3

RUN go get k8s.io/client-go/...
COPY create_job.go /
CMD ["go", "run", "/create_job.go"]
