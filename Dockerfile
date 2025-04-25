FROM golang AS compiler_stage
RUN mkdir -p /go/src/LogisticCenter
WORKDIR /go/src/LogisticCenter
ADD main.go .
ADD go.mod .
RUN go build -o logistic-center

FROM alpine:latest
LABEL version="1.0"
LABEL maintainer="DENIS KOZLOV"
COPY --from=compiler_stage /go/src/LogisticCenter/logistic-center .
ENTRYPOINT [ "./logistic-center" ]