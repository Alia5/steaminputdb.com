package filedetails_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Alia5/steaminputdb.com/api/steam/filedetails"
	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestFileDetails(t *testing.T) {
	successMock := func(t *testing.T) *httptest.Server {
		t.Helper()
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			w.Write(mustMarshalProto(t, resp))
		}))
	}

	type testCase struct {
		name           string
		path           string
		setupMock      func(t *testing.T) *httptest.Server
		expectedStatus int
		expectedBody   string
		assertBody     func(t *testing.T, body []byte)
		contains       string
	}

	testCases := []testCase{
		{
			name:           "SUCCESS_PROCESSED",
			path:           "/v1/steam/filedetails?file_id=123",
			setupMock:      successMock,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"title": "My Config",
				"description": "Cool config",
				"app_id": 440,
				"app_id_string": "440",
				"file_id": 123,
				"file_name": "controller.vdf",
				"file_url": "https://cdn.steamusercontent.com/ugc/abc",
				"file_size": 2048,
				"creator_id": "76561198000000000",
				"controller_type": "controller_xboxone",
				"controller_type_nice": "Xbox One",
				"controller_native": true,
				"time_created": "2023-11-14T22:13:20Z",
				"time_updated": "2023-11-14T22:15:23Z",
				"lifetime_playtime_seconds": 600,
				"lifetime_playtime_sessions": 42,
				"subscriptions": 9001,
				"votes": {"score": 4.25, "up": 100, "down": 5},
				"tags": ["controller_xboxone", "controller_native"],
				"playtime_seconds": null,
				"playtime_sessions": null
			}`,
		},
		{
			name:           "SUCCESS_RAW",
			path:           "/v1/steam/filedetails?file_id=123&raw=true",
			setupMock:      successMock,
			expectedStatus: http.StatusOK,
			assertBody: func(t *testing.T, body []byte) {
				t.Helper()
				var got steamapi.CPublishedFile_GetDetails_Response
				require.NoError(t, json.Unmarshal(body, &got))
				require.Len(t, got.Publishedfiledetails, 1)
				require.NotNil(t, got.Publishedfiledetails[0])
				assert.Equal(t, uint64(123), got.Publishedfiledetails[0].GetPublishedfileid())
			},
		},
		{
			name:           "MISSING_FILE_ID",
			path:           "/v1/steam/filedetails",
			expectedStatus: http.StatusUnprocessableEntity,
			contains:       "file_id",
		},
		{
			name: "STEAM_API_HTTP_ERROR",
			path: "/v1/steam/filedetails?file_id=123",
			setupMock: func(t *testing.T) *httptest.Server {
				t.Helper()
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte("nope"))
				}))
			},
			expectedStatus: http.StatusBadGateway,
			contains:       "failed to get steam file details",
		},
		{
			name: "FILE_NOT_FOUND_EMPTY_LIST",
			path: "/v1/steam/filedetails?file_id=123",
			setupMock: func(t *testing.T) *httptest.Server {
				t.Helper()
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/IPublishedFileService/GetDetails/v1/" {
						http.Error(w, "wrong path", http.StatusNotFound)
						return
					}
					resp := &steamapi.CPublishedFile_GetDetails_Response{Publishedfiledetails: []*steamapi.PublishedFileDetails{}}
					w.WriteHeader(http.StatusOK)
					w.Write(mustMarshalProto(t, resp))
				}))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"detail":"file not found", "status":404, "title":"Not Found"}`,
		},
		{
			name: "FILE_NOT_FOUND_WRONG_TYPE",
			path: "/v1/steam/filedetails?file_id=123",
			setupMock: func(t *testing.T) *httptest.Server {
				t.Helper()
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/IPublishedFileService/GetDetails/v1/" {
						http.Error(w, "wrong path", http.StatusNotFound)
						return
					}
					fileID := uint64(123)
					wrongType := uint32(1)
					resp := &steamapi.CPublishedFile_GetDetails_Response{
						Publishedfiledetails: []*steamapi.PublishedFileDetails{
							{Publishedfileid: &fileID, FileType: &wrongType},
						},
					}
					w.WriteHeader(http.StatusOK)
					w.Write(mustMarshalProto(t, resp))
				}))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"detail":"file not found", "status":404, "title":"Not Found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, api := humatest.New(t)
			filedetails.RegisterRoute(api)

			if tc.setupMock != nil {
				steamAPIServer := tc.setupMock(t)
				defer steamAPIServer.Close()
				redirectSteamAPITo(t, steamAPIServer.URL)
			}

			resp := api.Get(tc.path)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.contains != "" {
				assert.True(t, strings.Contains(resp.Body.String(), tc.contains), "response body should contain %q, got: %s", tc.contains, resp.Body.String())
			}
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, resp.Body.String())
			}
			if tc.assertBody != nil {
				tc.assertBody(t, resp.Body.Bytes())
			}

		})
	}

}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func redirectSteamAPITo(t testing.TB, targetBaseURL string) {
	t.Helper()

	orig := http.DefaultClient
	origSteam := steamapi.DefaultClient

	u, err := url.Parse(targetBaseURL)
	if err != nil {
		t.Fatalf("failed to parse target base URL: %v", err)
	}

	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "api.steampowered.com" {
				req.URL.Scheme = u.Scheme
				req.URL.Host = u.Host
			}
			return http.DefaultTransport.RoundTrip(req)
		}),
	}
	steamapi.DefaultClient = steamapi.NewClientWithBaseURL("test-key", targetBaseURL)

	t.Cleanup(func() {
		http.DefaultClient = orig
		steamapi.DefaultClient = origSteam
	})
}

func mustMarshalProto(t testing.TB, msg proto.Message) []byte {
	t.Helper()
	b, err := proto.Marshal(msg)
	require.NoError(t, err)
	return b
}
