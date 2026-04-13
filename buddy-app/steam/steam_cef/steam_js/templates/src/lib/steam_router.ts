export const BrowserInputSupport = {
    PageUnloading: 0,
    Unknown: 1,
    None: 2,
    Basic: 3,
    Full: 4
} as const;

export type BrowserInputSupportLevel = typeof BrowserInputSupport[keyof typeof BrowserInputSupport];

export const controllerButtonToHID: Record<number, number[]> = {
    0: [88],     // A => Enter
    // 1: [41],     // B => Escape
    2: [41],     // X => ESCAPE
    3: [43],     // Y => Tab
    20: [82],    // Up => Up Arrow
    21: [81],    // Down => Down Arrow
    22: [80],    // Left => Left Arrow
    23: [79]     // Right => Right Arrow
};

export interface SteamRouter {
    Navigate: unknown;
    NavigationManager: unknown;
    WindowStore: {
        GamepadUIMainWindowInstance: {
            NavigateToSteamWeb?: (url: string, name: string, newTab: boolean) => void;
            m_StoreBrowser: {
                m_fnGoBackOverride?: () => boolean;
                m_bCanGoBackward: boolean;
                m_browserView: {
                    SetSteamURLCallback: (callback: (url: string) => void) => void;
                    on: (event: string, callback: () => void) => void;
                    off: (event: string, callback: () => void) => void;
                    CanGoBackward: () => boolean;
                    GoBack: () => void;
                    SetFocus: (focus: boolean) => unknown;
                };
                m_gamepadBridge: {
                    SetGameInputSupportLevel: (level: BrowserInputSupportLevel, name: string) => void;
                    TakeFocus: () => void;
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    m_NavigationController: Record<string, any>;
                };
                m_URL: string;
            };
            m_FooterStore: {
                m_Instance: {
                    m_ActionDescriptionStore: {
                        SetActionsFromMap: (map: Record<number, string>) => void;
                    };
                };
            };
        };
    };
}
