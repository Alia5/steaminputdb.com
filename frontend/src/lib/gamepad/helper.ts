import type { NavDirection } from './gamepadNavigator.svelte';

const FOCUSABLE_SELECTOR = 'a[href], button, input, select, textarea, [tabindex]:not([tabindex="-1"]), [data-gamepadnav-focusable]';
export const TEXT_INPUT_TYPES = ['text', 'password', 'search', 'email', 'url', 'tel', 'number'] as const;
const SCROLL_MARGIN = '64px';
const GAMEPAD_NAV_FOCUSABLE = 'data-gamepadnav-focusable';
const GAMEPAD_NAV_FOCUSED = 'data-gamepadnav-focused';

export interface Rect {
    left: number;
    top: number;
    right: number;
    bottom: number;
    width: number;
    height: number;
}
export const ZERO_RECT: Rect = { left: 0, top: 0, right: 0, bottom: 0, width: 0, height: 0 };

export const scrollToEl = (el: HTMLElement) => {
    const prev = el.style.scrollMargin;
    const headerHeight = document.querySelector('header')?.getBoundingClientRect().height ?? 0;
    el.style.scrollMargin = `${headerHeight + 16}px ${SCROLL_MARGIN} ${SCROLL_MARGIN} ${SCROLL_MARGIN}`;
    el.scrollIntoView({ block: 'nearest', inline: 'nearest', behavior: 'smooth' });
    el.style.scrollMargin = prev;
};

const ensureTabIndex = (el: HTMLElement) => {
    if (!el.hasAttribute('tabindex')) {
        el.setAttribute('tabindex', '-1');
    }
};

export const initGamepadNavFocusable = () => {
    document.querySelectorAll<HTMLElement>(`[${GAMEPAD_NAV_FOCUSABLE}]`)
        .forEach(ensureTabIndex);

    const observer = new MutationObserver((mutations) =>
        mutations
            .flatMap((m) => Array.from(m.addedNodes) as HTMLElement[])
            .filter((n) => n.nodeType === Node.ELEMENT_NODE)
            .forEach((n) => {
                if (n.hasAttribute?.(GAMEPAD_NAV_FOCUSABLE)) {
                    ensureTabIndex(n);
                }
                n.querySelectorAll?.<HTMLElement>(`[${GAMEPAD_NAV_FOCUSABLE}]`)
                    .forEach(ensureTabIndex);
            })
    );
    observer.observe(document.body, { childList: true, subtree: true });
    return observer;
};

export const setGamepadNavFocused = (el: HTMLElement) => {
    document.querySelectorAll<HTMLElement>(`[${GAMEPAD_NAV_FOCUSED}]`)
        .forEach((prev) => prev.removeAttribute(GAMEPAD_NAV_FOCUSED));
    if (!el.hasAttribute(GAMEPAD_NAV_FOCUSABLE)) {
        return;
    }
    el.setAttribute(GAMEPAD_NAV_FOCUSED, '');
    el.addEventListener('blur', () => el.removeAttribute(GAMEPAD_NAV_FOCUSED), { once: true });
};

const rectCenter = (r: Rect) => ({
    x: r.left + r.width / 2,
    y: r.top + r.height / 2
});

const isInDir = (ref: Rect, r: Rect, dir: NonNullable<NavDirection>): boolean => {
    const { x: rx, y: ry } = rectCenter(ref);
    const { x: cx, y: cy } = rectCenter(r);
    return ({ up: cy < ry, down: cy > ry, left: cx < rx, right: cx > rx })[dir];
};

const primaryEdgeDist = (ref: Rect, r: Rect, dir: NonNullable<NavDirection>): number =>
    ({ right: r.left - ref.right, left: ref.left - r.right, down: r.top - ref.bottom, up: ref.top - r.bottom })[dir];

const crossOverlap = (ref: Rect, r: Rect, dir: NonNullable<NavDirection>): number => {
    const vert = dir === 'up' || dir === 'down';
    const [rMin, rMax, cMin, cMax] = vert
        ? [ref.left, ref.right, r.left, r.right]
        : [ref.top, ref.bottom, r.top, r.bottom];
    return Math.max(0, Math.min(rMax, cMax) - Math.max(rMin, cMin));
};

const crossDist = (ref: Rect, r: Rect, dir: NonNullable<NavDirection>): number => {
    const vert = dir === 'up' || dir === 'down';
    return vert
        ? Math.abs(rectCenter(r).x - rectCenter(ref).x)
        : Math.abs(rectCenter(r).y - rectCenter(ref).y);
};

const isVisible = (el: HTMLElement): boolean =>
    !!(el.offsetParent || el.offsetWidth || el.offsetHeight);

const isFocusable = (el: HTMLElement): boolean =>
    el.matches(FOCUSABLE_SELECTOR) && isVisible(el);

const hasFocusableDescendant = (el: HTMLElement): boolean =>
    !!el.querySelector(FOCUSABLE_SELECTOR);

const entryEdgeChild = (container: HTMLElement, dir: NonNullable<NavDirection>): HTMLElement | undefined =>
    Array.from(container.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR))
        .filter((el) => isVisible(el))
        .reduce<HTMLElement | undefined>((best, el) => {
            if (!best) {
                return el;
            }
            const r = el.getBoundingClientRect();
            const b = best.getBoundingClientRect();
            return ({
                right: r.left < b.left,
                left: r.right > b.right,
                down: r.top < b.top,
                up: r.bottom > b.bottom
            })[dir] ? el : best;
        }, undefined);

