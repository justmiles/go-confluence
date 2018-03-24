package confluence

import (
	"encoding/json"
	"github.com/google/go-querystring/query"
	log "github.com/sirupsen/logrus"
)

// GetContent Returns all content in a Confluence instance.
// https://developer.atlassian.com/cloud/confluence/rest/#api-content-get
func (client Client) GetContent(qp *GetContentQueryParameters) ([]Content, error) {
	var queryParams string
	if qp != nil {
		v, _ := query.Values(qp)
		queryParams = v.Encode()
	}
	body, err := client.request("GET", "/rest/api/content", queryParams, "")
	if err != nil {
		return nil, err
	}
	var contentResponse ContentResponse
	err = json.Unmarshal(body, &contentResponse)
	if err != nil {
		log.Error("Unable to unmarshal ContentResponse. Received: '", string(body), "'")
	}
	return contentResponse.Results, err
}

// GetContentQueryParameters query parameters for GetContent
type GetContentQueryParameters struct {
	QueryParameters
	Limit      int    `url:"limit,omitempty"`
	Orderby    string `url:"orderby,omitempty"`
	PostingDay string `url:"postingDay,omitempty"`
	Spacekey   string `url:"spaceKey,omitempty"`
	Start      int    `url:"start,omitempty"`
	Title      string `url:"title,omitempty"`
	Trigger    string `url:"trigger,omitempty"`
	Type       string `url:"type,omitempty"`
}

// CreateContent creates a new piece of content or publishes an existing draft.
// https://developer.atlassian.com/cloud/confluence/rest/#api-content-post
func (client Client) CreateContent(bp *CreateContentBodyParameters, qp *QueryParameters) (Content, error) {
	var res Content
	var queryParams string
	if qp != nil {
		v, _ := query.Values(qp)
		queryParams = v.Encode()
	}

	byteString, err := json.Marshal(bp)
	if err != nil {
		log.Error("Unable to marshal body. Received: '", err, "'")
	}

	body, err := client.request("POST", "/rest/api/content", queryParams, string(byteString))
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Error(body)
		log.Error(err)
		log.Error("Unable to unmarshal CreateContentResponse. Received: '", string(body), "'")
	}
	return res, err
}

// CreateContentBodyParameters query parameters for CreateContent
type CreateContentBodyParameters struct {
	Ancestors []struct {
		ID string `json:"id,omitempty"`
	} `json:"ancestors,omitempty"`
	Body struct {
		AnonymousExportView *struct {
		} `json:"anonymous_export_view,omitempty"`
		Editor2 *struct {
		} `json:"editor2,omitempty"`
		ExportView *struct {
		} `json:"export_view,omitempty"`
		Storage *struct {
		} `json:"storage,omitempty"`
		StyledView *struct {
		} `json:"styled_view,omitempty"`
		View *struct {
			Representation string `json:"representation,omitempty"`
			Value          string `json:"value,omitempty"`
		} `json:"view,omitempty"`
	} `json:"body,omitempty"`
	ID    string `json:"id,omitempty"`
	Space struct {
		Key string `json:"key,omitempty"`
	} `json:"space,omitempty"`
	Status string `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Type   string `json:"type,omitempty"`
}

func (client Client) DeleteContent(content Content) error {
	_, err := client.request("DELETE", "/rest/api/content/"+content.ID, "", "")
	return err
}

// ContentResponse represents the data returned from the Confluence API
type ContentResponse struct {
	Results []Content `json:"results"`
}

// Content represents the data returned from the Confluence API
type Content struct {
	Client
	ID                  string `json:"id"`
	Type                string `json:"type"`
	Status              string `json:"status"`
	Title               string `json:"title"`
	URL                 string `json:"url,omitempty"`
	MacroRenderedOutput struct {
	} `json:"macroRenderedOutput,omitempty"`
	Links struct {
		Self   string `json:"self"`
		Tinyui string `json:"tinyui"`
		Editui string `json:"editui"`
		Webui  string `json:"webui"`
	} `json:"_links"`
}
