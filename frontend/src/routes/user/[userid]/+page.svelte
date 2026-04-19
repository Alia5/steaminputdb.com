<script lang="ts">
import { afterNavigate, beforeNavigate } from '$app/navigation';
import { resolve } from '$app/paths';
import { page } from '$app/state';
import { client } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import { fetchConfigs, PAGE_SIZE } from '$lib/api/searchConfigs';
import SC2 from '$lib/assets/SC2_Googley.svg.svelte';
import { intersectionObserver } from '$lib/attachments/intersectionObserver.svelte';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import SearchForm from '$lib/components/search/searchform.svelte';
import Spinner from '$lib/components/Spinner.svelte';
import { log } from '$lib/log';
import { createUserSchemaJsonLd } from '$lib/schema/user';
import { configRating } from '$lib/snippets/configRating.svelte';
import { configurationFeatureList } from '$lib/snippets/configurationfeaturelist.svelte';
import { controllertype } from '$lib/snippets/controllertype.svelte';
import { statusCodeNames } from '$lib/statuscodes';
import { assetUrlBase, communityUrlBase } from '$lib/steamapi/const';
import Icon from '@iconify/svelte';
import { formatDistance } from 'date-fns';
import { onMount } from 'svelte';
import { cubicIn, cubicInOut, cubicOut } from 'svelte/easing';
import { fade, fly, slide } from 'svelte/transition';
import type { PageProps } from './$types';

let { data }: PageProps = $props();

const playerInfo = $derived(
	data.playerInfo as unknown as components['schemas']['UserInfoResponse'] | undefined
);

