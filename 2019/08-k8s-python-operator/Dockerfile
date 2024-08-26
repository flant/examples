FROM python:3.7.6-alpine3.9@sha256:d3953acf09227baf26919dbb4dc3637e8015ce46debbed4cf5444702728cfc16

ADD . /app

RUN pip3 install /app

ENTRYPOINT ["copyrator"]
