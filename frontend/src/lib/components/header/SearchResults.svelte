<script lang="ts">
import { resolve } from '$app/paths';
import { client, type ResponseType } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import SC2 from '$lib/assets/SC2_Alt.svg?component';
import { log } from '$lib/log';
import { controllertype } from '$lib/snippets/controllertype.svelte';
import { assetUrlBase } from '$lib/steamapi/const';
import Icon from '@iconify/svelte';
import { cubicInOut, cubicOut } from 'svelte/easing';
import { fade, slide } from 'svelte/transition';

type Config = components['schemas']['ConfigItem'];
type Game = components['schemas']['AppItem'];

let {
	results
}: {
	results: ResponseType<'POST', '/v1/search/'>['data'];
} = $props();

let resultAppIdMap = $derived.by(() => {
	let map = results?.games?.reduce(
		(acc, game) => {
			if (!game.app_id) {
				return acc;
			}
			acc[game.app_id] = game;
			return acc;
		},
		{} as Record<number, (typeof results.games)[0]>
	);
	return map;
});

let infoAppIdMap = $state<Record<number, components['schemas']['AppItem']>>({});

$effect(() => {
	results?.configs?.forEach((cfg) => {
		if (!cfg.app_id) {
			return;
		}
		if (resultAppIdMap?.[cfg.app_id] || infoAppIdMap[cfg.app_id]) {
			return;
		}
		const idCopy = cfg.app_id;
		infoAppIdMap[idCopy] = {
			app_id: idCopy
		} as components['schemas']['AppItem'];
		client
			.GET('/v1/steam/appinfo/{app_id}', {
				params: {
					path: {
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
</script>

<div class="results">
	{#each results?.games || [] as game (game.app_id)}
		{@render entry(game, 'game', game.name, `/app/${game.app_id}`, game.app_id, `${game.app_id}`)}
	{/each}
	{#each results?.configs || [] as config (config.file_id)}
		{@render entry(
			config,
			'config',
			config.title,
			`/config/${config.file_id}`,
			config.app_id,
			config.app_id_string
		)}
	{/each}
</div>

{#snippet entry(
	e: Config | Game,
	type: 'game' | 'config',
	title: string | null,
	link_suffix: string,
	app_id?: number | null,
	app_id_string?: string | null
)}
	<!-- TODO: fix types -->
	<a
		class="plain"
		href={resolve(link_suffix as '/')}
		transition:slide|global={{ duration: 196, easing: cubicInOut }}
	>
		<div class="thumb">
			{#if resultAppIdMap?.[app_id || 0]?.assets}
				{@const assets = resultAppIdMap[app_id || 0]!.assets!}
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
						<enhanced:img src={srcChosen} alt="Thumbnail" height="100%"></enhanced:img>
					</picture>
				{/if}
			{:else if infoAppIdMap?.[app_id || 0]}
				{@const assets = infoAppIdMap[app_id || 0]!.assets!}
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
						<enhanced:img src={srcChosen} alt="Thumbnail" height="100%"></enhanced:img>
					</picture>
				{/if}
			{/if}
			<div>
				{#if type === 'game'}
					<Icon icon="mdi:steam" width="100%" style="transform: scale(1.5);" />
				{:else}
					<SC2 />
				{/if}
				<span>{type.replace(/./, (c) => c.toUpperCase())}</span>
			</div>
		</div>
		<div class="info">
			<strong class={`title ${type}`}>{title}</strong>
			{#if type === 'config'}
				<i>
					{#if app_id}
						<Icon icon="mdi:steam" width="1.2em" />
					{:else}
						<Icon icon="mdi:link-variant" width="1.2em" />
					{/if}
					{resultAppIdMap?.[app_id || 0]?.name ??
						infoAppIdMap?.[app_id || 0]?.name ??
						app_id_string}
				</i>
				<span>
					{@render controllertype({ item: e as Config })}
				</span>
			{/if}
		</div>
	</a>
{/snippet}

<style lang="postcss">
.results {
	display: grid;
	isolation: isolate;
	& > :last-child {
		& .thumb {
			border-radius: 0 0 0 0.7em;
			& picture::before {
				border-radius: 0 0 0 0.7em;
			}
		}
		&::before {
			display: none;
		}
	}
	anchor-name: --hovered-link;
	a:hover,
	a:focus-visible {
		anchor-name: --hovered-link;
	}

	overflow: clip;
	border-radius: 0 0 1em 1em;

	--transition-delay: var(--transition-duration);

	&::after {
		content: '';
		position: absolute;

		top: anchor(top);
		left: anchor(left);
		right: anchor(right);
		bottom: anchor(bottom);

		opacity: 0;

		transition:
			top calc(var(--transition-duration) * 1) cubic-bezier(0.086, 1.037, 0.621, 0.903)
				var(--transition-delay),
			bottom calc(var(--transition-duration) * 1) cubic-bezier(0.086, 1.037, 0.621, 0.903)
				var(--transition-delay),
			opacity var(--transition-duration) var(--default-ease);

		position-anchor: --hovered-link;

		background: var(--color-primary);
		border-radius: var(--hover-radius);
		z-index: -2;
	}
	&:has(a:hover)::after,
	&:has(a:focus-visible)::after {
		top: anchor(top);
		left: anchor(left);
		right: anchor(right);
		bottom: anchor(bottom);
		opacity: 0.3;
		--transition-delay: 0ms;
	}
	&:has(a:hover:is(:last-child))::after,
	&:has(a:focus-visible:is(:last-child))::after {
		border-radius: 0 0 1em 1em;
	}
}

a {
	display: grid;
	--gap: 1em;
	padding: calc(0.5 * var(--gap)) calc(0.5 * var(--gap));
	column-gap: 1em;
	row-gap: 0;
	grid-template-columns: 33% 66%;
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
	&:is(.game) {
		font-size: 1.5em;
		align-self: center;
	}
	&:is(.config) {
		font-size: 1.2em;
	}
}

.thumb {
	aspect-ratio: 21 / 8;
	width: 100%;
	height: 100%;

	color: var(--text-color-dark);
	isolation: isolate;

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
			mask: linear-gradient(120deg, white 0%, white 20%, transparent 50%);
			mask-type: alpha;
			backdrop-filter: saturate(0.1) brightness(0.3);
			z-index: 1;
		}
	}

	display: grid;
	background: linear-gradient(rgb(128 128 1280 / 0.2), rgb(128 128 128 / 0.8));
	& > * {
		grid-row: 1;
		grid-column: 1;
	}

	& > :last-child {
		display: grid;
		place-items: center;
		max-width: 35%;
		padding: 0.5em;
		height: 100%;
		filter: drop-shadow(1px 0.1em 2px rgb(0 0 0 / 0.5)) drop-shadow(0 0px 8px rgb(0 0 0 / 0.5));
		& > :global(svg) {
			max-height: 50%;
			max-width: 100%;
			height: 100%;
		}
		& > :last-child {
			translate: 0 -0.5em;
		}
	}
}

.info {
	display: grid;
	grid-template-rows: min-content min-content auto;
	& span,
	i {
		display: grid;
		align-items: center;
		grid-auto-flow: column;
		width: fit-content;
		gap: 0.5em;
	}
	strong {
		padding-bottom: 0.5em;
	}

	& span {
		color: var(--highlight-color);
		font-weight: bold;
	}
}
</style>
