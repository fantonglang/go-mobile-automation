package apis

import (
	"encoding/json"

	m "github.com/fantonglang/go-mobile-automation"
)

type Service struct {
	req  *m.SharedRequest
	path string
}

const (
	SERVICE_UIAUTOMATOR = "uiautomator"
)

func NewService(name string, req *m.SharedRequest) *Service {
	return &Service{
		req:  req,
		path: "/services/" + name,
	}
}

type ServiceResponse struct {
	Success     bool   `json:"success"`
	Running     bool   `json:"running"`
	Description string `json:"description"`
}

func (s *Service) Start() (*ServiceResponse, error) {
	inputBytes := make([]byte, 0)
	bytes, err := s.req.Post(s.path, inputBytes)
	if err != nil {
		return nil, err
	}
	r := new(ServiceResponse)
	err = json.Unmarshal(bytes, r)
	return r, err
}

func (s *Service) Stop() (*ServiceResponse, error) {
	bytes, err := s.req.Delete(s.path)
	if err != nil {
		return nil, err
	}
	r := new(ServiceResponse)
	err = json.Unmarshal(bytes, r)
	return r, err
}

func (s *Service) Running() (bool, error) {
	bytes, err := s.req.Get(s.path)
	if err != nil {
		return false, err
	}
	r := new(ServiceResponse)
	err = json.Unmarshal(bytes, r)
	if err != nil {
		return false, err
	}
	return r.Running, nil
}
