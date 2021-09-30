FROM golang:latest
WORKDIR /usr/src/app/
COPY . /usr/src/app/
RUN go mod download
RUN go build -o forum .
EXPOSE 5555
ENV TZ Asia/Almaty
CMD ["./forum"]




# docker build -t forum .
# docker run -d -p 5555:5555 forum