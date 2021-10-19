FROM ubuntu

ENV TZ Asia/Tokyo
ENV DEBIAN_FRONTEND="noninteractive" 

RUN apt update \
 && apt install -y git libjpeg8-dev cmake \
 && git clone https://github.com/jacksonliam/mjpg-streamer.git \
 && cd mjpg-streamer/mjpg-streamer-experimental \
 && make \
 && make install

CMD ./start.sh
