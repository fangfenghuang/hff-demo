package http

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url string) ([]byte, error) {

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New(string(b))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", b)
	}

	return b, nil
}

func Delete(url string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return b, nil
}

func Put(url string, body []byte) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !(resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK) {
		return nil, fmt.Errorf("error: %s", b)
	}

	return b, nil
}
