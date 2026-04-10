<script lang="ts">
import { afterNavigate, beforeNavigate } from '$app/navigation';
import { resolve } from '$app/paths';
import { page } from '$app/state';
import { client } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import SC2 from '$lib/assets/SC2_Googley.svg.svelte';
import { intersectionObserver } from '$lib/attachments/intersectionObserver.svelte';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import SearchForm from '$lib/components/search/searchform.svelte';
import Spinner from '$lib/components/Spinner.svelte';
import { log } from '$lib/log';
import { createHomeSchemaJsonLd } from '$lib/schema/home';
import { configRating } from '$lib/snippets/configRating.svelte';
import { configurationFeatureList } from '$lib/snippets/configurationfeaturelist.svelte';
import { controllertype } from '$lib/snippets/controllertype.svelte';
import { statusCodeNames } from '$lib/statuscodes';
import { assetUrlBase } from '$lib/steamapi/const';
import Icon from '@iconify/svelte';
import { formatDistance } from 'date-fns';
import { onMount } from 'svelte';
import { cubicIn, cubicInOut, cubicOut } from 'svelte/easing';
import { fade, fly, slide } from 'svelte/transition';
import { fetchConfigs, PAGE_SIZE } from '../../../lib/api/searchConfigs';
import type { PageProps } from './$types';

let { data }: PageProps = $props();

// svelte-ignore state_referenced_locally
let results = $state(data?.results);
$effect(() => {
	if (data?.hasSearched) {
		results = data?.results;
	}
});

let loadingMore = $state(false);
let showBackToTop = $state(false);
const loadMore = async () => {
	log.debug('Load more triggered');
	if (!results?.items) {
		log.error('No results to load more for');
		return;
	}
	if (loadingMore) {
		return;
	}

	loadingMore = true;
	const searchParams = page.url.searchParams;
	searchParams.set('page', `${Math.floor(results.items.length / PAGE_SIZE) + 1}`);
	try {
		const res = await fetchConfigs(fetch, searchParams);
		results.items = results.items
			.concat(res?.items ?? [])
			.filter((item, index, self) => index === self.findIndex((i) => i.file_id === item.file_id));
	} catch (e) {
		log.error('Error loading more results', 'error', e);
	}
	loadingMore = false;
};

let form = $state<HTMLFormElement>()!;
let loading = $state(false);
let hasSearched = $derived(data.hasSearched);
let searchError = $derived(data.searchError);

let infoAppIdMap = $state<Record<number, components['schemas']['AppItem']>>({});

$effect(() => {
	results?.items?.forEach((cfg) => {
		if (!cfg.app_id) {
			return;
		}
		if (infoAppIdMap[cfg.app_id]) {
			return;
		}
		const idCopy = cfg.app_id;
		infoAppIdMap[idCopy] = {
			app_id: idCopy
		} as components['schemas']['AppItem'];
		client
			.GET('/v1/steam/appinfo', {
				params: {
					query: {
						app_id: idCopy
					}
				}
			})
			.then((resp) => {
				if (resp.error) {
					log.error('Error fetching store info', 'appid', idCopy, 'error', resp.error);
					return;
				}
				if (!resp.data) {
					log.error('No data in response when fetching store info', 'appid', idCopy);
					return;
				}
				infoAppIdMap[idCopy] = resp.data as components['schemas']['AppItem'];
			})
			.catch((err) => {
				log.error('Error fetching store info', 'appid', idCopy, 'error', err);
				delete infoAppIdMap[idCopy];
			});
	});
});

let formValues = $state(
	Array.from(page.url.searchParams).reduce((acc, [key, value]) => {
		return Object.assign(acc, { [key]: value });
	}, {})
);

let eyes = $state<{
	left: HTMLElement;
	right: HTMLElement;
}>({
	left: undefined!,
	right: undefined!
})!;
const findEyes = () => {
	const group = document.querySelector('#sc2>*>svg>g>:last-child');
	if (!group) {
		return;
	}
	eyes.left = group.children[0] as HTMLElement;
	eyes.right = group.children[1] as HTMLElement;
};
onMount(() => {
	findEyes();
});

beforeNavigate((event) => {
	if (event.type == 'form') {
		loading = true;
	}
});
afterNavigate(() => {
	loading = false;
});
</script>

