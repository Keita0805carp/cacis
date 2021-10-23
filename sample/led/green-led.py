import RPi.GPIO as GPIO
from time import sleep

PIN = 21
 
GPIO.setmode(GPIO.BCM)
GPIO.setup(PIN, GPIO.OUT)
 
while True:
    try:
        print(GPIO.HIGH)
        GPIO.output(PIN, GPIO.HIGH)
        sleep(1)
        print(GPIO.LOW)
        GPIO.output(PIN, GPIO.LOW)
        sleep(1)
    except KeyboardInterrupt:
        print("cleanup...")
        break
GPIO.cleanup()
