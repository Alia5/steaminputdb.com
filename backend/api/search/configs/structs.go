package configs

import (
	"time"

	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/steamapi"
)

type ConfigPagination struct {
	Limit uint32 `json:"limit,omitempty,omitzero" default:"5" max:"100" doc:"Maximum number of results to return"`
	Page  uint32 `json:"page,omitempty,omitzero" doc:"Page number for paginated results"`
}

type ConfigFilter struct {
	AppID        string   `json:"app_id,omitempty,omitzero" maxLength:"100" doc:"Only search configurations for this specific Steam App ID (Non-Steam Games use their name as AppID)"`
	Tags         []string `json:"tags,omitempty,omitzero" doc:"Filter results to only include configurations with all of the specified tags"`
	ExcludedTags []string `json:"excluded_tags,omitempty,omitzero" doc:"Filter results to only include configurations that do not have any of the specified tags"`
	Creator      *string  `json:"creator,omitempty,omitzero" doc:"Filter results to only include configurations created by the specified Steam User ID"`
}

type ConfigRank struct {
	By             RankBy `json:"by,omitempty,omitzero" default:"vote" enum:"vote,publication,trend,subscriptions,votes_asc,votes_up,text_search,playtime_trend,total_playtime,avg_playtime_trend,lifetime_avg_playtime,playtime_sessions_trend,lifetime_playtime_sessions,updated" doc:"Criterion to rank search results by"`
	TrendingPeriod uint32 `json:"trending_period,omitempty,omitzero" default:"30" doc:"Number of days to consider for trending rank"` // TODO: use?
}

type ConfigInclude struct {
	Votes bool `json:"votes,omitempty,omitzero" default:"false" doc:"Include vote data"`
	Tags  bool `json:"tags,omitempty,omitzero" default:"false" doc:"Include tags"`
}

type ConfigQueryBody struct {
	ConfigPagination
	Filter    ConfigFilter  `json:"filter,omitzero"`
	Rank      ConfigRank    `json:"rank,omitzero"`
	Include   ConfigInclude `json:"include,omitzero"`
	QueryText string        `json:"query_text" required:"true" maxLength:"100" doc:"The search query string"`
	Raw       bool          `json:"raw,omitempty,omitzero" default:"false" doc:"Return raw Steam API response"`
}

type Request struct {
	Body ConfigQueryBody `additionalProperties:"false"`
}

type ConfigItem struct {
	Title                    *string                    `json:"title" example:"My Awesome Config"`
	Description              *string                    `json:"description" example:"An awesome configuration for an awesome game"`
	AppID                    *uint32                    `json:"app_id,omitempty,omitzero" example:"420"`
	AppIDString              *string                    `json:"app_id_string" example:"420" doc:"AppID as string or name for Non-Steam Games"`
	FileID                   *uint64                    `json:"file_id" example:"1234567890"`
	FileName                 *string                    `json:"file_name" example:"controllerconfig.vdf"`
	FileURL                  *string                    `json:"file_url" example:"https://cdn.steamusercontent.com/ugc/UGCID/HASH"`
	FileSize                 *uint64                    `json:"file_size" example:"2048" doc:"File size in bytes"`
	CreatorID                string                     `json:"creator_id" example:"76561198000000000" doc:"Steam User ID of the configuration creator"`
	ControllerType           *steamtypes.ControllerType `json:"controller_type,omitempty,omitzero" example:"controller_steamcontroller_gordon" doc:"Type of controller this configuration is designed for"`
	ControllerTypeNice       string                     `json:"controller_type_nice,omitempty,omitzero" example:"Steam Controller (2015)" doc:"Human-friendly name of the controller type"`
	ControllerNative         bool                       `json:"controller_native,omitempty,omitzero" example:"true" doc:"Unsure"`
	TimeCreated              time.Time                  `json:"time_created"`
	TimeUpdated              time.Time                  `json:"time_updated"`
	LifetimePlaytimeSeconds  *uint64                    `json:"lifetime_playtime_seconds" example:"4474585092" doc:"Total playtime in seconds"`
	LifetimePlaytimeSessions *uint64                    `json:"lifetime_playtime_sessions" example:"1234" doc:"Total number of playtime sessions"`
	Subscriptions            *uint32                    `json:"subscriptions" example:"1500" doc:"Total number of Downloads"`

	Votes struct {
		Score *float32 `json:"score" example:"4.5" doc:"Average vote score"`
		Up    *uint32  `json:"up" example:"100" doc:"Number of upvotes"`
		Down  *uint32  `json:"down" example:"20" doc:"Number of downvotes"`
	} `json:"votes,omitzero"`

	Tags *[]string `json:"tags,omitempty,omitzero" example:"[\"controller_xboxone\", \"feature_gyro\"]" doc:"List of tags associated with the configuration"`
}

type ConfigsResponse struct {
	Total int          `json:"total,omitempty"`
	Items []ConfigItem `json:"items,omitempty"`
}

type raw steamapi.CPublishedFile_QueryFiles_Response

type responseBody interface {
	configQueryResponse()
}

func (r *raw) configQueryResponse()             {}
func (r *ConfigsResponse) configQueryResponse() {}

type SearchResponse struct {
	Body responseBody
}
