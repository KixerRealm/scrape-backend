package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"scrape-backend/src/main/internal/database"
	"sync"
	"time"
)

func (apiCfg *apiConfig) handlerCreateBugReport(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		ImageFilename string `json:"image_filename"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	bugReport, err := apiCfg.DB.CreateBugReport(r.Context(), database.CreateBugReportParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		Title:         params.Title,
		Description:   params.Description,
		ImageFilename: params.ImageFilename,
		UserID:        user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %s", err))
		return
	}

	err = createNewBugIssue(databaseBugReportToBugReport(bugReport))
	if err != nil {
		fmt.Println("Error:", err)
	}
	respondWithJSON(w, 200, databaseBugReportToBugReport(bugReport))
}

func (apiCfg *apiConfig) handlerGetBugReportsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	bugReports, err := apiCfg.DB.GetBugReportsByUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %s", err))
		return
	}
	respondWithJSON(w, 200, databaseBugReportsToBugReports(bugReports))
}

func (apiCfg *apiConfig) handlerGetAllBugReports(w http.ResponseWriter, r *http.Request) {
	bugReports, err := apiCfg.DB.GetBugReports(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %s", err))
		return
	}
	respondWithJSON(w, 200, databaseBugReportsToBugReports(bugReports))
}

func createNewBugIssue(bugReport BugReport) error {

	apiURL := "https://api.linear.app/graphql"
	accessToken := "enter-linear-API-key"

	var wg sync.WaitGroup

	teamIDCh := make(chan string, 1)
	projectIDCh := make(chan string, 1)
	labelIDCh := make(chan string, 1)
	stateIDCh := make(chan string, 1)

	wg.Add(4)

	go func() {
		defer wg.Done()
		id, err := getResponseId(accessToken, apiURL, "teams")
		if err != nil {
			fmt.Println("Error getting team ID:", err)
			return
		}
		teamIDCh <- id
	}()

	go func() {
		defer wg.Done()
		id, err := getResponseId(accessToken, apiURL, "projects")
		if err != nil {
			fmt.Println("Error getting project ID:", err)
			return
		}
		projectIDCh <- id
	}()

	go func() {
		defer wg.Done()
		id, err := getResponseId(accessToken, apiURL, "issueLabels")
		if err != nil {
			fmt.Println("Error getting label ID:", err)
			return
		}
		labelIDCh <- id
	}()

	go func() {
		defer wg.Done()
		id, err := getResponseId(accessToken, apiURL, "workflowStates")
		if err != nil {
			fmt.Println("Error getting state ID:", err)
			return
		}
		stateIDCh <- id
	}()

	go func() {
		wg.Wait()
		close(teamIDCh)
		close(projectIDCh)
		close(labelIDCh)
		close(stateIDCh)
	}()

	// Retrieve results from channels
	teamID := <-teamIDCh
	projectID := <-projectIDCh
	labelID := <-labelIDCh
	stateID := <-stateIDCh

	mutation := `mutation {
		issueCreate(input: {
			title: "` + bugReport.Title + `",
			description: "` + bugReport.Description + `",
			teamId: "` + teamID + `",
			projectId: "` + projectID + `",
			labelIds: "` + labelID + `"
			stateId: "` + stateID + `"
		}) {
			success
			issue {
				id
				title
			}
		}
	}`

	jsonData := map[string]string{"query": mutation}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", accessToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	fmt.Println("Response Body:", string(body))

	return nil
}

func getResponseId(accessToken string, apiURL string, typeQuery string) (string, error) {

	var query string
	if typeQuery == "teams" {
		query = `{
		"query": "query Teams { teams { nodes { id name } } }"
	}`
	} else if typeQuery == "projects" {
		query = `{
		"query": "query Projects { projects { nodes { id name } } }"
	}`
	} else if typeQuery == "issueLabels" {
		query = `{
		"query": "query IssueLabels { issueLabels { nodes { id name } } }"
	}`
	} else if typeQuery == "workflowStates" {
		query = `{
		"query": "query WorkflowStates { workflowStates { nodes { id name } } }"
	}`
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer([]byte(query)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err.Error(), nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", accessToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err.Error(), nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err.Error(), nil
	}

	fmt.Println("Response Body:", string(body))

	var id string
	if typeQuery == "teams" {
		type Response struct {
			Data struct {
				Teams struct {
					Nodes []struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"nodes"`
				} `json:"teams"`
			} `json:"data"`
		}
		bodyReader1 := bytes.NewReader(body)

		var teamResponse Response
		if err := json.NewDecoder(bodyReader1).Decode(&teamResponse); err != nil {
			fmt.Println("Error decoding response body:", err)
			return err.Error(), nil
		}

		if len(teamResponse.Data.Teams.Nodes) > 0 {
			id = teamResponse.Data.Teams.Nodes[0].ID
			fmt.Println("Extracted ID:", id)
		} else {
			fmt.Println("No nodes found in the response.")
		}
	} else if typeQuery == "projects" {
		type Response struct {
			Data struct {
				Projects struct {
					Nodes []struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"nodes"`
				} `json:"projects"`
			} `json:"data"`
		}
		bodyReader1 := bytes.NewReader(body)

		var projectResponse Response
		if err := json.NewDecoder(bodyReader1).Decode(&projectResponse); err != nil {
			fmt.Println("Error decoding response body:", err)
			return err.Error(), nil
		}

		if len(projectResponse.Data.Projects.Nodes) > 0 {
			id = projectResponse.Data.Projects.Nodes[5].ID
			fmt.Println("Extracted ID:", id)
		} else {
			fmt.Println("No nodes found in the response.")
		}
	} else if typeQuery == "issueLabels" {
		type Response struct {
			Data struct {
				Labels struct {
					Nodes []struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"nodes"`
				} `json:"issueLabels"`
			} `json:"data"`
		}
		bodyReader1 := bytes.NewReader(body)

		var labelsResponse Response
		if err := json.NewDecoder(bodyReader1).Decode(&labelsResponse); err != nil {
			fmt.Println("Error decoding response body:", err)
			return err.Error(), nil
		}

		if len(labelsResponse.Data.Labels.Nodes) > 0 {
			id = labelsResponse.Data.Labels.Nodes[2].ID
			fmt.Println("Extracted ID:", id)
		} else {
			fmt.Println("No nodes found in the response.")
		}
	} else if typeQuery == "workflowStates" {
		type Response struct {
			Data struct {
				States struct {
					Nodes []struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"nodes"`
				} `json:"workflowStates"`
			} `json:"data"`
		}
		bodyReader1 := bytes.NewReader(body)

		var statesResponse Response
		if err := json.NewDecoder(bodyReader1).Decode(&statesResponse); err != nil {
			fmt.Println("Error decoding response body:", err)
			return err.Error(), nil
		}

		if len(statesResponse.Data.States.Nodes) > 0 {
			id = statesResponse.Data.States.Nodes[1].ID
			fmt.Println("Extracted ID:", id)
		} else {
			fmt.Println("No nodes found in the response.")
		}
	}
	return id, nil
}
