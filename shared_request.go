package openatxclientgo

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func init() {
	http.DefaultClient.Timeout = 10 * time.Second
}

type SharedRequest struct {
	BaseUrl string
	onError func() error
}

func NewSharedRequest(baseUrl string) *SharedRequest {
	return &SharedRequest{
		BaseUrl: baseUrl,
		onError: nil,
	}
}

func (r *SharedRequest) Post(path string, data []byte) ([]byte, error) {
	resp, err := http.Post(r.BaseUrl+path, "application/json", bytes.NewBuffer(data))
	if err != nil {
		if r.onError != nil {
			if _err := r.onError(); _err == nil {
				time.Sleep(10 * time.Second)
				resp, err = http.Post(r.BaseUrl+path, "application/json", bytes.NewBuffer(data))
			}
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, errors.New("[status:" + strconv.Itoa(resp.StatusCode) + "] " + string(bytes))
	}
	return bytes, nil
}

func (r *SharedRequest) Get(path string) ([]byte, error) {
	resp, err := http.Get(r.BaseUrl + path)
	if err != nil {
		if r.onError != nil {
			if _err := r.onError(); _err == nil {
				time.Sleep(10 * time.Second)
				resp, err = http.Get(r.BaseUrl + path)
			}
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, errors.New("[status:" + strconv.Itoa(resp.StatusCode) + "] " + string(bytes))
	}
	return bytes, nil
}

func (r *SharedRequest) Delete(path string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", r.BaseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if r.onError != nil {
			if _err := r.onError(); _err == nil {
				time.Sleep(10 * time.Second)
				resp, err = http.DefaultClient.Do(req)
			}
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, errors.New("[status:" + strconv.Itoa(resp.StatusCode) + "] " + string(bytes))
	}
	return bytes, nil
}

func (r *SharedRequest) GetWithTimeout(path string, timeout time.Duration) ([]byte, bool, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(r.BaseUrl + path)
	if err != nil {
		if _err, ok := err.(*url.Error); ok {
			if _err.Timeout() {
				return nil, true, err
			}
		}
		return nil, false, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, false, errors.New("[status:" + strconv.Itoa(resp.StatusCode) + "] " + string(bytes))
	}
	return bytes, false, nil
}

func (r *SharedRequest) PostWithTimeout(path string, data []byte, contentType string, timeout time.Duration) ([]byte, bool, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Post(r.BaseUrl+path, contentType, bytes.NewBuffer(data))
	if err != nil {
		if _err, ok := err.(*url.Error); ok {
			if _err.Timeout() {
				return nil, true, err
			}
		}
		return nil, false, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, false, errors.New("[status:" + strconv.Itoa(resp.StatusCode) + "] " + string(bytes))
	}
	return bytes, false, nil
}
