import { log } from '$lib/log';
import {
    TEXT_INPUT_TYPES,
    ZERO_RECT,
    findBestInDir, initDialogFocusTrap, initGamepadNavFocusable, scrollToEl, setGamepadNavFocused, shouldTextInputConsume
} from './helper';

const DEFAULT_UPDATE_MS = 16;
const DEFAULT_AXIS_THRESHOLD = 0.5;

const DEV_MOCK_BPM = false;

export const STANDARD_BUTTON_NAMES = [
    'button-south', 'button-east', 'button-west', 'button-north',
    'shoulder-left', 'shoulder-right', 'trigger-left', 'trigger-right',
    'select', 'start', 'thumbstick-left', 'thumbstick-right',
    'dpad-up', 'dpad-down', 'dpad-left', 'dpad-right', 'home'
] as const;

export const STANDARD_AXIS_NAMES = [
    'left-stick-x', 'left-stick-y',
    'right-stick-x', 'right-stick-y'
] as const;

export type NavDirection = 'up' | 'down' | 'left' | 'right' | undefined;

export type GamepadNavEvent = Omit<Gamepad, 'buttons' | 'axes'> & {
    direction: NavDirection;
    source?: Event|KeyboardEvent;
};

export type GamepadInputEvent = Omit<Gamepad, 'buttons' | 'axes'> & {
    buttons?: {
        [index: number]: GamepadButton;
        [name: string]: GamepadButton;
    };
    axes?: {
        [index: number]: number;
        [name: string]: number;
    };
};

class gamepadNavigator {

    public gamepads = $state<Gamepad[]>([]);
    public axisThreshold = $state(DEFAULT_AXIS_THRESHOLD);
    public navigationDirection = $state<NavDirection | undefined>();
    private updateInterval: NodeJS.Timeout | undefined;
    private navFocusableObserver: MutationObserver | undefined;
    private dialogFocusTrapObserver: MutationObserver | undefined;

    private enabled = false;
    private enableKeyboardNav = false;
    private skipGamepads = false;
    private selectOpen = false;
    private lastNavTime = 0;

    public constructor() {

    }

    public Init() {
        if (this.enabled) {
            return;
        }
        log.debug('Initializing GamepadNavigator');
        this.enabled = true;
        this.navFocusableObserver = initGamepadNavFocusable();
        this.dialogFocusTrapObserver = initDialogFocusTrap();
        window.addEventListener('gamepadconnected', this.onGamepadConnected.bind(this));
        window.addEventListener('gamepaddisconnected', this.onGamepadDisconnected.bind(this));
        if (navigator?.userAgent?.includes('Steam Gamepad') || DEV_MOCK_BPM) {
        // if (navigator?.userAgent?.includes(' ')) {
            window.addEventListener('keydown', this.onKeyDown.bind(this));
            this.enableKeyboardNav = true;
            this.skipGamepads = true;
            clearInterval(this.updateInterval);
        } else {
            this.enableKeyboardNav = false;
            this.skipGamepads = false;
        }
    }

    public Dispose() {
        log.debug('Disposing GamepadNavigator');
        this.enabled = false;
        this.navFocusableObserver?.disconnect();
        this.dialogFocusTrapObserver?.disconnect();
        clearInterval(this.updateInterval);
        window.removeEventListener('gamepadconnected', this.onGamepadConnected.bind(this));
        window.removeEventListener('gamepaddisconnected', this.onGamepadDisconnected.bind(this));
        window.removeEventListener('keydown', this.onKeyDown.bind(this));
    }

    public onInput(event: GamepadInputEvent) {
        if (event.buttons?.['button-south']?.pressed || event.buttons?.[0]?.pressed) {
            this.activateFocused();
        }
        if (event.buttons?.['button-east']?.pressed || event.buttons?.[1]?.pressed) {
            const active = document.activeElement as HTMLElement | null;
            if (active && active !== document.body) {
                active.blur();
            } else {
                history.back();
            }
        }
    }


    public onNavEvent(event: GamepadNavEvent) {
        if (!event.direction) {
            return;
        }

        const now = performance.now();
        if (now - this.lastNavTime < 200) {
            return;
        }
        this.lastNavTime = now;

        const active = document.activeElement as HTMLElement | null;
        const current = (active && active !== document.body && active !== document.documentElement)
            ? active : undefined;
        if (shouldTextInputConsume(current, event.direction)) {
            return;
        }

        if (event.id === 'keyboard') {
            event.source?.preventDefault();
        }

        if (event.id !== 'keyboard' && current?.tagName === 'SELECT'
            && this.selectOpen
            && (event.direction === 'up' || event.direction === 'down')) {
            const select = current as HTMLSelectElement;
            const delta = event.direction === 'down' ? 1 : -1;
            select.selectedIndex = Math.max(0, Math.min(select.options.length - 1, select.selectedIndex + delta));
            select.dispatchEvent(new Event('change', { bubbles: true }));
            select.dispatchEvent(new Event('input', { bubbles: true }));
            return;
        }

        const ref = current?.getBoundingClientRect() ?? ZERO_RECT;
        const next = findBestInDir(event.direction, ref, current);
        if (next) {
            next.focus({ preventScroll: true });
            scrollToEl(next);
            setGamepadNavFocused(next);
            if (this.skipGamepads && (
                next.tagName === 'TEXTAREA'
                || next.isContentEditable
                || (next.tagName === 'INPUT'
                    && TEXT_INPUT_TYPES.includes(
                        (next as HTMLInputElement).type as typeof TEXT_INPUT_TYPES[number]
                    ))
            )) {
                open('steam://sidb/PrimeOSK');
            }
        }
    }

