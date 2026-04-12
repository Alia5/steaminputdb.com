

import { mount } from 'svelte';
import { ATTR } from '../components/sidb_button_bpm.svelte';
import BPMButton from '../components/sidb_button_bpm.svelte';
import { findWebpackModule, findWebpackModuleExport } from '$lib/webpack';
import { patchReactFiber } from '$lib/react_crap';

declare global {
    interface Window {
        __sidbOpenButtonObserver?: MutationObserver;
        __sidbCleanup?: (() => void)[];
        Router?: SteamRouter;
    }
}
interface SteamRouter {
    Navigate: unknown;
    NavigationManager: unknown;
    WindowStore: {
        GamepadUIMainWindowInstance: {
            NavigateToSteamWeb: (url: string, name: string, newTab: boolean) => void;
            m_StoreBrowser: {
                m_browserView: {
                    SetSteamURLCallback: (callback: (url: string) => void) => void;
                };
            };
        };
    };
}


interface React {
    createElement: (type: unknown, props: unknown, ...children: unknown[]) => HTMLElement;
}

let injectRunning = false;


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
        const React = findWebpackModule((m) => m.createElement && m.Fragment && m.Component) as React;
        if (!React) {
            console.error('React not found');
            return;
        }
        const Focusable = findWebpackModuleExport((e) => {
            const str = (typeof e === 'function' ? e?.toString?.() : e?.render?.toString?.()) ?? '';
            return /flow-children/.test(str) && /onActivate/.test(str) && /focusClassName/.test(str);
        });
        if (!Focusable) {
            console.error('Focusable not found');
            return;
        }
        const Router = findWebpackModuleExport((e) => e.Navigate && e.NavigationManager);
        if (!Router) {
            console.error('Router not found');
            return;
        }
        window.Router = Router as SteamRouter;


        const containers = [
            ...(
                [...document.querySelectorAll('.SVGIcon_BigPicture')]
                    .map((el) => el.closest('[role="button"]')?.parentElement?.parentElement)
                ??
                [...document.querySelectorAll('path[d^="M33 20.38V15.62"]')]
                    .map((el) => el.closest('[role="button"]')?.parentElement?.parentElement)
            )
        ].filter(Boolean) as HTMLElement[];
        containers.forEach((container) => {
            if (container.querySelector(`[${ATTR}]`)) {
                if (override) {
                    container.querySelector(`[${ATTR}]`)!.remove();
                } else {
                    return;
                }
            }

            const existingBtn = container.querySelector('[role="button"]');
            const btnClassName = existingBtn?.className ?? '';

            const openFn = window.Router?.WindowStore?.GamepadUIMainWindowInstance?.NavigateToSteamWeb
                ? ((url: string) => {
                    try {
                        window.Router!.WindowStore.GamepadUIMainWindowInstance.NavigateToSteamWeb(url, 'sidb', true);
                        window.Router!.WindowStore.GamepadUIMainWindowInstance.m_StoreBrowser.m_browserView.SetSteamURLCallback(
                            (cb_url) => {
                                if (cb_url.includes('https://store.steampowered.com')) {
                                    window.Router!.WindowStore.GamepadUIMainWindowInstance.NavigateToSteamWeb(
                                        'http' + cb_url.split('http').pop(),
                                        'sidb',
                                        true
                                    );
                                }
                            }
                        );
                    } catch {
                        (window.opener.open ?? window.open)(url);
                    }
                })
                : window.opener.open
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
        });

        override = false;
        injectRunning = false;
    };

    inject();

    const observer = new MutationObserver(() => {
        if (injectRunning) {
            return;
        }
        inject();
    });
    observer.observe(document.body, { childList: true, subtree: true });

    window.__sidbOpenButtonObserver = observer;
    window.__sidbCleanup = (window.__sidbCleanup || []).concat(() => {
        observer.disconnect();
        document.querySelectorAll(`[${ATTR}]`).forEach((el) => el.remove());
    });

})();
