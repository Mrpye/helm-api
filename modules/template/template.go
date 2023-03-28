package template

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/Mrpye/golib/encrypt"
	"github.com/Mrpye/golib/math"
	"github.com/Mrpye/golib/str"

	"github.com/Mrpye/helm-api/modules/body_types"
)

func ParseInterfaceMap(model body_types.InstallUpgradeRequest, val map[string]interface{}) (map[string]interface{}, error) {
	var err error
	for k, v := range val {
		_ = k
		switch data_val := v.(type) {
		case string:
			parsed_str, err := ParseTemplate(model, string(data_val))
			if err != nil {
				return nil, err
			}
			val[k] = parsed_str
		case map[string]interface{}:
			ParseInterfaceMap(model, data_val)
		}
	}
	return val, err
}

func ParseTemplate(model body_types.InstallUpgradeRequest, tpl string) (string, error) {
	//*********************
	//Create a function map
	//*********************
	funcMap := template.FuncMap{
		"base64enc":   encrypt.Base64EncString,
		"base64dec":   encrypt.Base64DecString,
		"gzip_base64": str.GzipBase64,
		"lc":          strings.ToLower,
		"uc":          strings.ToUpper,
		"domain":      str.GetDomainOrIP,
		"port":        str.GetPortString,
		"port_int":    str.GetPortInt,
		"clean":       str.Clean,
		"concat":      str.Concat,
		"replace":     strings.ReplaceAll,
		"contains":    str.CommaListContainsString,
		"not":         math.NOT,
		"or":          math.OR,
		"and":         math.AND,
		"plus":        math.Plus,
		"minus":       math.Minus,
		"multiply":    math.Multiply,
		"divide":      math.Divide,
	}

	//*****************
	//Pase the template
	//*****************
	tpl = strings.ReplaceAll(tpl, "<%", "{{")
	tpl = strings.ReplaceAll(tpl, "%>", "}}")
	tmpl, err := template.New("CodeRun").Funcs(funcMap).Parse(tpl)
	if err != nil {
		return "", err
	}

	//**************************************
	//Run the template to verify the output.
	//**************************************
	var tpl_buffer bytes.Buffer
	err = tmpl.Execute(&tpl_buffer, model)
	if err != nil {
		return "", err
	}

	return tpl_buffer.String(), nil
}
