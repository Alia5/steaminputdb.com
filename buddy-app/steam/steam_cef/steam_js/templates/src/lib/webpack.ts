// eslint-disable-next-line @typescript-eslint/no-explicit-any
const wpChunk = (window as any).webpackChunksteamui ?? window.opener?.webpackChunksteamui;
if (!wpChunk) {
    console.error('webpackChunksteamui not found');
}


const modules = new Map();
// eslint-disable-next-line @typescript-eslint/no-explicit-any
let wpRequire: any;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
wpChunk.push([[Symbol('@sidb')], {}, (r: any) => { wpRequire = r; }]);
Object.keys(wpRequire.m).forEach((mid) => {
    try {
        const m = wpRequire(mid);
        if (m) {
            modules.set(mid, m);
        }
    } catch  {
        // ignore
    }
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const findWebpackModule = <T>(filter: (m: any) => boolean) => {
    let found: T | undefined;
    modules.forEach((m) => {
        if (found) {
            return;
        }
        if (m.default && filter(m.default)) {
            found = m.default;
        } else if (filter(m)) {
            found = m;
        }
    });
    return found;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const findWebpackModuleExport = <T>(filter: (m: any) => boolean) => {
    let found: T | undefined;
    modules.forEach((m) => {
        if (found) {
            return;
        }
        [m.default, m].forEach((mod) => {
            if (found || typeof mod !== 'object' || !mod) {
                return;
            }
            Object.keys(mod).forEach((key) => {
                if (found) {
                    return;
                }
                try {
                    if (mod[key] && filter(mod[key])) {
                        found = mod[key];
                    }
                } catch  {
                    // ignore
                }
            });
        });
    });
    return found;
};
