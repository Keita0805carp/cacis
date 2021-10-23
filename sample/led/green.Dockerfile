FROM ubuntu

COPY led-green.py .

RUN apt update && apt install -y python3-rpi.gpio

CMD python3 led-green.py
