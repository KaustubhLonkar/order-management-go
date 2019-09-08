FROM busybox:latest

ENV GIN_MODE release
ENV NOAUTH false
ENV AWS_REGION us-east-1

WORKDIR /
# assumes that the availability-service-go binary was compiled via
# the Dockerfile-compile container with a linked volume of $builddir/bin
COPY bin/availability-service-go /
COPY *.json /

CMD ["/availability-service-go"]
EXPOSE 8888
USER nobody
