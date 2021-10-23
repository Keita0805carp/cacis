FROM ubuntu

COPY green-led.py .

RUN apt update && apt install -y python3-rpi.gpio

CMD python3 green-led.py
