# https://steele.blue/tiny-github-actions/
FROM dyweb/go-dev:1.13.6 as builder

LABEL maintainer="contact@dongyue.io"

ARG PROJECT_ROOT=/go/src/github.com/dyweb/weekly

WORKDIR $PROJECT_ROOT

COPY . $PROJECT_ROOT
RUN cd scripts/weekly && go install .

FROM ubuntu:18.04
LABEL maintainer="contact@dongyue.io"
WORKDIR /usr/bin
COPY --from=builder /go/bin/weekly .
ENTRYPOINT ["weekly"]
CMD ["help"]