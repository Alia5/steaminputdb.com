
import { log } from '$lib/log';
import createClient, { type FetchResponse } from 'openapi-fetch';
import type { paths } from './openapi';

const BUDDY_APP_TIMEOUT_MS = 2000;

const apiURL = 'http://localhost:5119';

const timeoutFetch = (baseFetch: typeof fetch, timeout = BUDDY_APP_TIMEOUT_MS): typeof fetch => (input, init) =>
    baseFetch(input, {
        ...init,
        signal: AbortSignal.timeout(timeout)
    });

log.debug('Creating API client with', 'url', apiURL);
export const client = createClient<paths>({
    baseUrl: apiURL,
    fetch: timeoutFetch(globalThis.fetch.bind(globalThis))
});


export const clientWithSvelteFetch = (fetch: typeof window.fetch, url?: string, timeout = BUDDY_APP_TIMEOUT_MS) => createClient<paths>({
    baseUrl: url || apiURL,
    fetch: timeoutFetch(fetch, timeout)
});

export type ResponseType<
    M extends keyof typeof client,
    P extends keyof paths,
> = P extends keyof paths
    ? Lowercase<M> extends keyof paths[P]
        ? paths[P][Lowercase<M>] extends Record<string | number, unknown>
            ? FetchResponse<paths[P][Lowercase<M>], unknown, `${string}/${string}`>
            : never
        : never
    : never;

