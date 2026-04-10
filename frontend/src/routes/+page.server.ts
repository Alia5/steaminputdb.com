import { resolve } from '$app/paths';
import { clientWithSvelteFetch } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import { log } from '$lib/log';
import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const actions = {
    search: async (event) => {
        const params = await event.request.formData();
        throw redirect(302, resolve( `/config/search?searchtext=${params.get('searchtext')}&sort-by=vote`));
    }
} satisfies Actions;


const HOT_UPDATE_PERIOD_HOURS = 6;
let lastHotUpdate = 0;
let hotToday = [] as components['schemas']['ConfigsResponse']['items'];


export const load: PageServerLoad = async (event) => {

    if (Date.now() - lastHotUpdate > HOT_UPDATE_PERIOD_HOURS * 3600 * 1000) {
    // if (true) {
        try {
            const res = await clientWithSvelteFetch(event.fetch).POST('/v1/search/configs', {
                body: {
                    limit: 15,
                    query_text: '',
                    raw: false,
                    page: 0,
                    rank: {
                        by: 'trend',
                        trending_period: 1
                    },
                    filter: {
                    },
                    include: {
                        votes: true,
                        tags: true
                    }
                }
            });
            if (res.error) {
                throw res.error;
            }
            lastHotUpdate = Date.now();
            hotToday = (res.data as components['schemas']['ConfigsResponse'])?.items
                ?.filter((c) => c.app_id && Number.isInteger(c.app_id))
                ?.sort((a, b) => (a.votes?.up || 0) > (b.votes?.up || 0) ? -1 : 1)
                ?? [];
        } catch (e) {
            log.error('Failed to fetch hot configs', e);
            hotToday = [];
        }
    }

    return {
        hotToday
    };

};
