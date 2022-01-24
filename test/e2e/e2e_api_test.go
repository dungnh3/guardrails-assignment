package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	testDataPath = "./data"

	apiCreateRepository = "/api/v1/repositories"
)

// Test_EndToEnd_CreatePutSession_Success be used to test api put session
func (s *e2eTestSuite) Test_EndToEnd_CreateRepository_Success() {
	b, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", testDataPath, "create_repository.json"))
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPost, uri(apiCreateRepository), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	s.Require().NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
	//
	//defer res.Body.Close()
	//
	//content, err := ioutil.ReadAll(res.Body)
	//s.Require().NoError(err)
	//
	//data := make(map[string]interface{})
	//err = json.Unmarshal(content, &data)
	//s.Require().NoError(err)
}
