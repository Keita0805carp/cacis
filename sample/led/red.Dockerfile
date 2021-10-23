FROM ubuntu

COPY led-red.py .

RUN apt update && apt install -y python3-rpi.gpio

CMD python3 led-red.py
