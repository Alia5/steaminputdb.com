

import { mount } from 'svelte';
import { ATTR } from '../components/sidb_button_bpm.svelte';
import BPMButton from '../components/sidb_button_bpm.svelte';
import { findWebpackModule, findWebpackModuleExport } from '$lib/webpack';
import { patchReactFiber } from '$lib/react_crap';
import { BrowserInputSupport, type BrowserInputSupportLevel, controllerButtonToHID, type SteamRouter } from '$lib/steam_router';

declare global {
    interface Window {
        __sidbOpenButtonObserver?: MutationObserver;
        __sidbCleanup?: (() => void)[];
        Router?: SteamRouter;
    }
}


interface React {
    createElement: (type: unknown, props: unknown, ...children: unknown[]) => HTMLElement;
}

let injectRunning = false;

let oskPrimed = false;

const setGameInputSupportLevel = (raw: string, router: SteamRouter) => {
    const level = Math.min(
        Math.max(Number.parseInt(raw, 10), BrowserInputSupport.PageUnloading),
        BrowserInputSupport.Full
    ) as BrowserInputSupportLevel;
    router.WindowStore.GamepadUIMainWindowInstance
        .m_StoreBrowser.m_gamepadBridge
        .SetGameInputSupportLevel(level, 'sidb');
};

const setFooterPrompts = (encoded: string, router: SteamRouter) => {
    const map = JSON.parse(decodeURIComponent(encoded));
    router.WindowStore.GamepadUIMainWindowInstance
        .m_FooterStore.m_Instance
        .m_ActionDescriptionStore.SetActionsFromMap(map);
};

const sidbActions: Record<string, (payload: string, router: SteamRouter) => void> = {
    SetGameInputSupportLevel: setGameInputSupportLevel,
    SetFooterPrompts: setFooterPrompts,
    PrimeOSK: () => { oskPrimed = true; }
};

const handleSteamURLCallback = (
    router: SteamRouter
) => (cb_url: string) => {
    const sidbMatch = cb_url.match(/^steam:\/\/sidb\/(\w+)(?:\/(.+))?$/);
    if (sidbMatch) {
        const [, action, payload] = sidbMatch;
        sidbActions[action!]?.(payload!, router);
        return;
    }
    if (cb_url.includes('https://store.steampowered.com')) {
        router.WindowStore.GamepadUIMainWindowInstance.NavigateToSteamWeb!(
            'http' + cb_url.split('http').pop(), 'sidb', true
        );
    }
};

