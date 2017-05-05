FROM golang

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/superliuwr/poc-go-kinesis"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/poc-go-kinesis/bin

WORKDIR /opt/poc-go-kinesis/bin

COPY bin/poc-go-kinesis /opt/poc-go-kinesis/bin/
RUN chmod +x /opt/poc-go-kinesis/bin/poc-go-kinesis

CMD /opt/poc-go-kinesis/bin/poc-go-kinesis