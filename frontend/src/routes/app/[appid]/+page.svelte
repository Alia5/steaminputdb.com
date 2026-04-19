<script lang="ts">
import { afterNavigate, beforeNavigate } from '$app/navigation';
import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { components } from '$lib/api/openapi';
import { fetchConfigs, PAGE_SIZE } from '$lib/api/searchConfigs';
import SC2 from '$lib/assets/SC2_Googley.svg.svelte';
import { intersectionObserver } from '$lib/attachments/intersectionObserver.svelte';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import SearchForm from '$lib/components/search/searchform.svelte';
import Spinner from '$lib/components/Spinner.svelte';
import { log } from '$lib/log';
import { createAppSchemaJsonLd } from '$lib/schema/app';
import { configRating } from '$lib/snippets/configRating.svelte';
import { configurationFeatureList } from '$lib/snippets/configurationfeaturelist.svelte';
import { controllertype } from '$lib/snippets/controllertype.svelte';
import { statusCodeNames } from '$lib/statuscodes';
import { assetUrlBase, storePageBackgroundBase } from '$lib/steamapi/const';
import Icon from '@iconify/svelte';
import { formatDistance } from 'date-fns';
import { onMount } from 'svelte';
import { cubicIn, cubicInOut, cubicOut } from 'svelte/easing';
import { fade, fly, slide } from 'svelte/transition';
import type { PageProps } from './$types';
import { sectionHead } from './sectionHead.svelte';

let { data }: PageProps = $props();

const appInfo: components['schemas']['AppInfoItem'] | undefined = $derived(data.appInfo);

const pageBGURL = $derived.by(() => {
	if (!appInfo?.assets) {
		return;
	}
	if (appInfo.assets.page_background) {
		return `${assetUrlBase}${appInfo.assets.asset_url_format?.replace(
			'${FILENAME}',
			appInfo.assets.page_background
		)}`;
	}
	if (appInfo.assets.raw_page_background) {
		return `${storePageBackgroundBase}${appInfo.assets.asset_url_format?.replace(
			'${FILENAME}',
			appInfo.assets.raw_page_background
		)}`;
	}
	if (appInfo.assets.page_background_path) {
		return `${storePageBackgroundBase}${appInfo.assets.page_background_path}`;
	}
});

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

let form = $state<HTMLFormElement>()!;
let loading = $state(false);
let searchError = $derived(data.searchError);
let formValues = $state(
	Array.from(page.url.searchParams).reduce((acc, [key, value]) => {
		return Object.assign(acc, { [key]: value });
	}, {})
);

// svelte-ignore state_referenced_locally
// eslint-disable-next-line svelte/prefer-writable-derived
let results = $state(data?.configs);
$effect(() => {
	results = data?.configs;
});

