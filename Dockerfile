FROM golang:latest

# Copy the local package files to the container’s workspace.

WORKDIR /go/src/github.com/KaustubhLonkar
RUN cd /go/src/github.com/KaustubhLonkar \
    && git https://github.com/KaustubhLonkar/order-management-go.git

RUN cd /go/src/github.com/KaustubhLonkar/order-management-go
# Install our dependencies
RUN go get github.com/go-sql-driver/mysql  
RUN go get github.com/gin-gonic/gin
RUN go get github.com/segmentio/kafka-go
RUN go get github.com/segmentio/kafka-go/snappy
RUN go get github.com/jinzhu/gorm/dialects/mysql
RUN go get github.com/jinzhu/gorm
RUN go get github.com/rs/zerolog/log
RUN go get github.com/gin-gonic/contrib/static
RUN go get github.com/zsais/go-gin-prometheus

# Install api binary globally within container 
RUN go install github.com/KaustubhLonkar/order-management-go

# Set binary as entrypoint
ENTRYPOINT /go/bin/order-management-go

# Expose default port (8888)
EXPOSE 8888 