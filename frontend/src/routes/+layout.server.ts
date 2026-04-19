import type { components } from '$lib/api/openapi';
import { log } from '$lib/log';
import type { LayoutServerLoad } from './$types';


export const load: LayoutServerLoad = async ({ cookies, url }) => {
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
    const mid = token.split('.')?.[1];
    if (!mid) {
        return res;
    }
    const decoded = atob(mid);
    const payload = JSON.parse(decoded) as components['schemas']['UserInfoResponse'] & {
        sub?: string;
        is_admin?: boolean;
    };
    const steamId = payload.sub as string | undefined;

    log.debug('Layout server load', 'steamid', steamId, 'payload', payload);

    return {
        ...res,
        steamId,
        userInfo: payload
    };

};


