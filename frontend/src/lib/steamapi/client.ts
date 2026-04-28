import {
    BuildRequestTypeString,
    BuildResponseTypeString,
    CtorForTypeString,
    type Endpoint,
    type RequestTypeEP,
    type ResponseTypeEP
} from './magic';


export async function GET<EP extends Endpoint, Req extends RequestTypeEP<EP>>(
    endpoint: EP,
    req?: Req,
    fetchFn: typeof fetch = fetch,
    apiKey?: string
): Promise<ResponseTypeEP<EP>> {
    try {
        const params = new URLSearchParams();
        if (req) {
            const requestCtor = CtorForTypeString(BuildRequestTypeString(endpoint));
            const reqBytes = requestCtor.encode(
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-ignore
                requestCtor?.create(req) ?? req
            ).finish();
            params.append(
                'input_protobuf_encoded',
                btoa(String.fromCharCode(...reqBytes))
            );
        }
        if (apiKey) {
            params.append('key', apiKey);
        }
        const URL = `https://api.steampowered.com/${
            // @ts-expect-error fuckoff
            endpoint.interface
        }/${
            // @ts-expect-error fuckoff
            endpoint.method
        }/v${
            // @ts-expect-error fuckoffI
            endpoint.version ?? '1'
        }/?${
            params.toString()
        }`;
        const resp = await fetchFn(URL, {
            method: 'GET'
        });
        if (!resp.ok) {
            throw new Error('Steam API HTTP Error', {
                cause: {
                    status: resp.status,
                    statusText: resp.statusText
                }
            });
        }
        const buf = new Uint8Array(await resp.arrayBuffer());
        const respCtor = CtorForTypeString(BuildResponseTypeString(endpoint));
        return respCtor.decode(buf) as ResponseTypeEP<EP>;
    } catch (e) {
        throw new Error('Failed to fetch from Steam API', { cause: e });
    }
}

export const demo = async () => {
    const resp = await GET(
        // @ts-expect-error just shut up
        { interface: 'IStoreQueryService', method: 'SearchSuggestions' },
        {
            context: {
                language: 'english',
                countryCode: 'US'
            },
            searchTerm: 'portal',
            maxResults: 2
        }
    );
    return resp;
};
