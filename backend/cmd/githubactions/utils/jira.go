package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
)

// JiraIssueTransition will try to transit an issue to the given state
func JiraIssueTransition(issueKey string, transitionName string, jiraClient *jira.Client) error {

	// get available transitions
	transitions, _, err := jiraClient.Issue.GetTransitions(issueKey)
	if err != nil {
		return errors.Wrap(err, "get transitions failed")
	}

	// search for transition
	transitionID := ""
	for _, transition := range transitions {
		if transition.Name == transitionName {
			transitionID = transition.ID
		}
	}

	// handle invalid transition
	if transitionID == "" {
		return errors.New("transition not allowed on current state")
	}

	// transition
	_, errTransition := jiraClient.Issue.DoTransition(issueKey, transitionID)
	if errTransition != nil {
		return errors.Wrap(errTransition, "transition failed")
	}

	return nil
}

func Encode(updateVo interface{}) (io.ReadWriter, error) {
	// json encoding
	var buffer io.ReadWriter
	buffer = new(bytes.Buffer)
	errEncode := json.NewEncoder(buffer).Encode(updateVo)
	if errEncode != nil {
		return nil, errors.Wrap(errEncode, "JSON encoding failed")
	}
	return buffer, nil
}

// JiraIssueUpdate custom update request
func JiraIssueUpdate(issueKey string, buffer io.ReadWriter, jiraClient *jira.Client) (string, error) {

	// update request
	var result interface{}
	req, errReq := jiraClient.NewRawRequest(http.MethodPut, "/rest/api/2/issue/"+issueKey+"?notifyUsers=false", buffer)
	if errReq != nil {
		return "", errors.Wrap(errReq, "JIRA req init failed")
	}
	resp, errDo := jiraClient.Do(req, &result)
	if errDo == io.EOF {
		errDo = nil
	}
	if errDo != nil {
		return "", errors.Wrap(errDo, "JIRA req failed")
	}

	// read response
	responseBody, errReadAll := ioutil.ReadAll(resp.Body)
	if errReadAll != nil {
		return "", errors.Wrap(errReadAll, "JIRA response failed")
	}

	// successful?
	if http.StatusNoContent != resp.StatusCode {
		return string(responseBody), errors.New("unexpected http status: " + resp.Status)
	}

	return string(responseBody), nil
}
