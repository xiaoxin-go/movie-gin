package utils

import (
	"fmt"
	"github.com/tebeka/selenium"
	"strings"
	"time"
)

func NewFilm(name string, wd selenium.WebDriver)*film{
	result := &film{Name: name, Wd: wd}
	result.get()
	return result
}

type FilmData struct {
	Name        string
	Title       string
	Length      string
	Genres      []string
	Actresses   []Actress
	ReleaseDate time.Time
	ImageUrl string
	Links       []Link
	Images      []Image
}

type Link struct {
	Name      string
	Size      string
	Magnet    string
	ShareDate time.Time
}
type Image struct{
	Name string
	Url string
	SimpleUrl string
}

type film struct {
	Name  string
	Wd    selenium.WebDriver
	error error
}

func (m *film) Error() error {
	return m.error
}
type Actress struct{
	Name string
	Url string
}

func (m *film) get() {
	m.error = m.Wd.Get("https://www.javbus.com/" + m.Name)
}

func (m *film) title() (result string) {
	el, err := m.Wd.FindElement(selenium.ByCSSSelector, ".container>h3")
	if err != nil {
		m.error = fmt.Errorf("get title error: %s", err.Error())
		return
	}
	result, m.error = el.Text()
	return
}

func (m *film) Data() (result FilmData) {
	infoEl, err := m.Wd.FindElement(selenium.ByCSSSelector, ".container>.movie>.info")
	if err != nil {
		m.error = fmt.Errorf("get info error: %s", err.Error())
		return
	}
	text, err := infoEl.Text()
	if err != nil {
		m.error = fmt.Errorf("get info error: %s", err.Error())
		return
	}
	result = FilmData{
		Name: m.Name,
		Title:     m.title(),
		Genres:    m.genres(),
		Actresses: m.actresses(),
		ImageUrl: m.imageUrl(),
		Links: m.links(),
		Images: m.images(),
		ReleaseDate: m.releaseDate(text),
		Length: m.length(text),
	}
	return
}

func (m *film) releaseDate(text string) time.Time {
	result := strings.Split(strings.Split(text, "發行日期:")[1], "\n")[0]
	result = strings.Trim(result, " ")
	return StrToDate(result)
}
func (m *film) length(text string) string {
	result := strings.Split(strings.Split(text, "長度:")[1], "\n")[0]
	return strings.Trim(result, " ")
}

func (m *film) genres() (result []string) {
	els, err := m.Wd.FindElements(selenium.ByCSSSelector, ".container>.movie>.info .genre>label")
	if err != nil{
		m.error = fmt.Errorf("get genres error: %s", err.Error())
		return
	}
	for _, el := range els{
		genre, err := el.Text()
		if err != nil{
			m.error = fmt.Errorf("get genres text error: %s", err.Error())
			return
		}
		result = append(result, genre)
	}
	return
}
func (m *film) actresses() (result []Actress) {

	els, err := m.Wd.FindElements(selenium.ByCSSSelector, ".container>.movie>.info>p:last-child a")
	if err != nil{
		m.error = fmt.Errorf("get genres error: %s", err.Error())
		return
	}
	for _, el := range els{
		name, err := el.Text()
		if err != nil{
			m.error = fmt.Errorf("get actress text error: %s", err.Error())
			return
		}
		url, err := el.GetAttribute("href")
		if err != nil{
			m.error = fmt.Errorf("get actress href error: %s", err.Error())
			return
		}
		result = append(result, Actress{Name: name, Url: url})
	}
	return
}
func (m *film) imageUrl()(result string){
	el, err := m.Wd.FindElement(selenium.ByCSSSelector, ".container>.movie .bigImage")
	if err != nil {
		m.error = fmt.Errorf("get big image url error: %s", err.Error())
		return
	}
	result, err = el.GetAttribute("href")
	if err != nil{
		m.error = fmt.Errorf("get big image url href error: %s", err.Error())
	}
	return
}
func (m *film) links()(result []Link){
	els, err := m.Wd.FindElements(selenium.ByCSSSelector, "#magnet-table>tr")
	if err != nil{
		m.error = fmt.Errorf("get links error: %s", err.Error())
		return
	}
	for _, el := range els{
		tds, err := el.FindElements(selenium.ByTagName, "td")
		if err != nil{
			m.error = fmt.Errorf("get links error: %s", err.Error())
			return
		}
		if len(tds) < 3{
			m.error = fmt.Errorf("get links error: len tds < 3")
			return
		}
		nameEl, err := tds[0].FindElement(selenium.ByTagName, "a")
		if err != nil{
			m.error = fmt.Errorf("get links tds 0 err: %s", err.Error())
			return
		}
		name, err := nameEl.Text()
		if err != nil{
			m.error = fmt.Errorf("get links error: %s", err.Error())
			return
		}
		magnet, err := nameEl.GetAttribute("href")
		if err != nil{
			m.error = fmt.Errorf("get links error: %s", err.Error())
			return
		}
		sizeEl, err := tds[1].FindElement(selenium.ByTagName, "a")
		if err != nil{
			m.error = fmt.Errorf("get links tds 1 error: %s", err.Error())
			return
		}
		size, err := sizeEl.Text()
		if err != nil{
			m.error = fmt.Errorf("get links error: %s", err.Error())
			return
		}
		shareDateEl, err := tds[2].FindElement(selenium.ByTagName, "a")
		shareDate, err := shareDateEl.Text()
		if err != nil{
			m.error = fmt.Errorf("get links tds 2 error: %s", err.Error())
			return
		}
		result = append(result, Link{Name: name, Size: size, Magnet: magnet, ShareDate: StrToDate(shareDate)})
	}
	return
}
func (m *film) images()(result []Image){
	els, err := m.Wd.FindElements(selenium.ByCSSSelector, "#sample-waterfall>a")
	if err != nil{
		m.error = fmt.Errorf("get images error: %s", err.Error())
		return
	}
	for i, el := range els{
		url, err := el.GetAttribute("href")
		if err != nil{
			m.error = fmt.Errorf("get images url error: %s", err.Error())
			return
		}
		name := fmt.Sprintf("%s-%d", m.Name, i+1)
		simpleEl, err := el.FindElement(selenium.ByTagName, "img")
		if err != nil{
			m.error = fmt.Errorf("get images simpleEl error: %s", err.Error())
			return
		}
		simpleUrl, err := simpleEl.GetAttribute("src")
		if err != nil{
			m.error = fmt.Errorf("get simple url error: %s", err.Error())
			return
		}
		result = append(result, Image{Name: name, Url: url, SimpleUrl: simpleUrl})
	}
	return
}
