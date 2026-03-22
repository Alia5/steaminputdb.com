import type { PageLoad } from './$types';
import { getLatestVersion } from './requests';


// eslint-disable-next-line arrow-body-style
export const load: PageLoad = async ({ fetch }) => {
    return {
        latestVersion: getLatestVersion(fetch)
    };
};