<svelte:head>
	<title>SteamInputDB - Config Search</title>
	<meta
		name="keywords"
		content="Steam Input DB, Steam DB, DB, Steam Deck, Steam Input, Steam controller configs, controller layouts, community database, Steam API, gamepad configurations, controller presets" />
	<meta
		name="robots"
		content="index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1" />

	<meta property="og:site_name" content="SteamInputDB" />
	<meta property="og:type" content="website" />
	<meta property="og:title" content="SteamInputDB" />
	<meta name="description" content="Search for Steam Input configurations for Steam and non Steam games" />
	<meta
		property="og:description"
		content="Search for Steam Input configurations for Steam and non Steam games" />
	<meta property="og:url" content="https://www.steaminputdb.com/" />
	<meta property="og:image" content="https://www.steaminputdb.com/ogimage.png" />
	<meta property="og:image:alt" content="SteamInputDB preview image" />
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="SteamInputDB - Config Search" />
	<meta
		name="twitter:description"
		content="Search for Steam Input configurations for Steam and non Steam games" />
	<meta name="twitter:image" content="https://www.steaminputdb.com/ogimage.png" />
	<meta name="twitter:image:alt" content="SteamInputDB preview image" />

	<svelte:element this={'script'} type="application/ld+json">{createHomeSchemaJsonLd()}</svelte:element>
</svelte:head>

<svelte:window
	onmousemove={(e) => {
		if (!eyes || !eyes.left || !eyes.right) {
			return;
		}
		Object.values(eyes).forEach((eye, idx) => {
			const rect = eye.getBoundingClientRect();
			const eyeX = rect.left + rect.width / 2;
			const eyeY = rect.top + rect.height / 2;
			const deltaX = e.clientX - eyeX;
			const deltaY = e.clientY - eyeY;
			const angle = Math.atan2(deltaY, deltaX);
			const distance = Math.min(16, Math.hypot(deltaX, deltaY) / 2);
			const translateX = Math.cos(angle) * distance * (searchError ? (idx === 0 ? 1 : -1) : 1);
			const translateY = Math.sin(angle) * distance * (searchError ? (idx === 0 ? -1 : 1) : 1);
			eye.style.transform = `translate(${translateX}px, ${translateY}px)`;
		});
	}}
	onscroll={() => {
		showBackToTop = window.scrollY > window.innerHeight;
	}} />

