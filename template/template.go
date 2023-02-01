package template

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/Mrpye/helm-api/lib"
)

func ParseInterfaceMap(model lib.InstallUpgradeRequest, val map[string]interface{}) (map[string]interface{}, error) {
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

func ParseTemplate(model lib.InstallUpgradeRequest, tpl string) (string, error) {
	//*********************
	//Create a function map
	//*********************
	funcMap := template.FuncMap{
		"base64enc":   lib.Base64EncString,
		"base64dec":   lib.Base64DecString,
		"gzip_base64": lib.GzipBase64,
		"lc":          strings.ToLower,
		"uc":          strings.ToUpper,
		"domain":      lib.GetDomainOrIP,
		"port":        lib.GetPortString,
		"port_int":    lib.GetPortInt,
		"clean":       lib.Clean,
		"concat":      lib.Concat,
		"replace":     strings.ReplaceAll,
		"contains":    lib.StringContainsStringListItem,
		"not":         lib.NOT,
		"or":          lib.OR,
		"and":         lib.AND,
		"plus":        lib.Plus,
		"minus":       lib.Minus,
		"multiply":    lib.Multiply,
		"divide":      lib.Divide,
	}

	//*****************
	//Pase the template
	//*****************
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
