package configs

import (
	"strings"

	"github.com/Alia5/steaminputdb.com/steamapi"
)

type RankBy string

const (
	RankedByVote                     RankBy = "vote"
	RankedByPublication              RankBy = "publication"
	RankedByTrend                    RankBy = "trend"
	RankedByTotalUniqueSubscriptions RankBy = "subscriptions"
	RankedByTotalVotesAsc            RankBy = "votes_asc"
	RankedByVotesUp                  RankBy = "votes_up"
	RankedByTextSearch               RankBy = "text_search"
	RankedByPlaytimeTrend            RankBy = "playtime_trend"
	RankedByTotalPlaytime            RankBy = "total_playtime"
	RankedByAveragePlaytimeTrend     RankBy = "avg_playtime_trend"
	RankedByLifetimeAveragePlaytime  RankBy = "lifetime_avg_playtime"
	RankedByPlaytimeSessionsTrend    RankBy = "playtime_sessions_trend"
	RankedByLifetimePlaytimeSessions RankBy = "lifetime_playtime_sessions"
	RankedByLastUpdated              RankBy = "updated"
)

func (s *RankBy) PublishedFileQueryType() steamapi.EPublishedFileQueryType {
	switch strings.ToLower(string(*s)) {
	case "vote":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByVote
	case "publication":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByPublicationDate
	case "trend":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByTrend
	case "subscriptions":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByTotalUniqueSubscriptions
	case "votes_asc":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByTotalVotesAsc
	case "votes_up":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByVotesUp
	case "text_search":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByTextSearch
	case "playtime_trend":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByPlaytimeTrend
	case "total_playtime":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByTotalPlaytime
	case "avg_playtime_trend":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByAveragePlaytimeTrend
	case "lifetime_avg_playtime":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByLifetimeAveragePlaytime
	case "playtime_sessions_trend":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByPlaytimeSessionsTrend
	case "lifetime_playtime_sessions":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByLifetimePlaytimeSessions
	case "updated":
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByLastUpdatedDate
	default:
		return steamapi.EPublishedFileQueryType_k_PublishedFileQueryType_RankedByVote
	}

}
