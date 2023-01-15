package aw

import "github.com/sclevine/agouti"

func GetWebDriver() (*agouti.WebDriver, error) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--no-sandbox",
			"--window-size=1280,720",
		}),
		agouti.Debug,
	)
	err := driver.Start()
	return driver, err
}
