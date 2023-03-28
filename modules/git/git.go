package git

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Mrpye/golib/file"
	"github.com/Mrpye/golib/net"
	"github.com/Mrpye/helm-api/modules/body_types"
)

func DownloadFromGitLab(url string, token string) (string, error) {
	//url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s?ref=%s", m.Host, project_data, File, branch)
	headers := []net.Header{
		{Key: "Content-Type", Value: "application/json"},
		{Key: "Accept", Value: "application/json"},
		{Key: "PRIVATE-TOKEN", Value: token},
	}
	res, _, err := net.CallApi(url, "GET", headers, nil, false)
	if string(res) == "{\"message\":\"401 Unauthorized\"}" {
		return string(res), errors.New("invalid auth token")
	}
	return string(res), err
}

func DownloadFromGitHub(url string, token string) (string, error) {
	//url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s?ref=%s", m.Host, project_data, File, branch)
	headers := []net.Header{
		{Key: "Content-Type", Value: "application/json"},
		{Key: "Accept", Value: "*/*"},
		{Key: "Authorization", Value: token},
	}
	res, _, err := net.CallApi(url, "GET", headers, nil, false)
	return string(res), err
}

func downloadGitFile(url string, token string) (string, error) {
	if strings.Contains(url, "/api/v4/projects") {
		//gitlab
		data, err := DownloadFromGitLab(url, token)
		if err != nil {
			return "", err
		}
		return data, nil
	} else if strings.Contains(url, "/api/v4/projects") {
		//github
		data, err := DownloadFromGitHub(url, token)
		if err != nil {
			return "", err
		}
		return data, nil
	} else {
		return "", errors.New("invalid url")
	}
}

func LoadHelmConfig(config_path string, answer_path string, params map[string]string, release_name string, namespace string, token string) (*body_types.InstallUpgradeRequest, error) {
	//**************************
	//Load the answer file first
	//**************************
	var obj body_types.InstallUpgradeRequest
	var data string
	var err error
	if answer_path != "" {
		var answer_obj map[string]string

		if strings.HasPrefix(answer_path, "http") {
			data, err = downloadGitFile(answer_path, token)
			if err != nil {
				return nil, err
			}
		} else if file.FileExists(answer_path) {
			data, err = file.ReadFileToString(answer_path)
			if err != nil {
				return nil, err
			}
		}
		json.Unmarshal([]byte(data), &answer_obj)
		if params == nil {
			params = answer_obj
		}
	}
	//*************************
	//Load the config file next
	//*************************
	if strings.HasPrefix(config_path, "http") {
		data, err = downloadGitFile(config_path, token)
		if err != nil {
			return nil, err
		}
	} else if file.FileExists(config_path) {
		data, err = file.ReadFileToString(config_path)
		if err != nil {
			return nil, err
		}
	}

	json.Unmarshal([]byte(data), &obj)
	if params != nil {
		obj.Params = params
	}
	if release_name != "" {
		obj.ReleaseName = release_name
	}
	if namespace != "" {
		obj.Namespace = namespace
	}
	return &obj, nil
}
