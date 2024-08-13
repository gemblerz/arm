from drivers.coral import GoogleCoralDriver
from drivers.coral_gpio_expander import GoogleCoralAM9523Driver

driver = GoogleCoralAM9523Driver()
driver.move(1, 200, forward=True)
driver.release()