<script lang="ts">
import { onNavigate } from '$app/navigation';
import favicon from '$lib/assets/favicon.svg?url';
import Footer from '$lib/components/Footer.svelte';
import Header from '$lib/components/header/Header.svelte';

import { browser } from '$app/environment';
import { page } from '$app/state';
import { BuddyState } from '$lib/buddy-app/buddyState.svelte';
import { GamepadNavigator } from '$lib/gamepad/gamepadNavigator.svelte';
import { toast } from '$lib/toaster/toaster.svelte';
import { onMount, type Snippet } from 'svelte';
import { quadIn, quadOut } from 'svelte/easing';
import { slide } from 'svelte/transition';
import 'unfonts.css';
import { links } from 'unplugin-fonts/head';
import '../css/main.pcss';
const { children } = $props();

onNavigate((navigation) => {
	if (!document.startViewTransition) {
		return;
	}

	// prevent view transition for same-page navigations,
	// there should not be a fucking transition if nothing changes... 🙄
	if (navigation.from?.url.pathname === navigation.to?.url.pathname) {
		return;
	}

	return new Promise((resolve) => {
		document.startViewTransition(async () => {
			resolve();
			await navigation.complete;
		});
	});
});

onMount(() => {
	if (browser && page.data.buddyAppEnabled && !page.route.id?.includes('(buddy-app)')) {
		BuddyState.pingBuddy(fetch).then((reachable) => {
			if (!reachable) {
				toast({
					snippet: failedToReachBuddy as Snippet,
					color: 'firebrick'
				});
			}
		});
	}
	if (!browser) {
		return;
	}
	GamepadNavigator.Init();
	return () => {
		GamepadNavigator.Dispose();
	};
});
</script>

{#snippet failedToReachBuddy({ color }: { color: string })}
	<div
		role="alert"
		style="--color: {color}"
		in:slide|global={{ duration: 196, delay: 196, easing: quadOut }}
		out:slide|global={{ duration: 196, easing: quadIn }}>
		<strong style="font-size: 1.4em; margin-bottom: 0.5em;">Failed to reach SteamInputDB-Buddy</strong>
		<p>Is SteamInputDB-Buddy running and your browser allowed to make requests to "localhost"?</p>
	</div>
{/snippet}

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="icon" type="image/png" sizes="64x64" href="/favicon.png" />
	<link rel="icon" type="image/x-icon" href="/favicon.ico" />
	<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
	{#each links as link (link?.attrs?.href)}
		{#if link?.attrs?.onload}
			<link
				{...link?.attrs || {}}
				onload={function () {
					this.rel = 'stylesheet';
				}} />
		{:else}
			<link {...link?.attrs || {}} />
		{/if}
	{/each}
	<link rel="canonical" href={page.url.toString()} />
	<link
		rel="search"
		type="application/opensearchdescription+xml"
		href="/opensearch.xml"
		title="SteamInputDB" />
	<meta property="og:url" content={page.url.toString()} />
	<meta property="og:site_name" content="SteamInputDB" />
	<style>
	body,
	main,
	header,
	footer {
		transition: all var(--transition-duration) var(--default-ease);
	}
	</style>
</svelte:head>

<Header />
{@render children()}
<Footer />

<style lang="postcss">
:global(body) {
	display: grid;
	grid-template-rows: auto 1fr auto;
	min-height: 100svh;
	max-width: 100dvw;
}

:global(main) {
	grid-row: 2 / span 1;
	grid-column: 1 / span 1;
}

:global(footer) {
	grid-row: 3 / span 1;
	grid-column: 1 / span 1;
}
</style>
