package helper

import (
	"bytes"
	"net/http"
)

func PostReq(postUrl string, body []byte) (response *http.Response, err error) {
	posturl := postUrl

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	return res, err
}
