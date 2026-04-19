import { PUBLIC_API_BASE_URL_LOCAL } from '$env/static/public';
import { clientWithSvelteFetch, type ResponseType } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import { fetchConfigs } from '$lib/api/searchConfigs';
import { log } from '$lib/log';
import { error, isHttpError } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';


export const load: PageServerLoad = async ({ params, fetch, url }) => {

    const appid = params.appid;

    const app_id = parseInt(appid, 10);
    const app_id_valid = !isNaN(app_id) && app_id > 0;


    const loadRes: {
        appInfo?: components['schemas']['AppInfoItem'];
        configs?: components['schemas']['ConfigsResponse'];
        searchError?: {
            status?: number;
            message?: string;
        } & Record<string, unknown>;
    } = {

    };
    // TODO: send both requests in parallel
    if (app_id_valid) {

        const client = clientWithSvelteFetch(fetch, PUBLIC_API_BASE_URL_LOCAL);
        let infoResp: Awaited<ResponseType<'GET', '/v1/steam/appinfo'>> & {
            data?: components['schemas']['AppInfoItem'];
        };
        try {
            infoResp = await client.GET('/v1/steam/appinfo', {
                params: {
                    query: {
                        app_id,
                        raw: false,
                        controller_support: true,
                        official_configs: true
                    }
                }
            }) as typeof infoResp;
        } catch (err) {
            log.error('Failed to fetch app details', 'app_id', app_id, 'error', err);
            error(500, {
                message: 'An unexpected error occurred while fetching app details',
                err
            });
        }
        if (infoResp.error) {
            log.error('Failed to fetch app details', 'app_id', app_id, 'error', infoResp.error);
            error(infoResp.error.status || 503, {
                message: infoResp.error.detail || 'Failed to fetch app details',
                error: infoResp.error
            });
        }
        if (!infoResp.data || !infoResp.data.name) {
            infoResp.data = {
                app_id,
                name: `App ID: ${app_id}`,
                store_url_path: '',
                type: 'game'
            };
        }

        loadRes.appInfo = infoResp.data;
    }


    try {
        const searchParams = url.searchParams;
        searchParams.set('appid', appid);
        loadRes.configs = await fetchConfigs(fetch, searchParams);
    } catch (e) {
        log.error('Error fetching search results', 'error', e);
        if (isHttpError(e)) {
            loadRes.searchError = {
                status: e.status || 502,
                message: e.body?.message || 'Error contacting search endpoint',
                error: `${e.body}`
            };
        } else {
            loadRes.searchError = {
                status: 502,
                message: 'Error contacting search endpoint',
                error: `${e}`
            };
        }
    }

    return loadRes;

};
