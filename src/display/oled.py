import time

from luma.core.interface.serial import i2c
from luma.core.render import canvas
from luma.oled.device import ssd1306

from PIL import Image

serial = i2c(port=1, address=0x3C)
device = ssd1306(serial)

with canvas(device) as draw:
    draw.rectangle(device.bounding_box, outline="white", fill="black")
    draw.text((30, 40), "Hello World", fill="white")

time.sleep(3)

photo = Image.open("data/pixelart1.png")
resized_photo = photo.resize(device.size)
device.display(resized_photo.convert(device.mode))

while True:
    time.sleep(1)
