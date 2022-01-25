package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dungnh3/guardrails-assignment/api"

	"github.com/golang/protobuf/jsonpb"
)

const (
	testDataPath = "./data"

	repositoryUrlPath = "/api/v1/repositories"
)

// Test_EndToEnd_CreateRepository_Success be used to test api create repository
func (s *e2eTestSuite) Test_EndToEnd_1_CreateRepository_Success() {
	b, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", testDataPath, "create_repository.json"))
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPost, uri(repositoryUrlPath), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	s.Require().NoError(err)
	client := http.Client{}
	resp, err := client.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var a api.CreateRepositoryResponse
	err = jsonpb.Unmarshal(resp.Body, &a)
	s.Require().NoError(err)

	s.Equal("demo", a.Data.SourceRepository.Name)
	s.Equal("https://github.com/guardrailsio/backend-engineer-challenge", a.Data.SourceRepository.Link)
}

// Test_EndToEnd_CreateRepository_Success be used to test api get repository by id
func (s *e2eTestSuite) Test_EndToEnd_2_GetRepositoryByID_Success() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/%v", uri(repositoryUrlPath), 1 /*:id*/), nil)
	req.Header.Set("Content-Type", "application/json")
	s.Require().NoError(err)

	client := http.Client{}
	resp, err := client.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var a api.CreateRepositoryResponse
	err = jsonpb.Unmarshal(resp.Body, &a)
	s.Require().NoError(err)

	s.Equal("demo", a.Data.SourceRepository.Name)
	s.Equal("https://github.com/guardrailsio/backend-engineer-challenge", a.Data.SourceRepository.Link)
}

// Test_EndToEnd_CreateRepository_Success be used to test api update repository by id
func (s *e2eTestSuite) Test_EndToEnd_3_UpdateRepositoryByID_Success() {
	b, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", testDataPath, "update_repository.json"))
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPut, uri(repositoryUrlPath), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	s.Require().NoError(err)
	client := http.Client{}
	response, err := client.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%v/%v", uri(repositoryUrlPath), 1 /*:id*/), nil)
	req.Header.Set("Content-Type", "application/json")
	s.Require().NoError(err)

	resp, err := client.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var a api.CreateRepositoryResponse
	err = jsonpb.Unmarshal(resp.Body, &a)
	s.Require().NoError(err)

	s.Equal("abc", a.Data.SourceRepository.Name)
}

// Test_EndToEnd_CreateRepository_Success be used to test api remove repository by id
func (s *e2eTestSuite) Test_EndToEnd_4_RemoveRepositoryByID_Success() {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/%v", uri(repositoryUrlPath), 1 /*:id*/), nil)
	req.Header.Set("Content-Type", "application/json")
	s.Require().NoError(err)

	client := http.Client{}
	response, err := client.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
}
