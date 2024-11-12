package utils

import (
	"encoding/json"
	urllib "net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func FindAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
func GetTextContent(doc *goquery.Document, selectors ...string) string {
	text := ""
	for _, selector := range selectors {
		text += doc.Find(selector).Text()
	}
	text = strings.Replace(text, "\t", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, "\n\n", "\n", -1)
	text = strings.Replace(text, "\n\n", "\n", -1)
	text = strings.Replace(text, "\n", " ", -1)
	return text
}
func FormToMap(doc *goquery.Document, selector string) map[string]string {
	selection := doc.Find(selector)
	form := make(map[string]string)
	for _, input := range selection.Find("input").Nodes {
		name := ""
		value := ""
		name = FindAttr(input, "name")
		value = FindAttr(input, "value")
		if name != "" {
			form[name] = value
		}
	}
	return form
}
func ConvertMac(hexString string) string {
	hexString = strings.ToUpper(hexString)
	var parts []string
	for i := 0; i < len(hexString); i += 2 {
		parts = append(parts, hexString[i:i+2])
	}
	return strings.Join(parts, ":")
}
func ChangeUrlPath(url string, path string) (string, error) {
	newUrl, err := urllib.Parse(url)
	if err != nil {
		return "", err
	}
	newUrl.Path = path
	newUrl.RawQuery = ""
	return newUrl.String(), nil
}

func RemoveComments(bytes []byte) []byte {
	var result []byte
	length := len(bytes)
	inString := false
	for i := 0; i < length; i++ {
		if bytes[i] == '"' {
			inString = !inString
		}
		if !inString && bytes[i] == '/' {
			if i < length-1 {
				if bytes[i+1] == '/' {
					i++
					for i < length-1 && bytes[i+1] != '\n' {
						i++
					}
					continue
				} else if bytes[i+1] == '*' {
					i += 2
					for i < length-1 && !(bytes[i] == '*' && bytes[i+1] == '/') {
						i++
					}
					i++
					continue
				}
			}
		}
		result = append(result, bytes[i])
	}
	return result
}
func FileNotExists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func SaveJson(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(data)
}
