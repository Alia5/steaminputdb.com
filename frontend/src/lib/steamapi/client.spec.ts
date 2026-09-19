import { AssertionError } from 'assert';
import { describe, expect, it } from 'vitest';
import { GET } from './client';
import type { Endpoint, RequestTypeEP } from './magic';
import { CStoreQuerySearchSuggestionsRequest, CStoreQuerySearchSuggestionsResponse } from './proto/service_storequery.pb';


describe('GET', () => {

    interface TestCase {
        name: string;
        endpoint: Endpoint;
        req: unknown;
        expectedResp: unknown;
        expectError?: Error;
    }

    const tests: TestCase[] = [
        // Disabled: asserts an exact result order from Steam's live search API, which changes over time
        // {
        //     name: 'SUCCESS_CanMakeRequest(SearchSuggestions)',
        //     endpoint: {
        //         interface: 'IStoreQueryService',
        //         method: 'SearchSuggestions'
        //     },
        //     req: {
        //         context:{
        //             language: 'english',
        //             countryCode: 'US'
        //         },
        //         searchTerm: 'isaac',
        //         maxResults: 5,
        //         filters: {
        //             typeFilters:  {
        //                 includeGames: true
        //             }
        //         }
        //     },
        //     expectedResp: CStoreQuerySearchSuggestionsResponse.create({
        //         metadata: {
        //             totalMatchingRecords: 4,
        //             start: 0,
        //             count: 4
        //         },
        //         ids: [
        //             { appid: 250900 },
        //             { appid: 113200 },
        //             { appid: 1273600 },
        //             { appid: 341260 }
        //         ]
        //     })
        // },
        {
            name: 'SUCCESS_EmptyRequestResponse(SearchSuggestions)',
            endpoint: {
                interface: 'IStoreQueryService',
                method: 'SearchSuggestions'
            },
            req: CStoreQuerySearchSuggestionsRequest.create(),
            expectedResp: CStoreQuerySearchSuggestionsResponse.create()
        },
        {
            name: 'ERROR_InvalidRequest(SearchSuggestions)',
            endpoint: {
                interface: 'IStoreQueryService',
                method: 'SearchSuggestions'
            },
            req: { invalidField: 123 },
            expectedResp: null,
            expectError: new AssertionError({ message: 'expected { metadata: undefined, ids: [], …(1) } to deeply equal null' })
        },
        {
            name: 'ERROR_InvalidEndpoint',
            endpoint: ({
                interface: 'invalidService',
                method: 'invalidMethod'
            }) as unknown as Endpoint,
            req: undefined,
            expectedResp: null,
            expectError: new Error('Failed to fetch from Steam API', { cause: 'meh' })
        }
    ];

    tests.forEach((tc) => {
        it(tc.name, async () => {
            try {
                const resp = await GET(
                    tc.endpoint,
                    tc.req as RequestTypeEP<typeof tc.endpoint>
                );
                expect(resp).toEqual(tc.expectedResp);
            } catch (e) {
                if (tc.expectError) {
                    expect(typeof e).toEqual(typeof tc.expectError);
                    expect((e as Error).message).toEqual(tc.expectError.message);
                } else {
                    throw e;
                }
            }
        });
    });
});

