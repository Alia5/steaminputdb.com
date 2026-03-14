import type { ControllerGroup } from './config.types';
import { gyroButtonFlags, niceInputMap } from './const';
export const niceInputName = (input?: string) => {
    const [type, key] = input?.split(' ') ?? [];
    if (!key || !type) {
        return input;
    }
    const strippedKey = key.replace(
        /^(?:LEFT_)?(.)(.*)/,
        (_, a, b) => a + b.toLowerCase()
            .replace(
                /_([a-z])/g,
                (__: unknown, c: string) => c.toUpperCase()
            )
    );

    const keyLabel = (niceInputMap as Record<string, string>)[key] ?? strippedKey;


    const buttons = [
        keyLabel,
        (niceInputMap as Record<string, string>)[type] ?? type
    ].join(' ');

    if (type === 'xinput_button') {
        const isABXY = ['A', 'B', 'X', 'Y'].includes(key.toUpperCase());
        return isABXY ? `${keyLabel} Button` : keyLabel;
    }

    return buttons;
};

export const mergeLayerInputs = (
    layerGroups: { group: ControllerGroup; binding: string }[],
    parentGroups: { group: ControllerGroup; binding: string }[]
): { group: ControllerGroup; binding: string }[] => layerGroups.map((entry) => {
    if (entry.binding.includes('modeshift')) {
        return entry;
    }
    const parentActive = parentGroups.find((pg) => !pg.binding.includes('modeshift'));
    if (!parentActive || !parentActive.group.inputs) {
        return entry;
    }
    const layerInputs = entry.group.inputs ?? {};
    const parentInputs = parentActive.group.inputs;
    const hasAllParentKeys = Object.keys(parentInputs).every((k) => k in layerInputs);
    if (hasAllParentKeys) {
        return entry;
    }
    const merged = { ...parentInputs, ...layerInputs };
    return {
        ...entry,
        group: { ...entry.group, inputs: merged }
    };
});

export const decodeGyroButtons = (value: string): string[] => {
    if (!value){
        return ['Always On'];
    }
    const num = BigInt(value);
    const names: string[] = [];
    for (let i = 0; i < 64; i++) {
        if (num & (1n << BigInt(i))) {
            names.push(gyroButtonFlags[i] ?? `Unknown (bit ${i})`);
        }
    }
    return names.length ? names : [`Button ${value}`];
};
