import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from '../$types';
import type { PageData } from './$types';

export const load: LayoutServerLoad = async ({ parent }) => {
    const data: PageData = await parent();
    if (data.steamId && data.userInfo) {
        throw redirect(302, '/');
    }
    return data;
};
