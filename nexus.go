package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Artifact public for xml parsing
type Artifact struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
	RepoID     string `xml:"repoId"`
}

// Artifacts public for xml parsing
type Artifacts struct {
	Artifacts []Artifact `xml:"artifact"`
}

// SearchResponse public for xml parsing
type SearchResponse struct {
	XMLName xml.Name  `xml:"search-results"`
	Results Artifacts `xml:"data"`
}

type artifactRequest struct {
	Repository string
	Group      string
	Artifact   string
}

//TODO: use struct for args
//func latest(baseURL, repository, groupID, artifactID string) (*Artifact, error) {
func latest(baseURL string, request *artifactRequest) (*Artifact, error) {
	artifacts, err := search(baseURL, request)
	if err != nil {
		return nil, err
	}
	maxVersion := ""
	var newest Artifact
	for _, artifact := range artifacts {
		if artifact.RepoID == request.Repository && artifact.Version > maxVersion {
			fmt.Println(artifact.RepoID, artifact.GroupID, artifact.ArtifactID, artifact.Version, request.Repository)
			newest = artifact
			maxVersion = artifact.Version
		}
	}
	if maxVersion == "" {
		return nil, nil
	}
	return &newest, nil
}

//TODO: check for 'tooManyResults'
func search(baseURL string, request *artifactRequest) ([]Artifact, error) {
	//func search(baseURL, repository, groupID, artifactID string) ([]Artifact, error) {
	url := fmt.Sprintf("%s/data_index?g=%s&a=%s", baseURL, request.Group, request.Artifact)
	fmt.Println("Requesting ", url)
	body, err := getResponse(url)
	if err != nil {
		return nil, err
	}
	var searchResponse SearchResponse
	if err := xml.Unmarshal(body, &searchResponse); err != nil {
		return nil, err
	}
	artifacts := searchResponse.Results.Artifacts
	filtered := make([]Artifact, len(artifacts))
	var i = 0
	for _, artifact := range artifacts {
		if artifact.RepoID == request.Repository {
			filtered[i] = artifact
			i++
		}
	}
	return filtered[:i], nil
}

func getResponse(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parseArtifactRequest(r *http.Request) (*artifactRequest, error) {
	q := r.URL.Query()
	for _, param := range []string{"group", "artifact", "repository"} {
		if len(q.Get(param)) == 0 {
			return nil, errors.New(param + " is required")
		}
	}
	return &artifactRequest{
		Repository: q.Get("repository"),
		Group:      q.Get("group"),
		Artifact:   q.Get("artifact"),
	}, nil
}
