package nbtverify

import (
	"encoding/json"
	"fmt"
	"net/http"
	urllib "net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nbtca/nbtverify/nbtverify/utils"
)

func GetBaseLoginUrl(testUrl string, mobile bool) (bool, string, error) {
	data, err := utils.RequestGet(testUrl, mobile)
	if err != nil {
		return false, "", err
	}
	//if offline:
	//<script>top.self.location.href='http://10.80.253.2:9090/zportal/login?...'</script>
	//or else if online:
	//<html>
	// <head><title>503 Service Temporarily Unavailable</title></head>
	// </body>
	// </html>
	str := string(data)
	if strings.HasPrefix(str, "<script") {
		start := strings.Index(str, "'")
		end := strings.LastIndex(str, "'")
		if start == -1 || end == -1 {
			return false, "", fmt.Errorf("can't find url in response data: %s", str)
		}
		urlRaw := str[start+1 : end]
		fixUrl := strings.ReplaceAll(urlRaw, " ", "")
		return true, fixUrl, nil
	} else {
		return false, "", nil
	}
}

type RequestLoginResult struct {
	Message  string `json:"message"`
	NextPage string `json:"nextPage"`
	Result   string `json:"result"`
}
type LoginInfo struct {
	Username string
	Password string
	AsMobile bool
}

type LoginResult struct {
	AlreadyOnline bool           `json:"alreadyOnline"`
	Success       bool           `json:"success"`
	Message       string         `json:"message"`
	NextPage      string         `json:"nextPage"`
	BaseUrl       string         `json:"baseUrl"`
	Cookies       []*http.Cookie `json:"cookies"`
	Mobile        bool           `json:"mobile"`
}
type OnlineDetail struct {
	Welcome       string            `json:"welcome"`
	Account       string            `json:"account"`
	LogoutUrl     string            `json:"logoutUrl"`
	Form          map[string]string `json:"form"`
	UserIP        string            `json:"userIp"`
	UserMac       string            `json:"userMac"`
	UserName      string            `json:"userName"`
	DeviceIP      string            `json:"deviceIp"`
	IsMacFastAuth bool              `json:"isMacFastAuth"`
	Mobile        bool              `json:"mobile"`
}
type LogoutResult struct {
	mobile  bool
	Message string
}

func (result *LoginResult) GetDetail() (*OnlineDetail, error) {
	bytes, err := utils.RequestGetReferer(result.NextPage, result.Mobile, result.BaseUrl, result.Cookies...)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(bytes))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	welcome := utils.GetTextContent(doc, "body > div.m_box > ul", //mobile
		"body > div.box_mod > div.contentsbox > h1", //pc
	)
	account := utils.GetTextContent(doc, "body > ul", //mobile
		"body > div.box_mod > div.zhxx_box", //pc
	)
	form := utils.FormToMap(doc, "#hidden_form")
	logoutUrl, err := utils.ChangeUrlPath(result.NextPage, "/zportal/logout")
	if err != nil {
		return nil, err
	}
	return &OnlineDetail{
		Form:          form,
		Welcome:       welcome,
		Account:       account,
		UserIP:        form["userIp"],
		UserMac:       utils.ConvertMac(form["userMac"]),
		UserName:      form["userName"],
		DeviceIP:      form["deviceIp"],
		IsMacFastAuth: form["isMacFastAuth"] == "true",
		Mobile:        result.Mobile,
		LogoutUrl:     logoutUrl,
	}, nil
}
func (status *OnlineDetail) Logout() (*LogoutResult, error) {
	res, _, err := utils.RequestPostForm(status.LogoutUrl, status.Form, status.Mobile)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(res))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	message := utils.GetTextContent(doc,
		"body > div > ul > li",                     //mobile
		"body > div.box_mod > div.zhxx_box > span", //pc
	)
	return &LogoutResult{
		mobile:  status.Mobile,
		Message: message,
	}, nil
}
func invokeLogin(baseUrl string, form map[string]string, mobile bool) (*LoginResult, error) {
	// baseAddr := strings.Split(baseUrl, "/")
	// newUrl := baseAddr[0] + "//" + baseAddr[2] + "/zportal/login/do"
	url, err := urllib.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	url.Path = "/zportal/login/do"
	url.RawQuery = ""
	newUrl := url.String()
	res, cookies, err := utils.RequestPostForm(newUrl, form, mobile)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Login Result:", string(res))
	//{
	//   "message": "您已经在线！请不要重复认证",
	//   "nextPage": "goToAuthResult",
	//   "result": "online"
	// }
	// or
	//{
	//   "message": "",
	//   "nextPage": "goToAuthResult",
	//   "result": "success"
	// }
	text := string(res)
	result := RequestLoginResult{}
	err = json.Unmarshal([]byte(text), &result)
	if err != nil {
		return nil, err
	}
	url.Path = "/zportal/" + result.NextPage
	return &LoginResult{
		AlreadyOnline: result.Result == "online",
		Success:       result.Result == "success",
		Message:       result.Message,
		NextPage:      url.String(),
		BaseUrl:       baseUrl,
		Cookies:       cookies,
		Mobile:        mobile,
	}, nil
}
func Login(baseUrl string, info LoginInfo) (*LoginResult, error) {
	bytes, err := utils.RequestGet(baseUrl, info.AsMobile)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(bytes))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	form := utils.FormToMap(doc, "#login_form")
	form["username"] = info.Username
	form["pwd"] = info.Password
	return invokeLogin(baseUrl, form, info.AsMobile)
}
