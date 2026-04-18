<script lang="ts" module>
export { sectionHead };
</script>

<script lang="ts">
import { browser } from '$app/environment';

import { resolve } from '$app/paths';

import type { components } from '$lib/api/openapi';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { client } from '$lib/buddy-api/client';
import type { components as buddyComponents } from '$lib/buddy-api/openapi';
import { BuddyState } from '$lib/buddy-app/buddyState.svelte';
import { assetUrlBase, communityUrlBase } from '$lib/steamapi/const';
import Icon from '@iconify/svelte';
import { cubicOut } from 'svelte/easing';
import { fade, slide } from 'svelte/transition';
import IconDownload from '~icons/mdi/download';
import IconSteam from '~icons/mdi/steam';
import BuddyApplyButton from './BuddyApplyButton.svelte';
</script>

{#snippet sectionHead({
	fileInfo,
	appInfo,
	isMobileBrowser
}: {
	fileInfo: components['schemas']['ConfigItem'];
	appInfo?: components['schemas']['AppItem'];
	isMobileBrowser?: boolean;
})}
	<section class="cfg-head">
		<div>
			{#if appInfo?.assets}
				{@const assets = appInfo?.assets}
				{@const assetChosen =
					assets.library_hero ??
					assets.header ??
					assets.package_header ??
					assets.main_capsule ??
					assets.small_capsule ??
					assets.hero_capsule ??
					assets.library_hero ??
					'none.svg'}
				{#if assetChosen}
					<picture class="capsule" transition:fade={{ duration: 196, easing: cubicOut }}>
						<enhanced:img
							src={`${assetUrlBase}${assets.asset_url_format?.replace(
								'${FILENAME}',
								assetChosen
							)}`}
							alt="Capsule"
							height="100%"></enhanced:img>
					</picture>
				{/if}
			{/if}
			<div>
				<h1>{fileInfo.title}</h1>
				{#if appInfo}
					{#if appInfo.assets?.community_icon}
						<enhanced:img
							src={`${communityUrlBase}${appInfo.app_id}/${appInfo.assets.community_icon}.jpg`}
							alt="Icon"
							style="min-width: 1.2em; height: 1.2em; margin-right: 0.1em;"></enhanced:img>
					{:else}
						<Icon icon="mdi:steam" width="1.2em" />
					{/if}
				{:else}
					<Icon icon="mdi:link-variant" width="1.2em" />
				{/if}
				<h2>
					<a href={resolve(`/app/${appInfo?.app_id || fileInfo.app_id_string}`)}>
						{appInfo?.name || fileInfo.app_id_string}
					</a>
				</h2>
			</div>
		</div>
		<div>
			{#if !isMobileBrowser}
				{#snippet tooltipContent()}
					<div
						style="display: grid; place-items: center;"
						in:slide={{ duration: 196, easing: cubicOut }}
						out:fade={{ duration: 196, easing: cubicOut }}>
						<p style="white-space: nowrap; text-align: center;">Preview this config in Steam</p>
						{#if !appInfo}
							<p style="text-align: center;">
								You must have a shortcut in Steam with the exact name "<em
									style="font-weight: bold;">{fileInfo.app_id_string}</em
								>"
							</p>
						{:else}
							<b>You must own the game</b>
						{/if}
						<em>Please note that Steam often bugs out when using this feature...</em>
						<p>In the worst case, you must restart Steam</p>
						<br />
						<strong
							>Alternatively, you should consider installing the
							<a href="https://steaminputdb.com/buddy-app/install">SteamInputDB Buddy App</a>
							for better and direct Steam integration</strong>
						<br />
						<code
							>steam://controllerconfig/{encodeURI(
								fileInfo.app_id_string ?? ''
							)}/{fileInfo.file_id}</code>
					</div>
				{/snippet}
				{#snippet defaultPreviewLinkButton()}
					<a
						href={`steam://controllerconfig/${encodeURI(fileInfo.app_id_string ?? '')}/${fileInfo.file_id}`}
						class="button blue"
						{@attach tooltip({
							snippet: tooltipContent,
							snippetInDefaultBackground: true,
							outDelay: 200,
							arrow: true,
							arrowFollowCursor: true
						})}>
						<IconSteam style="width: 1.4em; height: 1.4em;" />
						<span>Preview | Apply</span>
					</a>
				{/snippet}
				{#if browser && document.cookie?.includes('buddy-app=enabled') && BuddyState.reachable}
					<svelte:boundary pending={defaultPreviewLinkButton} failed={defaultPreviewLinkButton}>
						{@const controllers = (
							await client.GET('/v1/steam/controllers').then((r) => {
								if (r.error) {
									throw r.error;
								}
								return r;
							})
						).data as buddyComponents['schemas']['ControllerResponse'][]}
						{@const apps = (
							await client.GET('/v1/steam/apps').then((r) => {
								if (r.error) {
									throw r.error;
								}
								return r;
							})
						).data as buddyComponents['schemas']['AppResponse'][]}
						<BuddyApplyButton
							fileInfo={fileInfo}
							appInfo={appInfo}
							controllers={controllers}
							apps={apps} />
					</svelte:boundary>
				{:else}
					{@render defaultPreviewLinkButton()}
				{/if}
			{/if}

			{#if fileInfo.file_url && !(browser && navigator?.userAgent?.includes('Steam Gamepad'))}
				<a href={fileInfo.file_url} class="button" rel="external">
					<IconDownload style="width: 1.4em; height: 1.4em;" />
					<span>Download .vdf</span>
				</a>
			{/if}
		</div>
	</section>
{/snippet}

<style lang="postcss">
:global(section.cfg-head) {
	display: grid;
	width: 100%;
	gap: 1em;
	grid-template-columns: repeat(auto-fit, minmax(min(100%, 25ch), auto));
	padding: 1em;
	container-type: inline-size;

	& > :first-child {
		position: relative;
		display: grid;
		align-items: center;
		width: 100%;
		height: fit-content;
		min-height: 12em;
		grid-column-gap: 1em;
		grid-row-gap: 0.25em;
		padding: 1em 1.6em;
		margin-right: auto;
		color: var(--text-color-dark);

		& .capsule {
			position: absolute;
			inset: 0;
			height: 100%;
			width: 100%;
			object-fit: cover;
			object-position: center;
			z-index: -1;
			border-radius: 1em;
			overflow: hidden;
			box-shadow: 0 0.25em 0.5em black;

			object-fit: cover;
			object-position: center;
			width: 100%;
			box-shadow: 0 0.2em 0.7em 0em var(--shadow-color);
			& :global(img) {
				width: 100%;
				height: 100%;
				object-fit: cover;
				object-position: center;
			}
		}

		& > :nth-child(2) {
			margin-right: auto;
			display: grid;
			height: fit-content;
			grid-template-columns: min-content auto;
			place-items: center;
			gap: 0.5ch;
			filter: drop-shadow(2px 2px 2px black) drop-shadow(0px 0px 8px black)
				drop-shadow(0px 0px 24px black);

			& > :first-child {
				grid-column: 1 / span 2;
				text-align: start;
				width: 100%;
				color: var(--text-color-dark);
			}
			& :global(> :nth-child(1n + 2)) {
				color: color-mix(in srgb, var(--color-primary), var(--text-color-dark) 60%);
				font-size: 1.8em;
			}
			& > :last-child {
				margin-right: auto;
			}
			& a {
				color: color-mix(in srgb, var(--color-primary), var(--text-color-dark) 15%);
			}
		}
	}
	& > :last-child {
		display: grid;
		place-items: center;
		margin: auto;
		gap: 1em;
		width: 100%;
		grid-template-columns: repeat(auto-fit, minmax(19ch, auto));
		margin-right: 0;
		@container (width > 1200px) {
			max-width: 40cqw;
		}
		& > :global(.hov-over) {
			width: 100%;
		}

		& > a {
			width: 100%;
			white-space: nowrap;
			display: grid;
			align-items: center;
			justify-content: center;
			gap: 0.5ch;
			font-weight: bold;
			border: 10px solid transparent;
			border: none !important;

			& > span {
				width: fit-content;
			}
		}
		& .button {
			background:
				linear-gradient(
					215deg,
					color-mix(in srgb, var(--card-color), transparent 75%) 0%,
					color-mix(in srgb, var(--card-color), transparent 90%) 70%
				),
				var(--bg-noise-transparent);
			backdrop-filter: blur(6px);
			&:hover,
			&:focus-visible {
				color: var(--text-color-dark);
				background-color: var(--color-primary);
			}
		}
		& .button:is(.blue) {
			background-color: #1a9fff;
			&:hover,
			&:focus-visible {
				background-color: color-mix(in srgb, #1a9fff, var(--color-primary) 50%);
			}
		}
	}
	code {
		user-select: all;
		margin-top: 0.5em;
	}
}
</style>
