package appinfo_test

import (
	"context"
	"testing"
	"time"

	appinfodal "github.com/Alia5/steaminputdb.com/db/dal/appinfo"
	"github.com/Alia5/steaminputdb.com/db/models"
	sidbtest "github.com/Alia5/steaminputdb.com/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const testAppID uint32 = 250900

func fullInfo(appID uint32) *models.AppInfo {
	return &models.AppInfo{
		AppID:                   appID,
		ControllerSupportRating: models.ControllerSupportRatingGold,
		ControllerSupportNotes:  new("notes"),
		MixedInputInfo: &models.AppMixedInputInfo{
			MixedInputSupport:      models.MixedInputWithMod,
			MixedInputGlyphFlicker: true,
			MixedInputNotes:        "mixed notes",
			MixedInputModLinks: []*models.MixedInputModLinks{
				{
					Mod: "https://example.com/a",
				},
				{
					Mod: "https://example.com/b",
				},
			},
		},
		SteamInputAPISupport: &models.AppSteamInputAPISupport{
			SIAPICameraSupport: models.SteamInputCameraSupportFull,
			SIAPIPixelsPer360:  "1234",
			Notes:              "siapi notes",
			SIAPITypes: []*models.AppSIAPITypes{
				{
					Type: models.SteamInputAPIInGameActions,
				},
				{
					Type: models.SteamInputAPIActionSets,
				},
			},
		},
		Glyphs: &models.AppGlyphs{
			AutoGlyphDetect:   true,
			ManualGlyphSelect: true,
			GlyphNotes:        "glyph notes",
			GlyphCtrlSupport: []*models.AppGlyphCtrlSupport{
				{
					ControllerType: "controller_ps5",
					Notes:          "ps5",
				},
				{
					ControllerType: "controller_xboxone",
				},
			},
			GlyphTags: []*models.AppGlyphTag{
				{
					Tag: models.AppGlyphTagSteamInputAPIButtons,
				},
				{
					Tag: models.AppGlyphTagInGameButtons,
				},
			},
		},
		HWFeatures: []*models.AppHWFeatures{
			{
				HWFeature:               models.HWFeatureRumble,
				HWFeatureControllerType: models.HWFeatureControllerFamilyCommon,
				Notes:                   "rumble",
			},
			{
				HWFeature:               models.HWFeatureHDHaptics,
				HWFeatureControllerType: models.HWFeatureControllerFamilySteam,
			},
			{
				HWFeature:               models.HWFeatureHDHaptics,
				HWFeatureControllerType: models.HWFeatureControllerFamilyPlayStation,
			},
		},
	}
}

func TestUpdateSteamInputDBInfo(t *testing.T) {
	type testCase struct {
		name string
		// app (and a second, fully filled app) is not inserted
		noApp bool
		// applied using UpdateSteamInputDBInfo before update
		initial *models.AppInfo
		update  func(c context.Context, dal appinfodal.DAL) error
		// only steaminputdb infos are compared, without timestamps / AppIDs of relations
		expected    *models.AppInfo
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "APP_NOT_FOUND",
			noApp: true,
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateSteamInputDBInfo(c, fullInfo(testAppID))
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name:  "RATING_APP_NOT_FOUND",
			noApp: true,
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateControllerSupportRating(c, testAppID, models.ControllerSupportRatingGold, nil)
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name:  "NOTES_APP_NOT_FOUND",
			noApp: true,
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateControllerSupportNotes(c, testAppID, new("notes"), nil)
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "CREATE_ALL",
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateSteamInputDBInfo(c, fullInfo(testAppID))
			},
			expected: fullInfo(testAppID),
		},
		{
			name:    "NIL_RELATIONS_KEEP_EXISTING_ROWS",
			initial: fullInfo(testAppID),
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateSteamInputDBInfo(c, &models.AppInfo{
					AppID:                   testAppID,
					ControllerSupportRating: models.ControllerSupportRatingSilver,
					MixedInputInfo: &models.AppMixedInputInfo{
						MixedInputSupport: models.MixedInputUnsupported,
					},
				})
			},
			expected: func() *models.AppInfo {
				info := fullInfo(testAppID)
				info.ControllerSupportRating = models.ControllerSupportRatingSilver
				info.MixedInputInfo.MixedInputSupport = models.MixedInputUnsupported
				info.MixedInputInfo.MixedInputGlyphFlicker = false
				info.MixedInputInfo.MixedInputNotes = ""
				return info
			}(),
		},
		{
			name:    "EMPTY_RELATIONS_DELETE_EXISTING_ROWS",
			initial: fullInfo(testAppID),
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateSteamInputDBInfo(c, &models.AppInfo{
					AppID:                  testAppID,
					ControllerSupportNotes: new(""),
					MixedInputInfo: &models.AppMixedInputInfo{
						MixedInputModLinks: []*models.MixedInputModLinks{},
					},
					SteamInputAPISupport: &models.AppSteamInputAPISupport{
						SIAPITypes: []*models.AppSIAPITypes{},
					},
					Glyphs: &models.AppGlyphs{
						GlyphCtrlSupport: []*models.AppGlyphCtrlSupport{},
						GlyphTags:        []*models.AppGlyphTag{},
					},
					HWFeatures: []*models.AppHWFeatures{},
				})
			},
			expected: &models.AppInfo{
				AppID:                  testAppID,
				ControllerSupportNotes: new(""),
				MixedInputInfo:         &models.AppMixedInputInfo{},
				SteamInputAPISupport:   &models.AppSteamInputAPISupport{},
				Glyphs:                 &models.AppGlyphs{},
				HWFeatures:             []*models.AppHWFeatures{},
			},
		},
		{
			name:    "RATING",
			initial: fullInfo(testAppID),
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateControllerSupportRating(c, testAppID, models.ControllerSupportRatingWood, nil)
			},
			expected: func() *models.AppInfo {
				info := fullInfo(testAppID)
				info.ControllerSupportRating = models.ControllerSupportRatingWood
				return info
			}(),
		},
		{
			name:    "NOTES",
			initial: fullInfo(testAppID),
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateControllerSupportNotes(c, testAppID, new("changed"), nil)
			},
			expected: func() *models.AppInfo {
				info := fullInfo(testAppID)
				info.ControllerSupportNotes = new("changed")
				return info
			}(),
		},
		{
			name:    "NIL_NOTES_SET_NOTES_NULL",
			initial: fullInfo(testAppID),
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateControllerSupportNotes(c, testAppID, nil, nil)
			},
			expected: func() *models.AppInfo {
				info := fullInfo(testAppID)
				info.ControllerSupportNotes = nil
				return info
			}(),
		},
		{
			name:    "MIXED_INPUT",
			initial: fullInfo(testAppID),
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateMixedInputInfo(c, testAppID, &models.AppMixedInputInfo{
					MixedInputSupport: models.MixedInputSupported,
					MixedInputModLinks: []*models.MixedInputModLinks{
						{
							Mod: "https://example.com/b",
						},
						{
							Mod: "https://example.com/c",
						},
					},
				}, nil)
			},
			expected: func() *models.AppInfo {
				info := fullInfo(testAppID)
				info.MixedInputInfo = &models.AppMixedInputInfo{
					MixedInputSupport: models.MixedInputSupported,
					MixedInputModLinks: []*models.MixedInputModLinks{
						{
							Mod: "https://example.com/b",
						},
						{
							Mod: "https://example.com/c",
						},
					},
				}
				return info
			}(),
		},
		{
			name:    "GLYPHS",
			initial: fullInfo(testAppID),
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateGlyphs(c, testAppID, &models.AppGlyphs{
					GlyphNotes: "changed",
					GlyphCtrlSupport: []*models.AppGlyphCtrlSupport{
						{
							ControllerType: "controller_ps5",
							Notes:          "changed",
						},
						{
							ControllerType: "controller_ps4",
						},
					},
					GlyphTags: []*models.AppGlyphTag{
						{
							Tag: models.AppGlyphTagSteamInputAPITexts,
						},
					},
				}, nil)
			},
			expected: func() *models.AppInfo {
				info := fullInfo(testAppID)
				info.Glyphs = &models.AppGlyphs{
					GlyphNotes: "changed",
					GlyphCtrlSupport: []*models.AppGlyphCtrlSupport{
						{
							ControllerType: "controller_ps4",
						},
						{
							ControllerType: "controller_ps5",
							Notes:          "changed",
						},
					},
					GlyphTags: []*models.AppGlyphTag{
						{
							Tag: models.AppGlyphTagSteamInputAPITexts,
						},
					},
				}
				return info
			}(),
		},
		{
			name:    "STEAMINPUTAPI_SUPPORT",
			initial: fullInfo(testAppID),
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateSteamInputAPISupport(c, testAppID, &models.AppSteamInputAPISupport{
					SIAPICameraSupport: models.SteamInputCameraSupportPartial,
					SIAPITypes: []*models.AppSIAPITypes{
						{
							Type: models.SteamInputAPIActionSets,
						},
						{
							Type: models.SteamInputAPIXInputActions,
						},
					},
				}, nil)
			},
			expected: func() *models.AppInfo {
				info := fullInfo(testAppID)
				info.SteamInputAPISupport = &models.AppSteamInputAPISupport{
					SIAPICameraSupport: models.SteamInputCameraSupportPartial,
					SIAPITypes: []*models.AppSIAPITypes{
						{
							Type: models.SteamInputAPIActionSets,
						},
						{
							Type: models.SteamInputAPIXInputActions,
						},
					},
				}
				return info
			}(),
		},
		{
			name:    "HW_FEATURES",
			initial: fullInfo(testAppID),
			update: func(c context.Context, dal appinfodal.DAL) error {
				return dal.UpdateHWFeatures(c, testAppID, []*models.AppHWFeatures{
					{
						HWFeature:               models.HWFeatureRumble,
						HWFeatureControllerType: models.HWFeatureControllerFamilyCommon,
						Notes:                   "changed",
					},
					// same feature of the other family (PlayStation) must be removed
					{
						HWFeature:               models.HWFeatureHDHaptics,
						HWFeatureControllerType: models.HWFeatureControllerFamilySteam,
					},
					{
						HWFeature:               models.HWFeatureTouchpads,
						HWFeatureControllerType: models.HWFeatureControllerFamilySteam,
					},
				}, nil)
			},
			expected: func() *models.AppInfo {
				info := fullInfo(testAppID)
				info.HWFeatures = []*models.AppHWFeatures{
					{
						HWFeature:               models.HWFeatureRumble,
						HWFeatureControllerType: models.HWFeatureControllerFamilyCommon,
						Notes:                   "changed",
					},
					{
						HWFeature:               models.HWFeatureHDHaptics,
						HWFeatureControllerType: models.HWFeatureControllerFamilySteam,
					},
					{
						HWFeature:               models.HWFeatureTouchpads,
						HWFeatureControllerType: models.HWFeatureControllerFamilySteam,
					},
				}
				return info
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := context.Background()
			dal, err := sidbtest.NewMemDB(t, false)
			require.NoError(t, err)

			if !tc.noApp {
				err = dal.AppInfo().Insert(c, &models.AppInfo{
					AppID: testAppID,
				})
				require.NoError(t, err)
				// other app must never be affected
				err = dal.AppInfo().Insert(c, &models.AppInfo{
					AppID: testAppID + 1,
				})
				require.NoError(t, err)
				err = dal.AppInfo().UpdateSteamInputDBInfo(c, fullInfo(testAppID+1))
				require.NoError(t, err)
			}
			if tc.initial != nil {
				err = dal.AppInfo().UpdateSteamInputDBInfo(c, tc.initial)
				require.NoError(t, err)
			}

			err = tc.update(c, dal.AppInfo())
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)

			expectedInfos := []*models.AppInfo{
				tc.expected,
				fullInfo(testAppID + 1),
			}
			for _, expected := range expectedInfos {
				got, err := dal.AppInfo().Get(c, expected.AppID, appinfodal.AppInfoInclude{
					SteamInputDBInfo: true,
				})
				require.NoError(t, err)

				got.MixedInputInfo.AppID = 0
				got.MixedInputInfo.CreatedAt = time.Time{}
				got.MixedInputInfo.UpdatedAt = time.Time{}
				for _, link := range got.MixedInputInfo.MixedInputModLinks {
					link.AppID = 0
				}
				got.SteamInputAPISupport.AppID = 0
				got.SteamInputAPISupport.CreatedAt = time.Time{}
				got.SteamInputAPISupport.UpdatedAt = time.Time{}
				for _, siapiType := range got.SteamInputAPISupport.SIAPITypes {
					siapiType.AppID = 0
				}
				got.Glyphs.AppID = 0
				for _, ctrl := range got.Glyphs.GlyphCtrlSupport {
					ctrl.AppID = 0
				}
				for _, tag := range got.Glyphs.GlyphTags {
					tag.AppID = 0
				}
				for _, feature := range got.HWFeatures {
					feature.AppID = 0
					feature.CreatedAt = time.Time{}
					feature.UpdatedAt = time.Time{}
				}

				assert.Equal(t, expected, &models.AppInfo{
					AppID:                   got.AppID,
					ControllerSupportRating: got.ControllerSupportRating,
					ControllerSupportNotes:  got.ControllerSupportNotes,
					MixedInputInfo:          got.MixedInputInfo,
					SteamInputAPISupport:    got.SteamInputAPISupport,
					Glyphs:                  got.Glyphs,
					HWFeatures:              got.HWFeatures,
				})
			}
		})
	}
}
