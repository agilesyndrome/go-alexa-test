package main

import (
  "encoding/json"
  "io/ioutil"
  "fmt"
  "errors"
  "github.com/arienmalec/alexa-go"
  "github.com/agilesyndrome/go-alexa-dispatcher/dispatcher"
  "github.com/stretchr/testify/assert"
  "testing"
  "github.com/agilesyndrome/go-alexa-i18n/alexai18n"
)



var (
 TestCases map[string] ExpectedResponse = map[string] ExpectedResponse{}
)

func Add(tc string, e ExpectedResponse) {
  TestCases[tc] = e
}

func Similar(tc string, tc_dupe string) {
  TestCases[tc] = TestCases[tc_dupe]
}

//Simple Helper function to assume that e is nil
func check(e error) {
  if e != nil {
    panic(e)
  }
}


func checkResponse(testFile string) (alexa.Response, error) {
   test_data, json_err := ioutil.ReadFile(testFile)
   check(json_err)
   testRequest := new(alexa.Request)
   err := json.Unmarshal(test_data, &testRequest)
   response, err := dispatcher.Dispatch(*testRequest )
   return response, err

}

type ExpectedResponse struct {
  TestFile string
  Title string
  Text string
  Error error
  Culture string
  ShouldEndSession bool

  testRequest alexa.Request
}


func Expect(culture string, title string, text string, shouldEndSession bool) (ExpectedResponse) {

  e:= ExpectedResponse {
    Title : title,
    Text : text,
    Culture : culture,
    ShouldEndSession : shouldEndSession,
    Error : nil,
    testRequest : alexai18n.CultureRequest(culture),

  }
  return e
}




func ExpectError(culture string, title string, text string, shouldEndSession bool, errorText string) (ExpectedResponse) {
  e:= Expect(culture, title, text, shouldEndSession)
  e.Error = errors.New(errorText)
  return e
}

func assertResponse(e ExpectedResponse, t *testing.T) {
   response, err := checkResponse(e.TestFile)
   expectedAlexaResponse := alexa.NewSimpleResponse(e.Title, e.Text)
   assert.IsType(t, alexa.Response{}, response)
   assert.Equal(t, expectedAlexaResponse, response)
   assert.Equal(t, e.ShouldEndSession, response.Body.ShouldEndSession)
   assert.Equal(t, e.Error, err)
}

func Run(tc string, t *testing.T)(ExpectedResponse) {
  e:= TestCases[ tc ]
  fmt.Printf("Running %s...\n", tc)
  e.TestFile = fmt.Sprintf("../tests/data/%s.json", tc)
  assertResponse(e,t)
  return e
}

func RunTests(t *testing.T) {
  for k, _ := range TestCases {
    Run(k,t)
  }
}

func TestRun(t *testing.T) {
  RunTests(t)
}
