import { PUBLIC_API_BASE_URL_LOCAL } from '$env/static/public';
import { clientWithSvelteFetch } from '$lib/api/client';
import { log } from '$lib/log';
import type { LayoutServerLoad } from './$types';


export const load: LayoutServerLoad = async ({ cookies, url, fetch }) => {
    const res = {
        theme: cookies.get('theme'),
        buddyAppEnabled: false
    };

    if (url.toString().includes('buddy-app=enabled')) {
        cookies.set('buddy-app', 'enabled', { path: '/', httpOnly: false });
    }
    if (url.toString().includes('buddy-app=disabled')) {
        cookies.delete('buddy-app', { path: '/' });
    }
    res.buddyAppEnabled = cookies.get('buddy-app') === 'enabled';

    const token = cookies.get('token');
    if (!token) {
        return res;
    }

    const client = clientWithSvelteFetch(fetch, PUBLIC_API_BASE_URL_LOCAL);
    const infoResp = await client.GET('/v1/steam/userinfo', {
        headers: {
            cookie: `token=${token}`
        }
    });
    if (infoResp.error || !infoResp.data) {
        log.debug('Layout server load: failed to fetch user info', 'error', infoResp.error);
        return res;
    }
    const userInfo = infoResp.data;
    const steamId = userInfo.steamid;

    log.debug('Layout server load', 'steamid', steamId, 'userInfo', userInfo);

    return {
        ...res,
        steamId,
        userInfo
    };

};


