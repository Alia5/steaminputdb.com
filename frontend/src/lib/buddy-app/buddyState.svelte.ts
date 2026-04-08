import { browser } from '$app/environment';
import { clientWithSvelteFetch } from '$lib/buddy-api/client';
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

    public connectedControllers = $state<undefined|components['schemas']['ControllerResponse'][]>(undefined);
    public connectedControllersChanged = $state(false);
    public async fetchConnectedControllers(fetch: typeof globalThis.fetch = globalThis.fetch) {
        try {
            const r = await clientWithSvelteFetch(fetch).GET('/v1/steam/controllers');
            if (r.data) {
                this.connectedControllers = r.data as components['schemas']['ControllerResponse'][];
                this.connectedControllersChanged = true;
                setTimeout(() => {
                    this.connectedControllersChanged = false;
                }, 0);
            }
            return r;
        } catch (e) {
            log.error('Failed to fetch connected controllers from buddy app', 'error', e);
            return undefined;
        }
    }

}

export const BuddyState = new buddyState();