const pageBGURL = $derived.by(() => {
	if (!playerInfo?.profilebackground) {
		return undefined;
	}
	return `${communityUrlBase.replace('/apps/', '/')}${playerInfo.profilebackground}`;
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
	if (!page.params.userid) {
		return;
	}

	loadingMore = true;
	const searchParams = page.url.searchParams;
	searchParams.set('userid', page.params.userid);
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
	<title>SteamInputDB | UserProfile {playerInfo?.personaname ?? page.params.userid}</title>
	<meta property="og:site_name" content="SteamInputDB" />
	<meta property="og:type" content="website" />
	<meta property="og:url" content={page.url.href} />
	<meta property="og:title" content="SteamInputDB | {playerInfo?.personaname ?? page.params.userid} " />
	<meta
		name="description"
		content="Search for Steam Input configurations from {playerInfo?.personaname ??
			page.params.userid}" />
	<meta
		property="og:description"
		content="Search for Steam Input configurations from {playerInfo?.personaname ??
			page.params.userid}" />
	<meta name="twitter:card" content="summary_large_image" />
	<meta
		name="twitter:title"
		content="SteamInputDB | UserProfile {playerInfo?.personaname ?? page.params.userid} " />
	<meta
		name="twitter:description"
		content="Search for Steam Input configurations from {playerInfo?.personaname ??
			page.params.userid}" />
	{#if playerInfo?.avatarfull}
		<meta property="og:image" content={playerInfo.avatarfull} />
		<meta
			property="og:image:alt"
			content="SteamInputDB | UserProfile {playerInfo?.personaname ?? page.params.userid} " />
		<meta name="twitter:image" content={playerInfo.avatarfull} />
		<meta
			name="twitter:image:alt"
			content="SteamInputDB | UserProfile {playerInfo?.personaname ?? page.params.userid} " />
	{/if}
	<svelte:element this={'script'} type="application/ld+json">
		{createUserSchemaJsonLd({
			userId: page.params.userid ?? '',
			personName: playerInfo?.personaname ?? page.params.userid ?? '',
			imageUrl: playerInfo?.avatarfull
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
		<section class="user-header">
			<div>
				{#if playerInfo?.mini_profile_background?.movie_webm}
					<video autoplay loop muted playsinline>
						<source
							src={`${communityUrlBase.replace('/apps/', '/')}${playerInfo.mini_profile_background.movie_webm}`}
							type="video/webm" />
					</video>
				{/if}
				{#if playerInfo?.avatarfull}
					<enhanced:img src={playerInfo?.avatarfull}></enhanced:img>
					{#if playerInfo?.avatarframe?.small}
						<enhanced:img
							src={`${communityUrlBase.replace('/apps/', '/')}${playerInfo.avatarframe.small}`}
							class="avatar-frame"></enhanced:img>
					{/if}
				{:else}
					<div></div>
				{/if}
				<h1>{playerInfo?.personaname ?? page.params.userid}</h1>
			</div>
			<div>
				<a
					href={playerInfo?.profileurl}
					class="button"
					target="_blank"
					rel="external"
					{@attach tooltip({
						content: 'View Steam profile',
						outDelay: 200,
						arrow: true,
						placement: 'bottom',
						autoPlacement: false,

						arrowFollowCursor: true
					})}>
					<Icon icon="mdi:steam" width="1.4em" height="1.4em" />
					<!-- <Icon icon="mdi:local-grocery-store" width="1.4em" height="1.4em" /> -->
				</a>
			</div>
		</section>
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
							<span>No layouts found</span>
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

	& .title {
		text-overflow: ellipsis;
		overflow: hidden;
		white-space: nowrap;
		font-size: 1.3em;
	}

	& .thumb {
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

	& .info {
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

section.user-header {
	display: grid;
	width: 100%;
	gap: 1em;
	grid-template-columns: repeat(auto-fit, minmax(min(100%, 32em), auto));
	max-width: calc(100dvw -2em);
	container-type: inline-size;

	& > :first-child {
		position: relative;
		display: grid;
		grid-template-columns: minmax(min-content, 10em) auto;
		align-items: center;
		width: 100%;
		height: fit-content;
		grid-column-gap: 1em;
		grid-row-gap: 0.25em;
		padding: 1em 1.6em;
		margin-right: auto;
		overflow: hidden;
		@container (width > 1200px) {
			max-width: 40cqw;
		}

		& video {
			position: absolute;
			inset: 0;
			height: 100%;
			width: 100%;
			object-fit: cover;
			object-position: center;
			z-index: -1;
			border-radius: 1em;
			box-shadow: 0 0.25em 0.5em black;
		}

		& picture,
		& img {
			object-fit: cover;
			object-position: center;
			overflow: hidden;
			width: fit-content;
			grid-row: 1 / span 1;
			grid-column: 1 / span 1;
			min-height: 56px;
			min-width: 56px;
		}
		& .avatar-frame {
			scale: 1.2;
		}

		& > i {
			grid-row: 2 / span 1;
			grid-column: 1 / span 2;
		}

		& > :last-child {
			margin-right: auto;
			filter: drop-shadow(2px 2px 2px black) drop-shadow(0px 0px 8px black);
			text-overflow: ellipsis;
			width: 100%;
			overflow: hidden;
		}
	}
	& > :last-child {
		display: flex;
		flex-flow: row wrap;
		align-items: center;
		justify-content: end;
		justify-self: end;
		gap: 1em;

		& > a {
			white-space: nowrap;
			display: grid;
			align-items: center;
			justify-content: center;
			gap: 0.5ch;
			font-weight: bold;
			background: linear-gradient(
				215deg,
				color-mix(in srgb, var(--card-color), transparent 35%) 0%,
				color-mix(in srgb, var(--card-color), transparent 60%) 70%
			);
			& > span {
				width: fit-content;
			}
		}
		& .button {
			&:hover,
			&:focus-visible {
				color: var(--text-color-dark);
				background-color: var(--color-primary);
			}
		}
		/* & .button:is(:first-child) {
			background-color: #1a9fff;
			&:hover,
			&:focus-visible {
				background-color: color-mix(in srgb, #1a9fff, var(--color-primary) 50%);
			}
		} */
	}
}
</style>
