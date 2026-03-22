import { browser } from '$app/environment';
import type { components } from '$lib/buddy-api/openapi';
import { log } from '$lib/log';
import { pingBuddy } from '../../routes/(buddy-app)/buddy-app/requests';


class buddyState {
    public reachable = $derived(!browser || false);
    public pending = $state(!browser || false);
    public pingResponse = $state<undefined|components['schemas']['Ping']>(browser ? undefined : undefined);

    public async pingBuddy(fetch: typeof globalThis.fetch = globalThis.fetch) {
        this.pending = true;
        this.reachable = false;
        try {
            const resp = await pingBuddy(fetch);
            this.pingResponse = resp;
            this.reachable = true;
            return resp;
        } catch (e) {
            log.error('Failed to ping buddy app', 'error', e);
            this.pingResponse = undefined;
            this.reachable = false;
            return false;
        } finally {
            this.pending = false;
        }
    }
}

export const BuddyState = new buddyState();
