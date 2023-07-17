package api

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tebeka/selenium"
)

const (
	// Example: https://www.banki.ru/products/currency/map/moskva/.
	baseURL = "https://www.banki.ru/products/currency/map/%s/"

	seleniumHubURL = "http://selenium-hub:4444/wd/hub"
)

var (
	// Debug mode. Default: false.
	Debug bool
)

// Client is a client for banki.ru.
type Client struct {
	buildURL  func() string
	webDriver selenium.WebDriver
}

// New creates a new client.
func NewClient() (*Client, error) {
	c := &Client{}
	var err error

	c.webDriver, err = selenium.NewRemote(selenium.Capabilities{"browserName": "chrome"}, seleniumHubURL)
	if err != nil {
		return nil, fmt.Errorf("failed to start selenium: %s", err)
	}

	return c, nil
}

// Branches returns branches info.
// Currency is USD by default, city is Moscow by default.
func (c *Client) Branches(crnc Currency, ct Region) (*Branches, error) {
	branches := &Branches{}

	if len(crnc) > 0 {
		branches.Currency = crnc
	}

	if len(ct) > 0 {
		branches.City = ct
	}

	c.buildURL = func() string {
		return fmt.Sprintf(baseURL, branches.City)
	}

	if Debug {
		log.Debug().Msgf("Fetching the currency rate from %s", c.buildURL())
	}

	r := &Branches{Currency: crnc, City: ct}
	b, err := c.parseBranches()
	if err != nil {
		return r, err
	}

	r.Items = b

	if Debug {
		log.Debug().Msgf("Found %d branches", len(b))
	}

	return r, err
}

// parseBranches parses branches info.
func (c *Client) parseBranches() ([]Branch, error) {
	var (
		b   []Branch
		err error
	)

	defer c.webDriver.Quit()

	t := time.Now()
	if err = c.webDriver.Get(c.buildURL()); err != nil {
		return []Branch{}, fmt.Errorf("failed to load page %s: %v", c.buildURL(), err)
	}

	log.Debug().Msgf("webDriver.Get took %v", time.Since(t))

	moveMouse := `var element = document.querySelector(".mapListstyled__StyledMapList-sc-294xv0-0.fdpae");
    var rect = element.getBoundingClientRect();
    var centerX = rect.left + (rect.width / 2);
    var centerY = rect.top + (rect.height / 2);
    var evt = new MouseEvent('mousemove', {
        bubbles: true,
        cancelable: true,
        view: window,
        clientX: centerX,
        clientY: centerY
    });
    element.dispatchEvent(evt);`
	if _, err = c.webDriver.ExecuteScript(moveMouse, nil); err != nil {
		return []Branch{}, fmt.Errorf("error executing script for mousemove: %v", err)
	}

	scroll := "window.scrollTo(0, document.body.scrollHeight);"
	if _, err = c.webDriver.ExecuteScript(scroll, nil); err != nil {
		return []Branch{}, fmt.Errorf("error executing script for scroll: %v", err)
	}

	r, err := c.webDriver.FindElement(selenium.ByCSSSelector, ".fdpae")
	if err != nil {
		return []Branch{}, fmt.Errorf("failed find element by css selector .fdpae: %v", err)
	}

	data, err := r.FindElements(selenium.ByCSSSelector, ".cITBmP")
	if err != nil {
		return []Branch{}, fmt.Errorf("failed find element by css selector .cITBmP: %v", err)
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 3)
	t = time.Now()

	for _, d := range data {
		wg.Add(1)
		sem <- struct{}{}

		go func(d selenium.WebElement) {
			defer func() {
				<-sem
				wg.Done()
			}()

			branch, err := parseBranch(d)
			if err != nil {
				log.Warn().Msg(err.Error())
				return
			}

			if branch == (Branch{}) {
				return
			}

			b = append(b, branch)
		}(d)

	}

	wg.Wait()
	log.Debug().Msgf("parseBranches took %v", time.Since(t))

	return b, err
}

