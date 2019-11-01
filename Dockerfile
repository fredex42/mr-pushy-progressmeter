FROM alpine:3.10

ADD servercontent.tar.gz /opt/mrpushy-progressmeter
WORKDIR /opt/mrpushy-progressmeter
ENTRYPOINT /opt/mrpushy-progressmeter/server
