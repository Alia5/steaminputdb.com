import { PUBLIC_API_BASE_URL_LOCAL } from '$env/static/public';
import { clientWithSvelteFetch, type ResponseType } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import { log } from '$lib/log';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';


export const load: PageServerLoad = async ({ params, fetch }) => {

    const fileid = params.fileid;

    const file_id = parseInt(fileid, 10);
    if (isNaN(file_id) || !file_id) {
        error(400, 'Invalid file ID');
    }

    const client = clientWithSvelteFetch(fetch, PUBLIC_API_BASE_URL_LOCAL);
    let infoResp: Awaited<ResponseType<'GET', '/v1/steam/filedetails'>>;
    try {
        infoResp = await client.GET('/v1/steam/filedetails', {
            params: {
                query: {
                    file_id,
                    playtime_stats: 14,
                    raw: false
                }
            }
        });
    } catch (err) {
        log.error('Failed to fetch file details', 'file_id', file_id, 'error', err);
        error(500, {
            message: 'An unexpected error occurred while fetching file details',
            err
        });
    }
    if (infoResp.error) {
        log.error('Failed to fetch file details', 'file_id', file_id, 'error', infoResp.error);
        error(infoResp.error.status || 503, {
            message: infoResp.error.detail || 'Failed to fetch file details',
            error: infoResp.error
        });
    }
    if (!infoResp.data) {
        error(404, 'File not found');
    }

    const fileInfo = infoResp.data as components['schemas']['ConfigItem'];
    const resData: {
        fileInfo: components['schemas']['ConfigItem'];
        nonSteam: boolean;
        appInfo?: components['schemas']['AppItem'];
        creatorInfo?: components['schemas']['UserInfoResponse'];
    } = {
        fileInfo,
        nonSteam: ! (!!fileInfo.app_id && Number.isInteger(fileInfo.app_id))
    };

    await Promise.allSettled([(async () => {
        if (fileInfo.app_id) {
            try {
                const appInfoResp = await client.GET('/v1/steam/appinfo/{app_id}', {
                    params: {
                        path: {
                            app_id: fileInfo.app_id
                        },
                        query: {
                            raw: false
                        }
                    },
                    signal: AbortSignal.timeout(10000)
                });
                if (appInfoResp.error) {
                    log.error(
                        'Failed to fetch app details',
                        'file_id', file_id,
                        'app_id', fileInfo.app_id,
                        'status', appInfoResp.error.status,
                        'error', appInfoResp.error
                    );
                }
                if (!appInfoResp.data) {
                    log.error('No app details data received',' file_id', file_id, 'app_id', fileInfo.app_id);
                }
                resData.appInfo = appInfoResp.data;
            } catch (err) {
                log.error('Failed to fetch app details', 'file_id', file_id, 'app_id', fileInfo.app_id, 'error', err);

            }
        }
    })(), (async () => {
        if (fileInfo.creator_id) {
            try {
                const creatorInfoResp = await client.GET('/v1/steam/userinfo', {
                    params: {
                        query: {
                            user_id: `${fileInfo.creator_id}`
                        }
                    },
                    signal: AbortSignal.timeout(10000)

                });
                if (creatorInfoResp.error) {
                    log.error(
                        'Failed to fetch creator details',
                        'file_id', file_id,
                        'creator_id', fileInfo.creator_id,
                        'status', creatorInfoResp.error.status,
                        'error', creatorInfoResp.error
                    );
                }
                if (!creatorInfoResp.data) {
                    log.error('No creator details data received',' file_id', file_id,'creator_id', fileInfo.creator_id);
                }
                resData.creatorInfo = creatorInfoResp.data;
            } catch (err) {
                log.error('Failed to fetch creator details', 'file_id', file_id, 'creator_id', fileInfo.creator_id, 'error', err);
            }
        }
    })()]);

    return resData;

};
