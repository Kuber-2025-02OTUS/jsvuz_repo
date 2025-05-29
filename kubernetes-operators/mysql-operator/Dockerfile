FROM golang:1.21 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 go build -a -o manager main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER nonroot:nonroot
ENTRYPOINT ["/manager"]
