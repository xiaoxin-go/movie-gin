package utils

import (
	"fmt"
	"github.com/tebeka/selenium"
	"strconv"
	"strings"
)


func NewActressController(url string, wd selenium.WebDriver)*actress{
	return &actress{Url: url, wd: wd}
}

type actress struct{
	Url string
	error error
	wd selenium.WebDriver
}


func (t *actress) Error()error{
	return t.error
}

type ActressData struct{
	Name string
	Age int
	Birthday string
	Cup string
	Height string
	Films []string
}

func (t *actress) Data() (result ActressData){
	err := t.wd.Get(t.Url)
	if err != nil{
		t.error = fmt.Errorf("open url error=> %s", err.Error())
		return
	}
	info, err := t.wd.FindElement(selenium.ByCSSSelector, ".avatar-box>.photo-info")
	if err != nil{
		t.error = fmt.Errorf("get info css error=> %s", err.Error())
		return
	}
	infoText, err := info.Text()
	if err != nil{
		t.error = fmt.Errorf("get info text error=> %s", err.Error())
		return
	}
	result = ActressData{
		Name: t.name(infoText),
		Age: t.age(infoText),
		Height: t.height(infoText),
		Cup: t.cup(infoText),
		Birthday: t.birthday(infoText),
		Films: t.films(),
	}
	return result
}
func (t *actress) name(text string)string{
	return strings.Split(text, "\n")[0]
}
func (t *actress) age(text string)int{
	value := t.labelValue(text, "年齡:")
	if value == ""{
		return 0
	}
	result, _ := strconv.Atoi(value)
	return result
}
func (t *actress) height(text string)string{
	return t.labelValue(text, "身高:")
}
func (t *actress) birthday(text string)string{
	return t.labelValue(text, "生日:")
}
func (t *actress) cup(text string)string{
	return t.labelValue(text, "罩杯:")
}
func (t *actress) labelValue(text, label string)string{
	if !strings.Contains(text, label){
		return ""
	}
	result := strings.Split(strings.Split(text, label)[1], "\n")[0]
	return strings.Trim(result, " ")
}
func (t *actress) page()int{
	els, err := t.wd.FindElements(selenium.ByCSSSelector, ".pagination>li>a")
	if err != nil || len(els) < 2{
		return 1
	}
	fmt.Println("=======> ", els)
	el := els[len(els) - 2]
	pageStr, err := el.Text()
	fmt.Println("page => ", pageStr)
	if err != nil{
		t.error = fmt.Errorf("get pageStr error=> %s", err.Error())
		return 1
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil{
		t.error = fmt.Errorf("page str to int error=> %s", err.Error())
		return 1
	}
	return page
}
func (t *actress) films()(result []string){
	result = make([]string, 0)
	page := t.page()
	currentPage := 1
	for {
		filmEls, err := t.wd.FindElements(selenium.ByCSSSelector, "#waterfall>.item>.movie-box")
		if err != nil{
			t.error = fmt.Errorf("get film css error=> %s", err.Error())
			return
		}
		for _, el := range filmEls{
			filmHref, err := el.GetAttribute("href")
			if err != nil{
				t.error = fmt.Errorf("get film href error=> %s", err.Error())
				return
			}
			filmName := strings.Split(filmHref, "/")[3]
			result = append(result, filmName)
		}
		if currentPage >= page{
			break
		}
		currentPage += 1
		url := fmt.Sprintf("%s/%d", t.Url, currentPage)
		err = t.wd.Get(url)
		if err != nil{
			t.error = fmt.Errorf("get get url error=> %s, %s", err.Error(), url)
			return
		}
	}
	return
}