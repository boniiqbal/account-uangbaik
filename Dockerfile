FROM refactoryteam/docker-base:xenial-latest

LABEL Maintainer="Refactory Engineering <engineering@refactory.id>"
LABEL Name="account-microservice-v.1.0"
ENV GOLANG_VERSION 1.12.4

# Install tools
RUN apt-get update && apt-get install -y curl git mysql-client && \
    rm -rf /var/lib/apt/lists/*

RUN curl -sSL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz \
		| tar -v -C /usr/local -xz

# setup go app dir
ENV APP_DIR account-microservice

RUN mkdir -p go/app/$APP_DIR
WORKDIR /go/app/$APP_DIR
ADD . /go/app/$APP_DIR
ADD .env.prod /go/app/$APP_DIR/.env

# Copy account-microservice config
COPY ./docker-config/default /etc/nginx/sites-available/default
COPY ./docker-config/uangbaik.conf /etc/supervisor/conf.d/uangbaik.conf

ENV GO111MODULE on
ENV PATH /usr/local/go:$PATH
ENV GOROOT /usr/local/go
ENV GOPATH $HOME/go
ENV PATH PATH=$GOPATH/bin:$GOROOT/bin:/go/bin:$PATH

RUN go get . && go build .
RUN go get github.com/gobuffalo/pop && go install github.com/gobuffalo/pop/soda

RUN apt-get purge -y curl git && apt-get autoremove -y

CMD ["/usr/bin/supervisord"]
