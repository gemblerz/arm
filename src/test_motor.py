from drivers.coral import GoogleCoralDriver

driver = GoogleCoralDriver()
driver.move(1, 200, forward=False)
driver.release()