    private activateFocused() {
        const active = document.activeElement as HTMLElement | null;
        if (!active || active === document.body || active === document.documentElement) {
            return;
        }
        if (active.tagName === 'TEXTAREA' || active.isContentEditable) {
            return;
        }
        if (active.tagName === 'INPUT'
            && TEXT_INPUT_TYPES.includes(
                (active as HTMLInputElement).type as typeof TEXT_INPUT_TYPES[number]
            )) {
            return;
        }
        if (active.tagName === 'SELECT') {
            this.selectOpen = true;
            active.addEventListener('blur', () => { this.selectOpen = false; }, { once: true });
            (active as HTMLSelectElement).showPicker();
            return;
        }
        active.click();
    }

    private onKeyDown(event: KeyboardEvent) {
        if (!this.enableKeyboardNav) {
            return;
        }
        if (event.key === 'Enter') {
            event.preventDefault();
            this.activateFocused();
            return;
        }
        const direction = event.key?.includes('Arrow')
            ? (event.key.split('Arrow')?.[1]?.toLowerCase() as NavDirection)
            : undefined;
        if (direction) {
            this.onNavEvent({
                id: 'keyboard',
                source: event,
                timestamp: Date.now(),
                connected: true,
                direction
            } as GamepadNavEvent);
        }
    }

    private onGamepadConnected(event: GamepadEvent) {
        if (this.skipGamepads) {
            clearInterval(this.updateInterval);
            return;
        }
        log.debug('Gamepad connected', 'event', event);
        this.updateGamepads();
        clearInterval(this.updateInterval);
        this.updateInterval = setInterval(() => {
            this.updateGamepads();
        }, DEFAULT_UPDATE_MS);
    }

    private onGamepadDisconnected(event: GamepadEvent) {
        log.debug('Gamepad disconnected', 'event', event);
        this.updateGamepads();
        if (!this.gamepads.length) {
            clearInterval(this.updateInterval);
        }
    }

    public onkeypress = this.onKeyDown.bind(this);
    public ongamepadconnected = this.onGamepadConnected.bind(this);
    public ongamepaddisconnected = this.onGamepadDisconnected.bind(this);

    public updateGamepads() {
        if (document.hidden || !document.hasFocus()) {
            return;
        }
        const updated = navigator?.getGamepads()?.filter((g) => g) as Gamepad[]
            ?? [];

        const actualChange = updated.reduce((acc, gp) => {
            const prev = this.gamepads.find((p) => p.index === gp.index);
            const changedButtons = gp.buttons.reduce((btns, btn, i) => {
                if (btn.pressed !== (prev?.buttons[i]?.pressed ?? false)) {
                    btns[i] = btn;
                    if (gp.mapping === 'standard' && STANDARD_BUTTON_NAMES[i]) {
                        btns[STANDARD_BUTTON_NAMES[i]] = btn;
                    }
                }
                return btns;
            }, {} as Record<number | string, GamepadButton>);
            const changedAxes = gp.axes.reduce((axs, axis, i) => {
                if (Math.abs(axis - (prev?.axes[i] ?? 0)) > this.axisThreshold) {
                    axs[i] = axis;
                    if (gp.mapping === 'standard' && STANDARD_AXIS_NAMES[i]) {
                        axs[STANDARD_AXIS_NAMES[i]] = axis;
                    }
                }
                return axs;
            }, {} as Record<number | string, number>);
            const buttonsChanged = !!Object.keys(changedButtons).length;
            const axesChanged = !!Object.keys(changedAxes).length;

            let navDirection = Object.keys(changedButtons)
                .filter((b) => b.startsWith('dpad-') && changedButtons[b]?.pressed)?.[0]
                ?.split('-')?.[1] as NavDirection | undefined;

            if (buttonsChanged || axesChanged) {
                this.onInput({
                    ...gp,
                    ...(buttonsChanged ? { buttons: changedButtons } : undefined),
                    ...(axesChanged ? { axes: changedAxes } : undefined)
                } as GamepadInputEvent);
                acc = true;
            }
            if (!navDirection) {
                const lsX = gp.axes[0] ?? 0;
                const lsY = gp.axes[1] ?? 0;

                if (Math.abs(lsY) > this.axisThreshold) {
                    navDirection = lsY < 0 ? 'up' : 'down';
                } else if (Math.abs(lsX) > this.axisThreshold) {
                    navDirection = lsX < 0 ? 'left' : 'right';
                }
            }
            if (navDirection !== this.navigationDirection) {
                this.navigationDirection = navDirection;
                if (navDirection) {
                    this.onNavEvent({ ...gp, direction: navDirection } as GamepadNavEvent);
                }
            }
            return acc;
        }, false);
        if (actualChange) {
            this.gamepads = updated;
        }
    }

}

export const GamepadNavigator = new gamepadNavigator();

