FROM golang:latest

RUN mkdir /app
WORKDIR /app
ADD fyi /app/

EXPOSE 8888
CMD ["/app/fyi"]
