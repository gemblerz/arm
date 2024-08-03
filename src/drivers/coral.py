import time
from periphery import GPIO

from .driver import RobotDriver

class GoogleCoralDriver(RobotDriver):
    PIN_MOTOR1_DIR = GPIO("/dev/gpiochip0", 6, "out")
    PIN_MOTOR1_STEP = GPIO("/dev/gpiochip0", 7, "out")

    PIN_MOTOR2_DIR = GPIO("/dev/gpiochip0", 8, "out")
    PIN_MOTOR2_STEP = GPIO("/dev/gpiochip2", 13, "out")

    PIN_MOTOR3_DIR = GPIO("/dev/gpiochip2", 9, "out")
    PIN_MOTOR3_STEP = GPIO("/dev/gpiochip4", 10, "out")

    PIN_MOTOR4_DIR = GPIO("/dev/gpiochip4", 12, "out")
    PIN_MOTOR4_STEP = GPIO("/dev/gpiochip4", 13, "out")


    def __init__(self):
        microsteps = {
            'Full':1,
            'Half':2,
            '1/4':4,
            '1/8':8,
            '1/16':16,
            '1/32':32
        }
        self.delay = .005 / microsteps['Full']

    def _drive(self, pin_dir, pin_step, steps, clockwise=True):
        pin_dir.write(clockwise)
        for i in range(steps):
            pin_step.write(True)
            time.sleep(self.delay)
            pin_step.write(False)
            time.sleep(self.delay)

    def release(self):
        GoogleCoralDriver.PIN_MOTOR1_DIR.close()
        GoogleCoralDriver.PIN_MOTOR1_STEP.close()
        GoogleCoralDriver.PIN_MOTOR2_DIR.close()
        GoogleCoralDriver.PIN_MOTOR2_STEP.close()
        GoogleCoralDriver.PIN_MOTOR3_DIR.close()
        GoogleCoralDriver.PIN_MOTOR3_STEP.close()
        GoogleCoralDriver.PIN_MOTOR4_DIR.close()
        GoogleCoralDriver.PIN_MOTOR4_STEP.close()

    def move(self, motor, speed, forward=True):
        if motor == 1:
            self._drive(GoogleCoralDriver.PIN_MOTOR1_DIR, GoogleCoralDriver.PIN_MOTOR1_STEP, speed, forward)
        elif motor == 2:
            self._drive(GoogleCoralDriver.PIN_MOTOR2_DIR, GoogleCoralDriver.PIN_MOTOR2_STEP, speed, forward)
        elif motor == 3:
            self._drive(GoogleCoralDriver.PIN_MOTOR3_DIR, GoogleCoralDriver.PIN_MOTOR3_STEP, speed, forward)
        elif motor == 4:
            self._drive(GoogleCoralDriver.PIN_MOTOR4_DIR, GoogleCoralDriver.PIN_MOTOR4_STEP, speed, forward)