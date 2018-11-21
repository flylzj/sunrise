package spider

import (
	"bytes"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
)

func GetJsonData(url string, method string, headers map[string]string, body string) (*simplejson.Json, error){
	fmt.Println(url)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 4.4.2; HUAWEI MLA-AL10 Build/HUAWEIMLA-AL10) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36 Html5Plus/1.0")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("http error")
		fmt.Println(err)
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("io error")
		fmt.Println(err)
		return nil, err
	}

	jsonDate, err := simplejson.NewJson(responseBody)
	if err != nil {
		fmt.Println("json error")
		fmt.Println(err)
		return nil, err
	}
	return jsonDate, nil
}