const nearestInDir = (
    entries: { el: HTMLElement; rect: Rect }[],
    dir: NonNullable<NavDirection>,
    ref: Rect
): { el: HTMLElement; rect: Rect } | undefined =>
    entries
        .filter((e) => isInDir(ref, e.rect, dir))
        .sort((a, b) => {
            const aDist = primaryEdgeDist(ref, a.rect, dir);
            const bDist = primaryEdgeDist(ref, b.rect, dir);
            if (aDist !== bDist) {
                return aDist - bDist;
            }
            const aOvl = crossOverlap(ref, a.rect, dir);
            const bOvl = crossOverlap(ref, b.rect, dir);
            if (aOvl !== bOvl) {
                return bOvl - aOvl;
            }
            return crossDist(ref, a.rect, dir) - crossDist(ref, b.rect, dir);
        })[0];

const ancestorDirectChild = (ancestor: HTMLElement, descendant: HTMLElement): HTMLElement | undefined => {
    let node: HTMLElement | undefined = descendant;
    while (node && node.parentElement !== ancestor) {
        node = node.parentElement ?? undefined;
    }
    return node;
};

const isInViewport = (r: Rect): boolean =>
    r.bottom > 0 && r.top < window.innerHeight && r.right > 0 && r.left < window.innerWidth;

const getActiveDialog = (): HTMLElement | undefined =>
    Array.from(document.querySelectorAll<HTMLElement>('dialog[open]:not([data-no-focus-trap])'))
        .filter((d) => d.getBoundingClientRect().width > 0)
        .pop();

const firstFocusable = (root?: HTMLElement): HTMLElement | undefined => {
    const container = root ?? document.documentElement;
    const isDialog = root?.tagName === 'DIALOG';
    const all = Array.from(container.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR))
        .filter((el) => isVisible(el))
        .filter((el) => !isDialog || el.parentElement !== root);
    const visible = all
        .filter((el) => isInViewport(el.getBoundingClientRect()))
        .sort((a, b) => {
            const ar = a.getBoundingClientRect();
            const br = b.getBoundingClientRect();
            return ar.top !== br.top ? ar.top - br.top : ar.left - br.left;
        });
    return visible[0] ?? all.sort((a, b) => {
        const ar = a.getBoundingClientRect();
        const br = b.getBoundingClientRect();
        return ar.top !== br.top ? ar.top - br.top : ar.left - br.left;
    })[0];
};

export const initDialogFocusTrap = (): MutationObserver => {
    const observer = new MutationObserver((mutations) =>
        mutations
            .flatMap((m) => Array.from(m.addedNodes) as HTMLElement[])
            .filter((n) => n.nodeType === Node.ELEMENT_NODE)
            .forEach((n) => {
                const dialog = n.matches?.('dialog[open]:not([data-no-focus-trap])')
                    ? n
                    : n.querySelector?.('dialog[open]:not([data-no-focus-trap])');
                if (!dialog) {
                    return;
                }
                requestAnimationFrame(() => {
                    const target = firstFocusable(dialog as HTMLElement);
                    if (target) {
                        target.focus({ preventScroll: true });
                        scrollToEl(target);
                        setGamepadNavFocused(target);
                    }
                });
            })
    );
    observer.observe(document.body, { childList: true, subtree: true });
    return observer;
};

export const findBestInDir = (
    dir: NonNullable<NavDirection>,
    ref: Rect,
    current?: HTMLElement
): HTMLElement | undefined => {
    const dialog = getActiveDialog();

    if (!current || (dialog && !dialog.contains(current))) {
        return firstFocusable(dialog);
    }

    let scope = current.parentElement ?? document.body;
    const boundary = dialog ?? document.body;

    while (scope) {
        const skipChild = ancestorDirectChild(scope, current);

        const entries = (Array.from(scope.children) as HTMLElement[])
            .filter((child) => child !== skipChild && child !== current)
            .filter((child) => isFocusable(child) || hasFocusableDescendant(child))
            .map((child) => ({ el: child, rect: child.getBoundingClientRect() }));

        const best = nearestInDir(entries, dir, ref);
        if (best) {
            return isFocusable(best.el) ? best.el : entryEdgeChild(best.el, dir);
        }

        if (scope === boundary) {
            break;
        }
        scope = scope.parentElement ?? document.body;
    }
    return undefined;
};

export const shouldTextInputConsume = (el: HTMLElement | undefined, dir: NonNullable<NavDirection>): boolean => {
    if (!el || dir === 'up' || dir === 'down') {
        return false;
    }
    if (
        !(el.tagName === 'INPUT' && TEXT_INPUT_TYPES.includes((el as HTMLInputElement).type as typeof TEXT_INPUT_TYPES[number]))
        && el.tagName !== 'TEXTAREA'
        && !el.isContentEditable
    ) {
        return false;
    }
    const { selectionStart, selectionEnd, value } = el as HTMLInputElement;
    if (selectionStart !== selectionEnd) {
        return true;
    }
    if (dir === 'left') {
        return (selectionStart ?? 0) > 0;
    }
    return (selectionStart ?? 0) < (value?.length ?? 0);
};