let showBackToTop = $state(false);
let loadingMore = $state(false);
const loadMore = async () => {
	log.debug('Load more triggered');
	if (!results?.items) {
		log.error('No results to load more for');
		return;
	}
	if (loadingMore) {
		return;
	}
	if (!page.params.appid) {
		return;
	}

	loadingMore = true;
	const searchParams = page.url.searchParams;
	searchParams.set('appid', page.params.appid);
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

beforeNavigate((event) => {
	if (event.type == 'form') {
		loading = true;
	}
});
afterNavigate(() => {
	loading = false;
});

onMount(() => {
	findEyes();
});
</script>

<svelte:head>
	<title>SteamInputDB - {appInfo?.name ?? page.params.appid}</title>
	<meta property="og:site_name" content="SteamInputDB" />
	<meta property="og:type" content="website" />
	<meta property="og:url" content={page.url.href} />
	<meta property="og:title" content="SteamInputDB - {appInfo?.name ?? page.params.appid}" />
	<meta
		name="description"
		content="Search for Steam Input configurations for {appInfo?.name ?? page.params.appid}" />
	<meta
		property="og:description"
		content="Search for Steam Input configurations for {appInfo?.name ?? page.params.appid}" />
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="SteamInputDB - {appInfo?.name ?? page.params.appid}" />
	<meta
		name="twitter:description"
		content="Search for Steam Input configurations for {appInfo?.name ?? page.params.appid}" />
	{#if appInfo?.assets}
		{@const assets = appInfo?.assets}
		{@const assetChosen =
			assets.main_capsule ?? assets.header ?? assets.hero_capsule ?? assets.library_hero ?? 'none.svg'}
		{#if assetChosen}
			<meta
				property="og:image"
				content={`${assetUrlBase}${assets.asset_url_format?.replace('${FILENAME}', assetChosen)}`} />
			<meta property="og:image:alt" content="SteamInputDB - {appInfo?.name ?? page.params.appid}" />
			<meta
				name="twitter:image"
				content={`${assetUrlBase}${assets.asset_url_format?.replace('${FILENAME}', assetChosen)}`} />
			<meta name="twitter:image:alt" content="SteamInputDB - {appInfo?.name ?? page.params.appid}" />
		{/if}
	{/if}
	<svelte:element this={'script'} type="application/ld+json">
		{createAppSchemaJsonLd({
			appId: page.params.appid ?? '',
			appName: appInfo?.name ?? page.params.appid ?? '',
			imageUrl:
				appInfo?.assets?.asset_url_format &&
				(appInfo?.assets?.main_capsule ??
					appInfo?.assets?.header ??
					appInfo?.assets?.hero_capsule ??
					appInfo?.assets?.library_hero)
					? `${assetUrlBase}${appInfo.assets.asset_url_format.replace(
							'${FILENAME}',
							appInfo.assets.main_capsule ??
								appInfo.assets.header ??
								appInfo.assets.hero_capsule ??
								appInfo.assets.library_hero ??
								''
						)}`
					: undefined
		})}
	</svelte:element>
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

<main style={pageBGURL ? `--bg: url('${pageBGURL}')` : ''}>
	<div>
		<div>
			{@render sectionHead({ appInfo, fallbackName: page.params.appid })}
		</div>
		<!-- TODO: show controller support and whatnot -->
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
						{#if !searchError}
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
								<div class="info">
									<div>
										<strong class="title">{item.title}</strong>
										<span>
											<span>{@render controllertype({ item })}</span>
											{#if item.lifetime_playtime_seconds}
												<span
													>{formatDistance(
														new Date(
															(item.lifetime_playtime_seconds ?? 0) * 1000
														),
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
	</div>
</main>

<style lang="postcss">
main {
	position: relative;
	isolation: isolate;
	display: grid;
	padding: 1em 0;
	width: 100%;
	&::before {
		content: '';
		position: absolute;
		left: 0;
		right: 0;
		top: 0;
		height: min(100%, 100dvh);
		background: var(--bg, transparent) top/cover no-repeat;
		mask: linear-gradient(0deg, transparent, white 100%);
		mask-type: alpha;
		z-index: -2;
	}
	& > :first-child {
		display: grid;
		place-self: baseline center;
		gap: 1em;
		place-items: center;
		min-width: 50%;
		--max-width: 1440px;
		max-width: min(100%, var(--max-width));
		isolation: isolate;
		& > :first-child {
			width: 100%;
			color: var(--text-color-dark);
		}
	}
}

search {
	display: grid;
	place-self: center;
	grid-template-rows: min-content auto;
	gap: 1em;
	place-items: center;

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
	isolation: isolate;

	&::before {
		content: '';
		position: absolute;
		left: 0;
		right: 0;
		top: 0;
		height: min(100%, 100dvh);
		background: var(--background-noise);
		z-index: -1;
		opacity: 0.5;
		border-radius: 1em;
	}

	& #load-more-trigger {
		justify-self: center;
		padding: 2em 1em;
		overflow: hidden;
	}

	& > :first-child {
		&::after {
			border-radius: 1em 1em 0 0;
		}
	}
}

.list > a {
	display: grid;
	--gap: 1em;
	padding: calc(0.5 * var(--gap)) calc(0.5 * var(--gap));
	column-gap: 1em;
	row-gap: 0;
	grid-template-columns: auto;
	width: 100%;
	color: var(--text-color);
	position: relative;
	isolation: isolate;
	padding: 1em;

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

	& .title {
		text-overflow: ellipsis;
		overflow: hidden;
		white-space: nowrap;
		font-size: 1.3em;
	}
	& .info {
		display: flex;
		flex-flow: row wrap;
		gap: 0.5em;

		overflow: clip;
		overflow-clip-margin: 1em;

		& > :first-child {
			display: grid;
			gap: 0.5em;
			flex: 1;
			flex-shrink: 1;
			min-width: min(100%, 300px);

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
					font-size: 1.2em;
				}
				& > :nth-child(2) {
					opacity: 0.8;
				}
			}
			& > :nth-child(3) {
				text-overflow: ellipsis;
				overflow: hidden;
				max-width: 100%;
				height: 3em;
				padding: 0 3em;
			}

			& > :last-child {
				margin: auto 0;
				height: fit-content;
			}
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
