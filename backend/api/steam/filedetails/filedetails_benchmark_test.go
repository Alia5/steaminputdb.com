package filedetails_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alia5/steaminputdb.com/api/steam/filedetails"
	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/danielgtaylor/huma/v2/humatest"
)

func BenchmarkFileDetailsInfo(b *testing.B) {

	type testCase struct {
		name        string
		useMemCache bool
	}

	testCases := []testCase{
		{name: "With MemCache", useMemCache: true},
		{name: "Without MemCache", useMemCache: false},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {

			bWithoutLogging := BSupressedLogs{b}
			_, api := humatest.New(bWithoutLogging)
			filedetails.RegisterRoute(api, tc.useMemCache)

			mockSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/IPublishedFileService/GetDetails/v1/" {
					http.Error(w, "wrong path", http.StatusNotFound)
					return
				}

				resp := &steamapi.CPublishedFile_GetDetails_Response{
					Publishedfiledetails: []*steamapi.PublishedFileDetails{
						{
							Publishedfileid:          new(uint64(123)),
							FileType:                 new(uint32(12)),
							Title:                    new("My Config"),
							FileDescription:          new("Cool config"),
							Filename:                 new("controller.vdf"),
							FileUrl:                  new("https://cdn.steamusercontent.com/ugc/abc"),
							FileSize:                 new(uint64(2048)),
							Creator:                  new(uint64(76561198000000000)),
							TimeCreated:              new(uint32(1700000000)),
							TimeUpdated:              new(uint32(1700000123)),
							LifetimePlaytime:         new(uint64(600)),
							LifetimePlaytimeSessions: new(uint64(42)),
							LifetimeSubscriptions:    new(uint32(9001)),
							VoteData: &steamapi.PublishedFileDetails_VoteData{
								Score:     new(float32(4.25)),
								VotesUp:   new(uint32(100)),
								VotesDown: new(uint32(5)),
							},
							Tags: []*steamapi.PublishedFileDetails_Tag{
								{Tag: new(string(steamtypes.ControllerTypeXboxOne))},
								{Tag: new("controller_native")},
							},
							Kvtags: []*steamapi.PublishedFileDetails_KVTag{
								{Key: new("app"), Value: new("440")},
							},
						},
					},
				}

				w.WriteHeader(http.StatusOK)
				w.Write(mustMarshalProto(b, resp))
			}))
			defer mockSrv.Close()
			redirectSteamAPITo(b, mockSrv.URL)

			for b.Loop() {
				resp := api.Get("/v1/steam/filedetails?file_id=250900")
				b.StopTimer()
				if resp.Body == nil {
					b.Fatal("response body is nil")
				}
				if resp.Code != 200 {
					b.Fatalf("expected status 200, got %d", resp.Code)
				}
				b.StartTimer()
			}

		})
	}

}

type BSupressedLogs struct {
	*testing.B
}

func (b BSupressedLogs) Log(args ...any) {}
