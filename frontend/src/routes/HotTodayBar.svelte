<script lang="ts">
import { resolve } from '$app/paths';
import { client } from '$lib/api/client';
import type { components } from '$lib/api/openapi';
import { log } from '$lib/log';
import { controllertype } from '$lib/snippets/controllertype.svelte';
import { assetUrlBase } from '$lib/steamapi/const';
import { cubicInOut, cubicOut } from 'svelte/easing';
import { fade } from 'svelte/transition';
import IconLink from '~icons/mdi/link-variant';
import IconSteam from '~icons/mdi/steam';

type Config = components['schemas']['ConfigItem'];

const {
	hotConfigs
}: {
	hotConfigs: Exclude<components['schemas']['ConfigsResponse']['items'], null | undefined>;
} = $props();

let infoAppIdMap = $state<Record<number, components['schemas']['AppItem']>>({});

$effect(() => {
	hotConfigs?.forEach((cfg) => {
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
</script>

<div id="hot-today">
	<div>
		{#each hotConfigs as config (config.file_id)}
			{#if infoAppIdMap?.[config.app_id || 0]?.name}
				{@render entry(
					config,
					config.title,
					`/config/${config.file_id}`,
					config.app_id,
					config.app_id_string
				)}
			{/if}
		{/each}
	</div>
</div>

{#snippet entry(
	e: Config,
	title: string | null,
	link_suffix: string,
	app_id?: number | null,
	app_id_string?: string | null
)}
	<a
		class="plain card"
		href={resolve(link_suffix as '/')}
		transition:fade|global={{ duration: 196, easing: cubicInOut }}>
		<div class="thumb">
			{#if infoAppIdMap?.[app_id || 0]?.assets}
				{@const assets = infoAppIdMap[app_id || 0]!.assets!}
				{@const srcChosen = assets?.asset_url_format
					? `${assetUrlBase}${assets?.asset_url_format?.replace(
							'${FILENAME}',
							assets.library_capsule_2x || assets.library_capsule || 'undefined'
						)}`
					: undefined}
				{#if srcChosen}
					<picture transition:fade={{ duration: 196, easing: cubicOut }}>
						<enhanced:img src={srcChosen} alt="Thumbnail" height="100%"></enhanced:img>
					</picture>
				{/if}
			{/if}
		</div>
		<div class="info">
			<span>
				{@render controllertype({ item: e })}
			</span>
			<i>
				{#if app_id}
					<IconSteam style="width: 1.2em; height: 1.2em;" />
				{:else}
					<IconLink style="width: 1.2em; height: 1.2em;" />
				{/if}
				{infoAppIdMap?.[app_id || 0]?.name ?? app_id_string}
			</i>
			<div></div>
			<div>
				<strong>{title}</strong>
			</div>
		</div>
	</a>
{/snippet}

<style lang="postcss">
#hot-today {
	display: grid;
	place-items: center;
	padding: 0.5em;
	& > div {
		display: grid;
		grid-auto-flow: column;
		grid-auto-columns: min(70dvw, 240px);
		gap: 1em;
		justify-content: center;
	}

	overflow-x: auto;
	width: 100%;
}

.thumb {
	aspect-ratio: 2 / 3;
	width: 100%;
	height: 100%;

	& picture,
	& img {
		object-fit: cover;
		object-position: center;
		width: 100%;
		height: 100%;
		overflow: hidden;
	}

	& picture {
		isolation: isolate;
		position: relative;
	}

	z-index: -1;
	isolation: isolate;
	display: grid;
	background: linear-gradient(rgb(128 128 1280 / 0.2), rgb(128 128 128 / 0.8));
}

a {
	overflow: clip;
	display: grid;
	grid-template-areas: 'stack';
	padding: 0;
	isolation: isolate;
	color: var(--text-color);
	color: var(--text-color-dark);

	&:hover,
	&:focus-visible {
		outline: 2px solid var(--color-primary);
		box-shadow: 0 0 1.3em -0.4em var(--color-primary);
	}

	& > * {
		grid-area: stack;
	}
}

.info {
	overflow: clip;
	overflow-clip-margin: 1em;
	width: 100%;
	height: 100%;
	display: grid;
	grid-template-rows: min-content min-content 1fr min-content;
	padding: 1em;
	font-size: 1.2em;
	position: relative;
	font-size: 1.1em;
	font-weight: 500;

	&::before {
		content: '';
		position: absolute;
		inset: -1em;
		background: linear-gradient(180deg, rgba(15, 1, 26, 0.897), transparent 50%, transparent);
		z-index: -1;
		mix-blend-mode: multiply;
	}

	gap: 0.5em;

	& > *:not(:last-child) {
		display: grid;
		grid-template-columns: auto auto;
		gap: 1ch;
		align-items: center;
		justify-content: start;
		width: 100%;
		filter: drop-shadow(1px 0.1em 4px rgb(0 0 0 /1)) drop-shadow(0 0 8px rgb(0 0 0 / 1));
	}

	& span :global(svg) {
		font-size: 1.2em;
	}

	strong {
		text-overflow: ellipsis;
		overflow: hidden;
		width: 100%;
		align-self: end;
		isolation: isolate;
		position: relative;
		filter: drop-shadow(1px 0.1em 4px rgb(0 0 0 /1)) drop-shadow(0 0 8px rgb(0 0 0 / 1));
	}

	& > :last-child {
		position: relative;
		grid-column: 1 / -1;
		align-self: end;
		display: grid;
		&::before {
			content: '';
			position: absolute;
			inset: -1em;
			background: linear-gradient(
				180deg,
				transparent,
				rgba(10, 1, 20, 0.24) 10%,
				rgba(15, 1, 26, 0.758)
			);
			z-index: -1;
			mix-blend-mode: multiply;
		}
	}
}
</style>
