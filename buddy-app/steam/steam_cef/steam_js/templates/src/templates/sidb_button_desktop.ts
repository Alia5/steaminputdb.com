declare global {
    interface Window {
        __sidbOpenButtonObserver?: MutationObserver;
    }
}

import { mount } from 'svelte';
import { ATTR } from '../components/sidb_button_desktop.svelte';
import DesktopButton from '../components/sidb_button_desktop.svelte';

const openUrl = (url: string) => {
    goTmpl(' if .UseSteamBrowser ');
    window.opener.SteamUIStore.ActiveWindowInstance.m_Navigator.SteamWebTab(url);
    goTmpl(' else ');
    open(url);
    goTmpl(' end ');
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
        document.querySelectorAll('.SVGIcon_Settings').forEach((el) => {
            const container = el.parentElement!.parentElement!.parentElement!;
            if (container.querySelector(`[${ATTR}]`)) {
                if (override) {
                    container.querySelector(`[${ATTR}]`)!.remove();
                } else {
                    return;
                }
            }
            const node = container.children[1] as HTMLElement;
            const insertIdx = container.children.length >= 4 ? 2 : 1;
            const anchor = container.children[insertIdx];
            mount(DesktopButton, {
                target: container,
                anchor: anchor,
                props: {
                    class: node.className,
                    tabindex: node.tabIndex,
                    role: node.getAttribute('role'),
                    openUrl
                }
            });

        });
        override = false;
    };

    inject();

    const observer = new MutationObserver(inject);
    observer.observe(document.body, { childList: true, subtree: true });

    window.__sidbOpenButtonObserver = observer;
    window.__sidbCleanup = (window.__sidbCleanup || []).concat(() => {
        observer.disconnect();
        document.querySelectorAll(`[${ATTR}]`).forEach((el) => el.remove());
    });

})();
