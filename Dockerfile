FROM golang:1.7.1

# install nginx
RUN apt-get update -y -qq
RUN apt-get install -y -qq nginx

COPY docker/site.conf /etc/nginx/sites-available/default

# install summar
COPY  . /go/src/github.com/hankjacobs/summar
WORKDIR /go/src/github.com/hankjacobs/summar

RUN go get -d -v
RUN go install -v

# install requestgen
WORKDIR /go/src/github.com/hankjacobs/summar/docker/requestgen
RUN go install -v

# make sure /var/log/stats.log is there so tail doesn't complain
RUN touch /var/log/stats.log

CMD nginx & requestgen & summar & tail -f /var/log/stats.log
