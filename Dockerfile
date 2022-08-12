FROM gocv/opencv:4.6.0
LABEL maintainer chenpan1@cmschina.com.cn

# docker run --rm -it -p 2333:2333 xshrim/ocr bash

ARG LOAD_LANG=eng

ENV GOPATH /go
ENV PATH=${PATH}:${GOPATH}/bin

COPY . /root/ocr

WORKDIR /root/
RUN apt update && apt install -y lsb-release apt-transport-https wget && wget -O - https://notesalexp.org/debian/alexp_key.asc | apt-key add - && echo "deb https://notesalexp.org/tesseract-ocr4/$(lsb_release -cs)/ $(lsb_release -cs) main"| tee /etc/apt/sources.list.d/notesalexp.list > /dev/null && apt update && apt install -y ca-certificates libtesseract-dev=4.1.1-2.1 tesseract-ocr=4.1.1-2.1 tesseract-ocr-eng tesseract-ocr-jpn tesseract-ocr-chi-sim tesseract-ocr-chi-tra && cd /root/ocr && go build -o /root/gocr main.go

CMD ["/root/gocr"]