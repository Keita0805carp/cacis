FROM ubuntu

COPY red-led.py .

RUN apt update && apt install -y python3-rpi.gpio

CMD python3 red-led.py
