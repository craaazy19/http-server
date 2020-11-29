package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ParseRequest struct {
	URL []string `json:"url"`
}

type ParseTask struct {
	URL          string `json:"url"`
	ResponseBody []byte `json:"responseBody"`
}

type ParseResponse []*ParseTask

func (r *Router) Parse(w http.ResponseWriter, req *http.Request) {
	r.shutdownMutex.RLock()
	defer r.shutdownMutex.RUnlock()

	request := new(ParseRequest)
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
		return
	}

	err = json.Unmarshal(reqBody, request)
	if err != nil {
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
		return
	}

	if len(request.URL) == 0 {
		internalError := http.StatusBadRequest
		http.Error(w, "No urls to parse", internalError)
		return
	}
	if len(request.URL) > 20 {
		internalError := http.StatusBadRequest
		http.Error(w, "The amount of urls to parse is greater than 20", internalError)
		return
	}

	httpClient := &http.Client{
		Timeout: clientWaitTimeout,
	}

	response := make(ParseResponse, 0)
	for _, url := range request.URL {
		responseBody, err := requestURL(httpClient, url)
		if err != nil {
			internalError := http.StatusBadRequest
			http.Error(w, fmt.Sprintf("Error when parsing '%s': %s", url, err.Error()), internalError)
			return
		}
		response = append(response, &ParseTask{URL: url, ResponseBody: responseBody})
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
		return
	}

	_, err = fmt.Fprint(w, string(responseJSON))
	if err != nil {
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
		return
	}
}

func requestURL(httpClient *http.Client, url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response is %d: '%s'", resp.StatusCode, responseBody)
	}

	return responseBody, nil
}
