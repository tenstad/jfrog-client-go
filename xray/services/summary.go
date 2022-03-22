package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	servicesutils "github.com/tenstad/jfrog-client-go/artifactory/services/utils"
	"github.com/tenstad/jfrog-client-go/auth"
	"github.com/tenstad/jfrog-client-go/http/jfroghttpclient"
	"github.com/tenstad/jfrog-client-go/utils"
	"github.com/tenstad/jfrog-client-go/utils/errorutils"
)

const (
	summaryAPI = "api/v2/summary/"
)

func (ss *SummaryService) getSummeryUrl() string {
	return ss.XrayDetails.GetUrl() + summaryAPI
}

// SummaryService returns the https client and Xray details
type SummaryService struct {
	client      *jfroghttpclient.JfrogHttpClient
	XrayDetails auth.ServiceDetails
}

// NewSummaryService creates a new service to retrieve the version of Xray
func NewSummaryService(client *jfroghttpclient.JfrogHttpClient) *SummaryService {
	return &SummaryService{client: client}
}

func (ss *SummaryService) GetBuildSummary(params XrayBuildParams) (*SummaryResponse, error) {
	httpDetails := ss.XrayDetails.CreateHttpClientDetails()
	url := fmt.Sprintf("%sbuild?build_name=%s&build_number=%s", ss.getSummeryUrl(), params.BuildName, params.BuildNumber)
	if params.Project != "" {
		url += "&" + projectKeyQueryParam + params.Project
	}
	resp, body, _, err := ss.client.SendGet(url, true, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	var summaryResponse SummaryResponse
	err = json.Unmarshal(body, &summaryResponse)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	if summaryResponse.Errors != nil && len(summaryResponse.Errors) > 0 {
		return nil, errorutils.CheckErrorf("Getting build-summery for build: %s failed with error: %s", summaryResponse.Errors[0].Identifier, summaryResponse.Errors[0].Error)
	}
	return &summaryResponse, nil
}

func (ss *SummaryService) GetArtifactSummary(params ArtifactSummaryParams) (*ArtifactSummaryResponse, error) {
	httpDetails := ss.XrayDetails.CreateHttpClientDetails()
	servicesutils.SetContentType("application/json", &httpDetails.Headers)
	// TODO: Check if required
	// utils.AddAuthHeaders(httpDetails.Headers, ss.XrayDetails)

	requestBody, err := json.Marshal(params)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	url := fmt.Sprintf("%sartifact", ss.getSummeryUrl())
	resp, body, err := ss.client.SendPost(url, requestBody, &httpDetails)
	if err != nil {
		return nil, err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return nil, errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, utils.IndentJson(body)))
	}
	var response ArtifactSummaryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}
	if response.Errors != nil && len(response.Errors) > 0 {
		return nil, errorutils.CheckErrorf("Getting artifact-summery for artifact: %s failed with error: %s", response.Errors[0].Identifier, response.Errors[0].Error)
	}
	return &response, nil
}

type ArtifactSummaryParams struct {
	Checksums []string `json:"checksums,omitempty"`
	Paths     []string `json:"paths,omitempty"`
}

type ArtifactSummaryResponse struct {
	Artifacts []Artifact `json:"artifacts,omitempty"`
	Errors    []Error    `json:"errors,omitempty"`
}

type Artifact struct {
	General General `json:"general,omitempty"`
	Issues  []Issue `json:"issues,omitempty"`
	// TODO: Create License struct with correct fields for api/v2/summary/artifact endpoint
	// Licenses []License `json:"licenses,omitempty"`
}

type General struct {
	ComponentId string `json:"component_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Path        string `json:"path,omitempty"`
	PkgType     string `json:"pkg_type,omitempty"`
	Sha256      string `json:"sha256,omitempty"`
}

type SummaryResponse struct {
	Issues []Issue
	Errors []Error
}

type Issue struct {
	IssueId     string             `json:"issue_id,omitempty"`
	Summary     string             `json:"summary,omitempty"`
	Description string             `json:"description,omitempty"`
	IssueType   string             `json:"issue_type,omitempty"`
	Severity    string             `json:"severity,omitempty"`
	Provider    string             `json:"provider,omitempty"`
	Cves        []SummeryCve       `json:"cves,omitempty"`
	Created     string             `json:"created,omitempty"`
	ImpactPath  []string           `json:"impact_path,omitempty"`
	Components  []SummeryComponent `json:"components,omitempty"`
}

type Error struct {
	Error      string `json:"error,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

type SummeryCve struct {
	Id          string `json:"cve,omitempty"`
	CvssV2Score string `json:"cvss_v2,omitempty"`
	CvssV3Score string `json:"cvss_v3,omitempty"`
}

type SummeryComponent struct {
	ComponentId   string   `json:"component_id,omitempty"`
	FixedVersions []string `json:"fixed_versions,omitempty"`
}
