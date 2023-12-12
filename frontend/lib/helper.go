package h

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

	return res, err
}

func PutReq(postUrl string, body []byte) (response *http.Response, err error) {
	posturl := postUrl

	r, err := http.NewRequest("PUT", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	return res, err
}

func GetReq(getUrl string, body []byte) (response *http.Response, err error) {
	r, err := http.NewRequest("GET", getUrl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	return res, err
}
