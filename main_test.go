package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI(t *testing.T) {
	suite.Run(t, &APISuite{})
}

type APISuite struct {
	suite.Suite

	client http.Client
}

func (s *APISuite) SetupSuite() {
	srv := NewServer()
	go func() {
		log.Printf("Start serving on %s", srv.Addr)
		log.Fatal(srv.ListenAndServe())
	}()
}

func (s *APISuite) TestCreateAndGet() {
	// setup:
	targetContent := []byte("biba kuka")
	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write(targetContent)
		s.Require().NoError(err)
	}))

	var key string
	s.Run("CreateShortLink", func() {
		// when:
		reqBody := io.NopCloser(strings.NewReader(fmt.Sprintf( /* language=json */ `{"url": "%s"}`, testServer.URL)))
		resp, err := s.client.Post("http://localhost:8080/urls", "application/json", reqBody)

		// then:
		s.Require().NoError(err)
		rawBody, err := ioutil.ReadAll(resp.Body)
		s.Require().NoError(err)
		var body map[string]string
		s.Require().NoError(json.Unmarshal(rawBody, &body))
		key = body["key"]
		s.Require().NotEmpty(key)
	})
}
