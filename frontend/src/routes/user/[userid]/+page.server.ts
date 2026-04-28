import { PUBLIC_API_BASE_URL_LOCAL } from '$env/static/public';
import { clientWithSvelteFetch, type ResponseType } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import { fetchConfigs } from '$lib/api/searchConfigs';
import { log } from '$lib/log';
import { error, isHttpError } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';


export const load: PageServerLoad = async ({ params, fetch, url }) => {

    const userid = params.userid;

    const user_id = parseInt(userid, 10);
    const user_id_valid = !isNaN(user_id) && user_id > 0;

    if (!user_id_valid) {
        log.error('Invalid user ID in URL', 'userid', userid);
        error(400, 'Invalid user ID');
    }

    const loadRes: {
        playerInfo?: components['schemas']['AppItem'];
        configs?: components['schemas']['ConfigsResponse'];
        searchError?: {
            status?: number;
            message?: string;
        } & Record<string, unknown>;
    } = {

    };

    const client = clientWithSvelteFetch(fetch, PUBLIC_API_BASE_URL_LOCAL);
    let infoResp: Awaited<ResponseType<'GET', '/v1/steam/userinfo'>> & {
        data?: components['schemas']['AppItem'];
    };


    await Promise.all([ (async () => {
        try {
            infoResp = await client.GET('/v1/steam/userinfo', {
                params: {
                    query: {
                        user_id: userid,
                        raw: false,
                        include_avatar_frame: true,
                        include_mini_profile_background: true,
                        include_profile_background: true
                    }
                }
            });
        } catch (err) {
            log.error('Failed to fetch user details', 'user_id', user_id, 'error', err);
            error(500, {
                message: 'An unexpected error occurred while fetching user details',
                err
            });
        }
        if (infoResp.error) {
            log.error('Failed to fetch user details', 'user_id', user_id, 'error', infoResp.error);
            error(infoResp.error.status || 503, {
                message: infoResp.error.detail || 'Failed to fetch user details',
                error: infoResp.error
            });
        }
        if (!infoResp.data || !infoResp.data.personaname) {
            error(404, 'User not found');
        }
        loadRes.playerInfo = infoResp.data;

    })(), (async () => {
        try {
            const searchParams = url.searchParams;
            searchParams.set('userid', userid);
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
    })()]);

    return loadRes;
};
