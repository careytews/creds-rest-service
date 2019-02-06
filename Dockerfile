FROM fedora:27

RUN mkdir /creds-rest-svc
WORKDIR /creds-rest-svc

COPY creds-rest-service /creds-rest-svc/

CMD ./creds-rest-service

EXPOSE 8080

