FROM golang:1.13.4
RUN go get -u github.com/mingrammer/flog
VOLUME "/tmp"
CMD ["--format", "apache_common", "--type", "log", "-n",  "100000000", "-s","2", "--output",  "/tmp/access.log"]
ENTRYPOINT [ "flog" ]