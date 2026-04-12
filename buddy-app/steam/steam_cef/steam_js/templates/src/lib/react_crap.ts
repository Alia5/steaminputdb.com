export interface ReactFiber {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    type: any;
    return: ReactFiber | null;
    memoizedState?: {
        queue?: {
            dispatch?: (action: unknown) => void;
        };
    };
    alternate?: ReactFiber;
}

export const getReactFiber = (
    el: HTMLElement
// eslint-disable-next-line @typescript-eslint/no-explicit-any
) => (el as any)?.[Object.keys(el).find((k) => k.startsWith('__reactFiber')) || ''] as ReactFiber | undefined;

export const rerender = (fiber: ReactFiber) => {
    let sf = fiber?.return;
    while (sf && !sf.memoizedState?.queue?.dispatch) {
        sf = sf.return;
    }
    sf?.memoizedState?.queue?.dispatch?.((x: unknown) => x);
};

export const patchReactFiber = (parentEl: HTMLElement, childElement: HTMLElement, insertIdx: number) => {
    let fiber = getReactFiber(parentEl) as ReactFiber | null;
    while (fiber && typeof fiber.type !== 'function') {
        fiber = fiber.return;
    }
    if (!fiber) {
        return null;
    }
    const orig = fiber.type;
    fiber.type = function (...args: unknown[]) {
        const ret = orig.apply(this, args);
        try {
            const children = Array.isArray(ret?.props?.children)
                ? [...ret.props.children]
                : [ret?.props?.children];
            children.splice(insertIdx, 0, childElement);
            return { ...ret, props: { ...ret.props, children } };
        } catch {
            return ret;
        }
    };
    Object.assign(fiber.type, orig);
    fiber.type.toString = () => orig.toString();
    if (fiber.alternate) {
        fiber.alternate.type = fiber.type;
    }
    rerender(fiber);
    return () => {
        fiber.type = orig;
        if (fiber.alternate) {
            fiber.alternate.type = orig;
        }
        rerender(fiber);
    };
};
