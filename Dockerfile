FROM golang:1.16

RUN useradd -ms /bin/bash deduplicator
USER deduplicator
WORKDIR /go/src/


USER deduplicator
RUN mkdir -p ~/.ssh
RUN mkdir -p deduplicator
COPY . deduplicator/
WORKDIR /go/src/deduplicator
ENV GO111MODULE=on


USER root
ENV export PATH=$PATH:/go/bin
ENV PROJECT_PATH /go/src/deduplicator
RUN chmod +x ./scripts/globalCoverage.sh && ./scripts/globalCoverage.sh
RUN CGO_ENABLED=0 GOOS=linux go build -v -o deduplicator

EXPOSE 8080

CMD ["/go/src/deduplicator/deduplicator"]