const registerControllerForwarding = (
    Router: SteamRouter
) => {
    const unreg: { unregister: () => void } | null = SteamClient.Input.RegisterForControllerInputMessages!(
        (_idx: number, button: number, pressed: boolean) => {
            const url = Router.WindowStore.GamepadUIMainWindowInstance.m_StoreBrowser
                .m_URL;
            if (!url?.includes('steaminput')
                && !url?.includes('localhost:5173')
            ) {
                return;
            }

            if ((
                Router.WindowStore.GamepadUIMainWindowInstance
                    .m_StoreBrowser.m_gamepadBridge.m_NavigationController
                    ?.m_ActiveContext?.m_activeBrowserView)
                    !== 'MainBrowser'
            ) {
                return;
            }

            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            if ((Router.WindowStore.GamepadUIMainWindowInstance as any)
                .m_VirtualKeyboardManager?.IsShowingVirtualKeyboard?.m_currentValue) {
                return;
            }

            if (pressed && button !== 0) {
                oskPrimed = false;
            }
            if (button === 0 && pressed && oskPrimed) {
                oskPrimed = false;
                try {
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    (Router.WindowStore.GamepadUIMainWindowInstance.m_StoreBrowser as any)
                        .m_refKeyboard?.ShowVirtualKeyboard();
                } catch (e) {
                    console.error('Failed to open OSK', e);
                }
                return;
            }

            const hids = controllerButtonToHID[button];
            if (hids !== undefined) {
                if (hids.length > 1) {
                    if (pressed) {
                        hids.forEach((hid) => {
                            window.opener.SteamClient.Input.ControllerKeyboardSetKeyState!(hid, true);
                            window.opener.SteamClient.Input.ControllerKeyboardSetKeyState!(hid, false);
                        });
                    }
                } else {
                    window.opener.SteamClient.Input.ControllerKeyboardSetKeyState!(hids[0]!, pressed);
                }
            }
            if (button === 1 && pressed) {
                if (Router.WindowStore.GamepadUIMainWindowInstance.m_StoreBrowser.m_fnGoBackOverride) {
                    if (Router.WindowStore.GamepadUIMainWindowInstance.m_StoreBrowser.m_fnGoBackOverride()) {
                        setTimeout(() => {
                            if (!Router.WindowStore.GamepadUIMainWindowInstance.m_StoreBrowser.m_fnGoBackOverride && unreg) {
                                unreg.unregister();
                                window.__sidbGamepadCallbackClean = undefined;
                            }
                        }, 1000);
                    }
                } else if (Router.WindowStore.GamepadUIMainWindowInstance.m_StoreBrowser.m_bCanGoBackward) {
                    Router.WindowStore.GamepadUIMainWindowInstance.m_StoreBrowser.m_browserView.GoBack();
                }
            }
        }
    ) as { unregister: () => void };
    if (window.__sidbGamepadCallbackClean) {
        window.__sidbGamepadCallbackClean.unregister?.();
        window.__sidbGamepadCallbackClean = undefined;
    }
    window.__sidbGamepadCallbackClean = unreg;
    window.__sidbCleanup?.push(() => {
        if (window.__sidbGamepadCallbackClean) {
            window.__sidbGamepadCallbackClean.unregister?.();
            window.__sidbGamepadCallbackClean = undefined;
        }
    });
};


