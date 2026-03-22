import { clientWithSvelteFetch } from '$lib/buddy-api/client';
import { log } from '$lib/log';
import { error } from '@sveltejs/kit';


export const pingBuddy = async (fetch: typeof globalThis.fetch = globalThis.fetch) => {
    const client = clientWithSvelteFetch(fetch);
    const resp = await client.GET('/v1/ping');
    if (resp.error) {
        error(resp?.response?.status ?? 502, {
            message: 'Buddy APP request failed',
            err: resp.error
        });
    }

    log.debug('Ping buddy response', 'status', resp.response.status, 'data', resp.data);
    return resp.data;
};

export const steamClientStatus = async (fetch: typeof globalThis.fetch = globalThis.fetch) => {
    const client = clientWithSvelteFetch(fetch);
    const resp = await client.GET('/v1/steam/status');
    if (resp.error) {
        error(resp?.response?.status ?? 502, {
            message: 'Buddy APP request failed',
            err: resp.error
        });
    }
    log.debug('Steam client status response', 'status', resp.response.status, 'data', resp.data);
    return resp.data;
};

export const getLatestVersion = async (fetch: typeof globalThis.fetch = globalThis.fetch) => {
    const resp = await fetch('https://api.github.com/repos/Alia5/steaminputdb.com/releases/latest');
    if (!resp.ok) {
        log.error('Failed to fetch latest release from GitHub', 'status', resp.status, 'resp', await resp.json());
        error(resp.status, {
            message: 'Failed to fetch latest release from GitHub',
            err: resp.statusText
        });
    }

    return (await resp.json()).tag_name;
};

export const getBuddySettings = async (fetch: typeof globalThis.fetch = globalThis.fetch) => {
    const client = clientWithSvelteFetch(fetch);
    const resp = await client.GET('/v1/settings');
    if (resp.error) {
        error(resp?.response?.status ?? 502, {
            message: 'Buddy APP request failed',
            err: resp.error
        });
    }
    log.debug('Buddy settings response', 'status', resp.response.status, 'data', resp.data);
    return resp.data;
};