<main>
	<search>
		<SearchForm
			bind:form={form}
			disabled={loading}
			method="GET"
			bind:values={formValues}
			submitOnChange={true}
			showTotalCount={results?.total} />
		<div class="results">
			{#if loading}
				<div
					class="loading"
					in:fade|global={{ duration: 196, easing: cubicOut }}
					out:fade|global={{ duration: 196, easing: cubicIn }}>
					<Spinner size="12em" />
				</div>
			{/if}
			{#if !loading && (searchError || (results?.items?.length ?? 0) == 0)}
				<div
					id="sc2"
					class={searchError ? 'error' : ''}
					in:fade|global={{ duration: 196, easing: cubicOut }}
					out:fade|global={{ duration: 196, easing: cubicIn }}
					onintrostartcapture={findEyes}>
					{#if hasSearched && !searchError}
						<span>No results found</span>
					{/if}
					{#if searchError}
						{#if searchError.status}
							<h2>
								{searchError.status}
								{#if statusCodeNames[searchError.status] && statusCodeNames[searchError.status] !== searchError?.message}
									{statusCodeNames[searchError.status]}
								{/if}
							</h2>
						{/if}
						<span>{searchError?.message ?? 'An unknown error occurred'}</span>
					{/if}
					<SC2
						height="100%"
						--eyes-color="black"
						--eyes-white-color="var(--text-color-dark)"
						--eyes-border-color="light-dark(var(--text-color-light), transparent)" />
				</div>
			{/if}
			{#if !searchError}
				<div class="list" transition:fade={{ duration: 196, easing: cubicInOut }}>
					{#each results?.items ?? [] as item (item.file_id)}
						<a
							class="plain"
							style={loading ? 'pointer-events: none; opacity: 0.4;' : ''}
							href={resolve(`/config/${item.file_id}`)}
							transition:slide|global={{ duration: 196, easing: cubicInOut }}>
							<div class="thumb">
								{#if infoAppIdMap?.[item.app_id || 0]}
									{@const assets = infoAppIdMap[item.app_id || 0]!.assets!}
									{@const srcChosen = assets?.asset_url_format
										? `${assetUrlBase}${assets?.asset_url_format?.replace(
												'${FILENAME}',
												assets.small_capsule ??
													assets.main_capsule ??
													assets.header ??
													assets.hero_capsule ??
													assets.library_hero ??
													assets.library_capsule ??
													assets.page_background ??
													'undefined'
											)}`
										: undefined}
									{#if srcChosen}
										<picture transition:fade={{ duration: 196, easing: cubicOut }}>
											<enhanced:img src={srcChosen} alt="Thumbnail" height="100%"
											></enhanced:img>
										</picture>
									{/if}
								{/if}
								<i>
									{#if item.app_id}
										<Icon icon="mdi:steam" width="1.2em" />
									{:else}
										<Icon icon="mdi:link-variant" width="1.2em" />
									{/if}
									{infoAppIdMap?.[item.app_id || 0]?.name ?? item.app_id_string}
								</i>
							</div>
							<div class="info">
								<div>
									<strong class="title">{item.title}</strong>
									<span>
										<span>{@render controllertype({ item })}</span>
										{#if item.lifetime_playtime_seconds}
											<span
												>{formatDistance(
													new Date((item.lifetime_playtime_seconds ?? 0) * 1000),
													new Date(0)
												)}
												combined playtime
											</span>
										{/if}
									</span>
									<span>{item.description}</span>
									<div>
										{@render configurationFeatureList({ fileInfo: item })}
									</div>
								</div>
								<div>
									{@render configRating({ item })}
								</div>
							</div>
						</a>
					{/each}
					{#if (results?.items?.length ?? 0) < (results?.total ?? 0)}
						<div
							id="load-more-trigger"
							{@attach intersectionObserver(() => {
								loadMore();
							})}>
							<Spinner size="12em" />
						</div>
					{/if}
				</div>
			{/if}
		</div>
		{#if showBackToTop}
			<button
				out:fade={{ duration: 196, easing: cubicIn }}
				in:fly={{ y: '2dvh', duration: 196, easing: cubicOut }}
				id="back-to-top"
				{@attach tooltip({ content: 'Back to top' })}
				onclick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}>
				<Icon icon="mdi:arrow-up" width="1.5em" />
			</button>
		{/if}
	</search>
</main>

<style lang="postcss">
main {
	display: grid;
	place-items: center;
}

search {
	display: grid;
	place-self: center;
	grid-template-rows: min-content auto;
	gap: 1em;
	place-items: center;
	min-width: 60%;
	--max-width: 1440px;
	max-width: min(100%, var(--max-width));
	isolation: isolate;

	min-height: min(100%, 55dvh);

	& > :first-child {
		width: 100%;
	}
	margin-top: 1em;
}

.results {
	display: grid;
	place-items: center;
	width: 100%;
	grid-template-columns: auto;
	grid-template-rows: 1fr;
	padding: 1em 0;
	height: 100%;

	& > * {
		grid-row: 1;
		grid-column: 1;
	}
}

#sc2 {
	filter: drop-shadow(0px 0.25em 0.2em var(--shadow-color));
	width: 100%;
	display: grid;
	place-items: center;
	padding: 0.5em;
	gap: 1em;

	:global(svg) {
		max-height: 26em;
		max-width: min(100%, 26em);
		margin-bottom: auto;
	}

	@media (any-pointer: coarse) {
		:global(ellipse) {
			transition: transform calc(var(--transition-duration) * 2) var(--default-ease);
		}
	}

	& span {
		color: var(--highlight-color);
		font-size: 2em;
		font-weight: bold;
	}

	&.error {
		:global(svg) {
			transform: rotate(180deg);
		}

		:global(#qam),
		:global(#pad_r),
		:global(#y),
		:global(#x) {
			display: none;
		}

		:global(#logo_overlay) {
			fill: var(--background-color);
		}

		:global(#dpad) {
			transform: rotate(3deg);
			transform-origin: 25% 25%;
			translate: -0.1em 0.1em;
		}

		h2 {
			color: var(--highlight-color);
			font-size: 2em;
		}

		& span {
			color: var(--text-color);
			font-size: 1.4em;
			font-weight: bold;
		}
	}
}

.results .list {
	width: 100%;
	display: grid;
	position: relative;

	& #load-more-trigger {
		justify-self: center;
		padding: 2em 1em;
		overflow: hidden;
	}
}

a {
	display: grid;
	grid-template-columns: minmax(3em, 24em) auto;

	@media (max-width: 1000px) {
		grid-template-columns: 100%;
	}

	--gap: 1em;
	column-gap: var(--gap);
	row-gap: 0;
	padding: calc(0.5 * var(--gap)) calc(0.5 * var(--gap));
	width: 100%;
	color: var(--text-color);
	position: relative;
	isolation: isolate;
	&::before {
		content: '';
		height: 0.1em;
		width: 100%;
		bottom: calc(0.1em * -0.5);
		left: 0;
		position: absolute;
		background: color-mix(in srgb, currentColor, transparent 80%);
	}

	&:hover,
	&:focus-visible {
		&::after {
			opacity: 0.3;
		}
	}
	&::after {
		content: '';
		position: absolute;
		inset: 0px;
		background: var(--color-primary);
		opacity: 0;
		z-index: -1;
		transition: opacity var(--transition-duration) var(--default-ease);
	}

	& > div {
		overflow: hidden;
		width: 100%;
		display: grid;
	}
}

.title {
	text-overflow: ellipsis;
	overflow: hidden;
	white-space: nowrap;
	font-size: 1.2em;
}

.thumb {
	aspect-ratio: 21 / 8;
	width: auto;
	height: 100%;
	width: 100%;

	border-radius: 0.5em;
	box-shadow: var(--card-shadow);

	color: var(--text-color-dark);
	isolation: isolate;

	align-self: center;

	& picture,
	& img {
		object-fit: cover;
		object-position: center;
		width: 100%;
		height: 100%;
		overflow: hidden;
		z-index: -1;
	}

	& picture {
		isolation: isolate;
		position: relative;
		&::before {
			content: '';
			position: absolute;
			inset: 0px;
			mask: linear-gradient(
				170deg,
				rgb(255 255 255 / 0.9) 0%,
				rgb(255 255 255 / 0.7) 25%,
				transparent 50%
			);
			mask-type: alpha;
			backdrop-filter: saturate(0.1) brightness(0.3);
			z-index: 1;
		}
	}

	display: grid;
	background: linear-gradient(rgb(90 90 90 / 1), rgb(150 150 150 / 0.9));
	& > * {
		grid-row: 1;
		grid-column: 1;
	}

	& i {
		font-size: 1em;
		padding: 0.5em;
		filter: drop-shadow(1px 1px 1px black) drop-shadow(0px 0px 4px rgb(0 0 0 / 0.25))
			drop-shadow(0px 0px 2px black);
		display: grid;
		grid-auto-flow: column;
		align-items: center;
		justify-content: start;
		height: fit-content;
		gap: 0.25em;
		font-weight: bold;
	}
}

.info {
	display: flex;
	flex-flow: row wrap;
	gap: 0.5em;

	& > :first-child {
		display: grid;
		gap: 0.5em;
		flex: 1;
		min-width: min(calc(100dvw - 2em), 300px);

		& > strong {
			padding-bottom: 0.5em;
		}

		overflow: clip;
		overflow-clip-margin: 1em;

		& > div {
			display: flex;
			flex-flow: row wrap;
			gap: 0.5em;
			overflow: clip;
			overflow-clip-margin: 1em;
		}

		& > :nth-child(2) {
			display: flex;
			flex-flow: row wrap;
			align-items: center;
			width: fit-content;
			gap: 1em;
			& > :first-child {
				display: flex;
				flex-flow: row wrap;
				align-items: center;
				color: var(--highlight-color);
				gap: 0.25em;
				font-weight: bold;
				font-size: 1.1em;
			}
			& > :nth-child(2) {
				opacity: 0.8;
			}
		}
		& > :nth-child(3) {
			white-space: nowrap;
			text-overflow: ellipsis;
			overflow: hidden;
			max-width: 100%;
		}

		& > :last-child {
			margin: auto 0;
			height: fit-content;
		}
	}
}

.loading {
	align-self: baseline;
	z-index: 1;
}

#back-to-top {
	position: fixed;
	bottom: 2dvh;
	right: 2dvw;
	border-radius: 100dvh;
	background:
		linear-gradient(
			215deg,
			color-mix(in srgb, var(--card-color), transparent 75%) 0%,
			color-mix(in srgb, var(--card-color), transparent 90%) 70%
		),
		var(--bg-noise-transparent);
	background-color: color-mix(in srgb, var(--color-primary), transparent 20%);
}
</style>
