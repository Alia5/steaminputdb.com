
export const ControllerSupportRating = {
    None: 0,
    Wood: 1,
    Stone: 2,
    Bronze: 3,
    Silver: 4,
    Gold: 5,
    Platinum: 6
} as const;
export type ControllerSupportRatingValue = (typeof ControllerSupportRating)[keyof typeof ControllerSupportRating];

export const MixedInputSupportType = {
    Unknown: 0,
    Unsupported: 1,
    Supported: 2,
    WithMod: 3,
    Other: 4
} as const;
export type MixedInputSupportTypeValue = (typeof MixedInputSupportType)[keyof typeof MixedInputSupportType];

export const AppGlyphTagType = {
    Unknown: 0,
    SteamInputAPIButtons: 1,
    InGameButtons: 2,
    SteamInputAPITexts: 3
} as const;
export type AppGlyphTagTypeValue = (typeof AppGlyphTagType)[keyof typeof AppGlyphTagType];

export const SteamInputCameraSupport = {
    Unknown: 0,
    None: 1,
    Partial: 2,
    Full: 3
} as const;
export type SteamInputCameraSupportValue = (typeof SteamInputCameraSupport)[keyof typeof SteamInputCameraSupport];

export const SteamInputAPISupportType = {
    Unknown: 0,
    None: 1,
    InGameActions: 2,
    ActionSets: 3,
    XInputActions: 4,
    Discontinued: 5
} as const;
export type SteamInputAPISupportTypeValue = (typeof SteamInputAPISupportType)[keyof typeof SteamInputAPISupportType];

export const HWFeature = {
    MotionInputs: 0,
    NativeGyroCamera: 1,
    Rumble: 2,
    HDHaptics: 3,
    Lightbar: 4,
    AdaptiveTriggersWired: 5,
    ImpulseTriggers: 6,
    Touchpads: 7,
    AudioHaptics: 8
} as const;
export type HWFeatureValue = (typeof HWFeature)[keyof typeof HWFeature];

export const HWFeatureControllerType = {
    Common: 0,
    Steam: 1,
    PlayStation: 2,
    Xbox: 3,
    Nintendo: 4
} as const;
export type HWFeatureControllerTypeValue = (typeof HWFeatureControllerType)[keyof typeof HWFeatureControllerType];
