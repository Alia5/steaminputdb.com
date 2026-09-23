import { describe, expect, it } from 'vitest';
import { BuildRequestTypeString, BuildResponseTypeString, CtorForTypeString, type Endpoint } from './magic';
import { CStoreQuerySearchSuggestionsRequest, CStoreQuerySearchSuggestionsResponse } from './proto/service_storequery.pb';


describe('TestMagicRuntime', () => {

    type TestFunc = <T extends Endpoint|string>(endpoint: T) => unknown;
    interface TestCase {
        name: string;
        input: Endpoint | string;
        expected: unknown;
        func: TestFunc;

    }

    const testCases: TestCase[] = [
        {
            name: 'BuildRequestTypeString',
            input: {
                interface: 'IStoreQueryService',
                method: 'SearchSuggestions'
            },
            expected: 'CStoreQuerySearchSuggestionsRequest',
            func: BuildRequestTypeString as TestFunc
        },
        {
            name: 'BuildResponseTypeString',
            input: {
                interface: 'IStoreQueryService',
                method: 'SearchSuggestions'
            },
            expected: 'CStoreQuerySearchSuggestionsResponse',
            func: BuildResponseTypeString as TestFunc
        },
        {
            name: 'CtorForTypeString(Request)',
            input: 'CStoreQuerySearchSuggestionsRequest',
            expected: CStoreQuerySearchSuggestionsRequest,
            func: CtorForTypeString
        },
        {
            name: 'CtorForTypeString(Response)',
            input: 'CStoreQuerySearchSuggestionsResponse',
            expected: CStoreQuerySearchSuggestionsResponse,
            func: CtorForTypeString
        }
    ];

    testCases.forEach((tc) => {
        it(tc.name, async () => {
            const result = tc.func(tc.input);
            expect(result).toBe(tc.expected);
        });
    });
});

