package helix

import "fmt"

// SegmentType A segment configuration type
type ExtensionSegmentType string

// Types of segments datastores for the configuration service
const (
	ExtensionConfigrationBroadcasterSegment ExtensionSegmentType = "broadcaster"
	ExtensionConfigurationDeveloperSegment  ExtensionSegmentType = "developer"
	ExtensionConfigurationGlobalSegment     ExtensionSegmentType = "global"
)

func (s ExtensionSegmentType) String() string {
	return string(s)
}

type ExtensionSetConfigurationParams struct {
	Segment     ExtensionSegmentType `json:"segment"`
	ExtensionID string               `json:"extension_id"`
	// BroadcasterID is only populated if segment is of type 'developer' || 'broadcaster'
	BroadcasterID string `json:"broadcaster_id,omitempty"`
	Version       string `json:"version"`
	Content       string `json:"content"`
}

type ExtensionConfigurationSegment struct {
	Segment ExtensionSegmentType `json:"segment"`
	Version string               `json:"version"`
	Content string               `json:"content"`
}

type ExtensionGetConfigurationParams struct {
	ExtensionID   string                 `query:"extension_id"`
	BroadcasterID string                 `query:"broadcaster_id"`
	Segments      []ExtensionSegmentType `query:"segment"`
}

type ExtensionSetRequiredConfigurationParams struct {
	BroadcasterID         string `query:"broadcaster_id" json:"-"`
	ExtensionID           string `json:"extension_id"`
	RequiredConfiguration string `json:"required_version"`
	ExtensionVersion      string `json:"extension_version"`
	ConfigurationVersion  string `json:"configuration_version"`
}

type ExtensionSetRequiredConfigurationResponse struct {
	ResponseCommon
}

type ExtensionGetConfigurationSegmentResponse struct {
	ResponseCommon
	Data ManyExtensionConfigurationSegments
}

type ManyExtensionConfigurationSegments struct {
	Segments []ExtensionConfigurationSegment `json:"data"`
}

type ExtensionSetConfigurationResponse struct {
	ResponseCommon
}

// https://dev.twitch.tv/docs/extensions/reference/#set-extension-configuration-segment
func (c *Client) SetExtensionSegmentConfig(params *ExtensionSetConfigurationParams) (*ExtensionSetConfigurationResponse, error) {
	if params.BroadcasterID != "" {
		switch params.Segment {
		case ExtensionConfigurationDeveloperSegment, ExtensionConfigrationBroadcasterSegment:
		default:
			return nil, fmt.Errorf("error: developer or broadcaster extension configuration segment type must be provided for broadcasters")
		}
	}

	resp, err := c.putAsJSON("/extensions/configurations", &ManyPolls{}, params)
	if err != nil {
		return nil, err
	}

	setExtCnfgResp := &ExtensionSetConfigurationResponse{}
	resp.HydrateResponseCommon(&setExtCnfgResp.ResponseCommon)

	return setExtCnfgResp, nil
}

func (c *Client) GetExtensionConfigurationSegment(params *ExtensionGetConfigurationParams) (*ExtensionGetConfigurationSegmentResponse, error) {

	if params.BroadcasterID != "" {
		for _, segment := range params.Segments {
			switch segment {
			case ExtensionConfigurationDeveloperSegment, ExtensionConfigrationBroadcasterSegment:
			default:
				return nil, fmt.Errorf("error: only developer or broadcaster extension configuration segment type must be provided for broadcasters")
			}
		}
	}

	resp, err := c.get("/extensions/configurations", &ManyExtensionConfigurationSegments{}, params)
	if err != nil {
		return nil, err
	}

	extCfgSegResp := &ExtensionGetConfigurationSegmentResponse{}
	resp.HydrateResponseCommon(&extCfgSegResp.ResponseCommon)
	extCfgSegResp.Data.Segments = resp.Data.(*ManyExtensionConfigurationSegments).Segments

	return extCfgSegResp, nil
}

func (c *Client) SetExtensionRequiredConfiguration(params *ExtensionSetRequiredConfigurationParams) (*ExtensionSetRequiredConfigurationResponse, error) {

	resp, err := c.putAsJSON("/extensions/configurations/required_configuration", &ExtensionSetRequiredConfigurationResponse{}, params)
	if err != nil {
		return nil, err
	}

	extReqCfgResp := &ExtensionSetRequiredConfigurationResponse{}
	resp.HydrateResponseCommon(&extReqCfgResp.ResponseCommon)

	return extReqCfgResp, nil
}
