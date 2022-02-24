package htmlops

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"io/ioutil"
	"log"
)

func CreateHtmlTemplate(tmpl string, data interface{}) (string, error) {
	t := template.New("template")

	var err error
	t, err = t.Parse(tmpl)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	result := tpl.String()
	return result, nil
}

func Base64Encode(data []byte) string {
	b64 := base64.StdEncoding.EncodeToString(data)
	return b64
}

func GetLocalHtml(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("GetLocalHtml failed: ", err)
		return []byte{}, err
	}

	return data, nil
}
