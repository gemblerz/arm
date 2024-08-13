import time

import board
import busio
import adafruit_aw9523
from .driver import RobotDriver

class GoogleCoralAM9523Driver(RobotDriver):
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
        i2c = board.I2C()
        # Or, to create I2C bus on specific pins
        # i2c = busio.I2C(board.I2C2_SCL, board.I2C2_SDA)

        self.aw = adafruit_aw9523.AW9523(i2c)
        self.motor_1_dir = self.aw.get_pin(1)
        self.motor_1_dir.switch_to_output(value=True)
        self.motor_1_step = self.aw.get_pin(2)
        self.motor_1_step.switch_to_output(value=True)


    def _drive(self, pin_dir, pin_step, steps, clockwise=True):
        pin_dir.value = clockwise
        for i in range(steps):
            pin_step.value = True
            time.sleep(self.delay)
            pin_step.value = False
            time.sleep(self.delay)

    def release(self):
        pass

    def move(self, motor, speed, forward=True):
        if motor == 1:
            self._drive(self.motor_1_dir, self.motor_1_step, speed, forward)
        elif motor == 2:
            self._drive(GoogleCoralDriver.PIN_MOTOR2_DIR, GoogleCoralDriver.PIN_MOTOR2_STEP, speed, forward)
        elif motor == 3:
            self._drive(GoogleCoralDriver.PIN_MOTOR3_DIR, GoogleCoralDriver.PIN_MOTOR3_STEP, speed, forward)
        elif motor == 4:
            self._drive(GoogleCoralDriver.PIN_MOTOR4_DIR, GoogleCoralDriver.PIN_MOTOR4_STEP, speed, forward)