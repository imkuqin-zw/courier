#alpine.Dockerfile

WORKDIR /app

ADD dist/leaf .
ADD dist/conf ./conf/

ENTRYPOINT ["/app/leaf"]