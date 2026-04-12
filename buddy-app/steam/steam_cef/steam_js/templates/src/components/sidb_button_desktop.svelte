<script lang="ts" module>
export const ATTR = 'data-sidb-button-desktop';
</script>

<script lang="ts">
import Logo from '$lib/assets/logo.svg?component';

const {
	openUrl,
	...props
}: {
	openUrl: (url: string) => void;
} & Record<string, unknown> = $props();

let container = $state<HTMLElement>()!;

const onclick = () => {
	try {
		const appIdPath = window.opener.SteamUIStore.ActiveWindowInstance.m_locationPathname;
		const appId = Number.parseInt(
			appIdPath.split('/').find((p: string) => Number.isInteger(Number.parseInt(p, 10))) ?? '0',
			10
		);
		const displayName = window.opener.appDetailsStore.m_mapAppData?.get(appId)?.details?.strDisplayName;
		const isNonSteam = !!window.opener.appDetailsStore.m_mapAppData?.get(appId)?.details?.strShortcutExe;
		if (isNonSteam) {
			openUrl(`https://steaminputdb.com/app/${encodeURIComponent(displayName)}?buddy-app=enabled`);
			return;
		}
		openUrl(`https://steaminputdb.com/app/${appId}?buddy-app=enabled`);
	} catch {
		openUrl('https://steaminputdb.com/config/search?buddy-app=enabled');
	}
};
</script>

<div bind:this={container} {...props} data-sidb-button-desktop aria-label="SteamInputDB" {onclick}>
	<Logo />
</div>

<style>
div {
	&:hover,
	&:focus-visible {
		:global(> *) {
			opacity: 1;
		}
	}
	:global(> *) {
		opacity: 0.6;
	}
}
</style>