// parseBranch parses branch info from the HTML element.
func parseBranch(e selenium.WebElement) (Branch, error) {

	eUpdatedDate, err := e.FindElement(selenium.ByCSSSelector, ".cURBaH")
	if err != nil {
		return Branch{}, fmt.Errorf("failed find eUpdatedDate: %v", err)
	}

	sUpdatedDate, err := eUpdatedDate.Text()
	if err != nil {
		return Branch{}, fmt.Errorf("failed get text from eUpdatedDate: %v", err)
	}

	sUpdatedDate = sanitaze(sUpdatedDate)
	if len(sUpdatedDate) == 0 {
		return Branch{}, errors.New("sUpdatedDate is empty")
	}

	if u := strings.Split(sUpdatedDate, " "); len(u) >= 3 {
		sUpdatedDate = strings.Join(u[len(u)-2:], " ")
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return Branch{}, fmt.Errorf("failed to load location: %v", err)
	}

	updatedDate, err := time.ParseInLocation("02.01.2006 15:04", sUpdatedDate, loc)
	if err != nil {
		return Branch{}, fmt.Errorf("failed to parse sUpdatedDate: %v", err)
	}

	if time.Now().In(loc).Sub(updatedDate) > 24*time.Hour {
		return Branch{}, fmt.Errorf("updatedDate is out of date for 24 hours: %v", updatedDate)
	}

	eBank, err := e.FindElement(selenium.ByCSSSelector, ".dPnGDN")
	if err != nil {
		return Branch{}, fmt.Errorf("failed find sBank: %v", err)
	}

	sBank, err := eBank.Text()
	if err != nil {
		return Branch{}, fmt.Errorf("failed to get text from eBank: %v", err)
	}

	eSubway, err := e.FindElement(selenium.ByCSSSelector, ".eybsgm")
	if err != nil {
		return Branch{}, fmt.Errorf("failed find eSubway: %v", err)
	}

	sSubway := ""
	if eSubway != nil {
		sSubway, err = eSubway.Text()
		if err != nil {
			log.Warn().Msgf("failed to get text from eSubway: %v", err)
		}
	}

	eRates, err := e.FindElements(selenium.ByCSSSelector, ".fvORFF")
	if err != nil {
		return Branch{}, fmt.Errorf("failed find eRates: %v", err)
	}

	if len(eRates) != 2 {
		return Branch{}, fmt.Errorf("invalid eRates count: %d", len(eRates))
	}

	sBuyRate, err := eRates[0].Text()
	if err != nil {
		return Branch{}, fmt.Errorf("failed to get text from eRates[0]: %v", err)
	}

	sBuyRate = strings.Replace(strings.ReplaceAll(strings.ReplaceAll(sBuyRate, "—", "0"), ",", "."), " ₽", "", -1)
	buyRate, err := strconv.ParseFloat(sBuyRate, 64)
	if err != nil || buyRate <= 0 {
		return Branch{}, fmt.Errorf("failed to parse float from sBuyRate: %v", err)
	}

	sSellRate, _ := eRates[1].Text()
	if err != nil {
		return Branch{}, fmt.Errorf("failed to get text from eRates[1]: %v", err)
	}

	sSellRate = strings.Replace(strings.ReplaceAll(strings.ReplaceAll(sSellRate, "—", "0"), ",", "."), " ₽", "", -1)
	sellRate, err := strconv.ParseFloat(sSellRate, 64)
	if err != nil || sellRate <= 0 {
		return Branch{}, fmt.Errorf("failed to parse float from sSellRate: %v", err)
	}

	return newBranch(sBank, sanitaze(sSubway), buyRate, sellRate, updatedDate), nil
}

// sanitaize string.
func sanitaze(s string) string {
	if len(s) == 0 {
		return s
	}

	s = strings.Replace(s, "\n", "", -1)
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")

	return s
}