__INJECT_RETURN = (() => {
    let override = goTmpl(' .Override ');
    if (window.__sidbOpenButtonObserver) {
        if (override) {
            window.__sidbOpenButtonObserver.disconnect();
        } else {
            return;
        }
    }

    const inject = () => {
        if (injectRunning) {
            return;
        }
        injectRunning = true;
        try {
            const React = findWebpackModule((m) => m.createElement && m.Fragment && m.Component) as React;
            if (!React) {
                console.error('React not found');
                injectRunning = false;
                return;
            }
            const Focusable = findWebpackModuleExport((e) => {
                const str = (typeof e === 'function' ? e?.toString?.() : e?.render?.toString?.()) ?? '';
                return /flow-children/.test(str) && /onActivate/.test(str) && /focusClassName/.test(str);
            });
            if (!Focusable) {
                console.error('Focusable not found');
                injectRunning = false;

                return;
            }
            const Router = findWebpackModuleExport((e) => e.Navigate && e.NavigationManager) as SteamRouter;
            if (!Router) {
                console.error('Router not found');
                injectRunning = false;

                return;
            }
            window.Router = Router;

            const steamWebOpener = () => {
                const mainWindow = Router.WindowStore.GamepadUIMainWindowInstance;
                if (!mainWindow.NavigateToSteamWeb) {
                    return undefined;
                }


                return (url: string) => {
                    try {
                        mainWindow.NavigateToSteamWeb!(url, 'sidb', true);
                        setGameInputSupportLevel(String(BrowserInputSupport.Full), Router);
                        Router.WindowStore.GamepadUIMainWindowInstance.m_StoreBrowser.m_browserView.on('finished-request', () => {
                            setTimeout(() =>
                            {

                                window.opener.SteamClient.Window.SetKeyFocus(false);
                                Router.WindowStore.GamepadUIMainWindowInstance
                                    .m_StoreBrowser.m_browserView.SetFocus(true);
                                window.opener.SteamClient.Input.SetWebBrowserActionset!(true);

                                setGameInputSupportLevel(String(BrowserInputSupport.Full), Router);
                                Router.WindowStore.GamepadUIMainWindowInstance
                                    .m_FooterStore.m_Instance
                                    .m_ActionDescriptionStore.SetActionsFromMap({ 0: 'Select', 1: 'Back' });

                            }, 1000);
                        });
                        registerControllerForwarding(Router);
                        Router.WindowStore.GamepadUIMainWindowInstance
                            .m_StoreBrowser.m_browserView
                            .SetSteamURLCallback(handleSteamURLCallback(Router));
                    } catch {
                        (window.opener.open ?? window.open)(url);
                    }
                };
            };

            const containers = [
                ...(
                    [...document.querySelectorAll('.SVGIcon_BigPicture')]
                        .map((el) => el.closest('[role="button"]')?.parentElement?.parentElement)
                ??
                [...document.querySelectorAll('path[d^="M33 20.38V15.62"]')]
                    .map((el) => el.closest('[role="button"]')?.parentElement?.parentElement)
                )
            ].filter(Boolean) as HTMLElement[];
            const container = [...containers].pop();
            if (!container) {
                injectRunning = false;
                return;
            }
            if (container.querySelector(`[${ATTR}]`)) {
                if (override) {
                    container.querySelector(`[${ATTR}]`)!.remove();
                } else {
                    injectRunning = false;
                    return;
                }
            }

            const existingBtn = container.querySelector('[role="button"]');
            const btnClassName = existingBtn?.className ?? '';

            const openFn = steamWebOpener()
                ?? window.opener.open
                ?? window.open;
            const openSidb = () => {
                try {
                    const appIdPath = window.opener.SteamUIStore.ActiveWindowInstance.m_locationPathname;
                    const appId = Number.parseInt(
                        appIdPath.split('/').find((p: string) => Number.isInteger(Number.parseInt(p, 10))) ?? '0',
                        10
                    );
                    const displayName = (
                        window.opener.appDetailsStore.m_mapAppData?.get(appId)?.details?.strDisplayName
                    );
                    const isNonSteam = (
                        !!window.opener.appDetailsStore.m_mapAppData?.get(appId)?.details?.strShortcutExe
                    );
                    if (isNonSteam) {
                        openFn(`https://steaminputdb.com/app/${encodeURIComponent(displayName)}?buddy-app=enabled`);
                        return;
                    }
                    openFn(`https://steaminputdb.com/app/${appId}?buddy-app=enabled`);
                } catch (e) {
                    console.error(e);
                    openFn('https://steaminputdb.com/config/search?buddy-app=enabled');
                }
            };

            const svelteMountSuffix = Math.random().toString(36).substring(2, 15);
            const panel = React.createElement(
                'div',
                {
                    className: 'Panel',
                    [ATTR]: '1',
                    key: '__sidb'
                },
                React.createElement(
                    Focusable,
                    {
                        'className': btnClassName,
                        'onClick': openSidb,
                        'onActivate': openSidb,
                        'aria-label': 'SteamInputDB'
                    },
                    React.createElement('div', {
                        className: `__sidb-svelte-mount-${svelteMountSuffix}`
                    })
                )
            );
            patchReactFiber(container, panel, 1);
            const mountObserver = new MutationObserver(() => {
                const mountPoint = container.querySelector(`.__sidb-svelte-mount-${svelteMountSuffix}`);
                if (mountPoint && !mountPoint.hasChildNodes()) {
                    mountObserver.disconnect();
                    mount(BPMButton, { target: mountPoint });
                }
            });
            mountObserver.observe(container, { childList: true, subtree: true });
        } finally {
            setTimeout(() => {
                injectRunning = false;
            }, 200);
        }

        override = false;

    };

    inject();

    const observer = new MutationObserver(() => {
        inject();
    });
    observer.observe(document.body, { childList: true, subtree: true });

    if (window.__sidbOpenButtonObserver) {
        window.__sidbOpenButtonObserver.disconnect();
    }
    window.__sidbOpenButtonObserver = observer;
    window.__sidbCleanup = (window.__sidbCleanup || []).concat(() => {
        observer?.disconnect();
        document.querySelectorAll(`[${ATTR}]`)?.forEach((el) => el.remove());
    });

})();
