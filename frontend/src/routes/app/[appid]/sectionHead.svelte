<script lang="ts" module>
export { sectionHead };

const CONTROLLER_SUPPORT_LEVEL_FULL = 2;
const CONTROLLER_SUPPORT_LEVEL_PARTIAL = 1;
</script>

<script lang="ts">
import type { components } from '$lib/api/openapi';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { assetUrlBase, communityUrlBase, steamStoreUrlBase } from '$lib/steamapi/const';
import Icon from '@iconify/svelte';
import { cubicOut } from 'svelte/easing';
import { fade } from 'svelte/transition';
import IcoDesktop from '~icons/mdi/monitor';

import IcoFullController from '$lib/assets/steam_controller_type_svgs/xbox.svg?component';
import IcoPartialController from '$lib/assets/steam_controller_type_svgs/xbox_partial.svg?component';
import IcoForbidden from '~icons/mdi/do-not-disturb-alt';

import { resolve } from '$app/paths';
import IcoDs4Full from '$lib/assets/steam_controller_type_svgs/ps4.svg?component';
import IcoDs4Partial from '$lib/assets/steam_controller_type_svgs/ps4_partial.svg?component';
import IcoDs5Full from '$lib/assets/steam_controller_type_svgs/ps5.svg?component';
import IcoDs5Partial from '$lib/assets/steam_controller_type_svgs/ps5_partial.svg?component';
import IcoSIAPI from '$lib/assets/steam_controller_type_svgs/siapi.svg?component';
import { CONTROLLER_LIST } from '$lib/components/search/controllerlist.svelte';
import IcoGeneric from '~icons/mdi/controller';

import IcoSteam from '~icons/mdi/steam';
import IcoPCGW from '~icons/simple-icons/pcgamingwiki';
import IcoProtonDB from '~icons/simple-icons/protondb';
import IcoSteamDB from '~icons/simple-icons/steamdb';
</script>

{#snippet sectionHead({
	appInfo,
	fallbackName
}: {
	appInfo?: components['schemas']['AppInfoItem'];
	fallbackName?: string;
})}
	<section class="app-header">
		<div>
			{#if appInfo?.assets?.library_hero || appInfo?.assets?.header || appInfo?.assets?.package_header || appInfo?.assets?.main_capsule}
				{@const srcChosen = appInfo?.assets?.asset_url_format
					? `${assetUrlBase}${appInfo.assets?.asset_url_format?.replace(
							'${FILENAME}',
							appInfo?.assets?.library_hero ??
								appInfo.assets.header ??
								appInfo.assets.package_header ??
								appInfo.assets.main_capsule ??
								'undefined'
						)}`
					: undefined}
				{#if srcChosen}
					<picture class="capsule" transition:fade={{ duration: 196, easing: cubicOut }}>
						<enhanced:img src={srcChosen} alt="Thumbnail" height="100%"></enhanced:img>
					</picture>
				{/if}
			{/if}
			<div>
				{#if appInfo?.assets?.community_icon}
					<picture transition:fade={{ duration: 196, easing: cubicOut }}>
						<enhanced:img
							src={`${communityUrlBase}${appInfo.app_id}/${appInfo.assets?.community_icon}.jpg`}
							alt="Icon"></enhanced:img>
					</picture>
				{:else}
					<!-- KEEP! -->
					{#if appInfo?.app_id !== 413080}
						<Icon icon="mdi:link-variant" width="2.5em" height="2.5em" />
					{:else}
						<IcoDesktop />
					{/if}
				{/if}
				{#if !appInfo && fallbackName}
					<i>(Non Steam Shortcut)</i>
				{/if}
				<h1>{appInfo?.name ?? fallbackName}</h1>
				<div>
					{#if appInfo?.controller_support?.support_level === CONTROLLER_SUPPORT_LEVEL_FULL}
						<div
							{@attach tooltip({
								content: 'Full controller support',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}>
							<IcoFullController width="2em" />
						</div>
					{:else if appInfo?.controller_support?.support_level === CONTROLLER_SUPPORT_LEVEL_PARTIAL}
						<div
							{@attach tooltip({
								content: 'Partial controller support',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}>
							<IcoPartialController width="2em" />
						</div>
						<!-- HACK if is real steam game -->
					{:else if appInfo?.assets?.community_icon}
						<div
							{@attach tooltip({
								content: 'No controller support',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}
							class="stacked">
							<IcoFullController width="2em" opacity="0.5" />
							<IcoForbidden width="2em" height="2em" color="red" style="z-index: 1;" />
						</div>
					{/if}
					{#if appInfo?.controller_support?.steam_input_api_support}
						<div
							{@attach tooltip({
								content: 'Steam Input API support',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}>
							<IcoSIAPI width="2em" />
						</div>
					{/if}
					{#if appInfo?.controller_support?.ds4_wired_support}
						<div
							{@attach tooltip({
								content: 'Native DualShock Controller support (USB only)',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}>
							<IcoDs4Partial width="2em" />
						</div>
					{/if}
					{#if appInfo?.controller_support?.ds4_wireless_support}
						<div
							{@attach tooltip({
								content: 'Native DualShock Controller support',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}>
							<IcoDs4Full width="2em" />
						</div>
					{/if}
					{#if appInfo?.controller_support?.ds5_wired_support}
						<div
							{@attach tooltip({
								content: 'Native DualSense Controller support (USB only)',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}>
							<IcoDs5Partial width="2em" />
						</div>
					{/if}
					{#if appInfo?.controller_support?.ds5_wireless_support}
						<div
							{@attach tooltip({
								content: 'Native DualSense Controller support',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}>
							<IcoDs5Full width="2em" />
						</div>
					{/if}
				</div>
			</div>
		</div>
		<div>
			{#if Object.entries(appInfo?.official_configs ?? {}).length}
				<div class="official-configs">
					<h3>Official Configs</h3>
					{#each Object.entries(appInfo?.official_configs ?? {}) as [controller_type, config_id] (config_id)}
						{@const controller_list_entry = CONTROLLER_LIST.find(
							(controller) => controller.type === controller_type
						)}
						<a href={resolve(`/config/${config_id}`)} class="button">
							{#if controller_list_entry}
								<controller_list_entry.icon width="2em" height="2em" />
							{:else}
								<IcoGeneric style="width: 2em; height: 2em;" />
							{/if}
							<span>{controller_list_entry?.niceName ?? 'Generic'}</span>
						</a>
					{/each}
				</div>
			{/if}
			{#if appInfo?.store_url_path}
				<a
					href={steamStoreUrlBase + appInfo?.store_url_path}
					class="button"
					target="_blank"
					rel="external"
					{@attach tooltip({
						content: 'View Steam store page',
						outDelay: 200,
						arrow: true,
						placement: 'bottom',
						autoPlacement: false,

						arrowFollowCursor: true
					})}>
					<IcoSteam style="width: 1.4em; height: 1.4em;" />
					<!-- <Icon icon="mdi:local-grocery-store" width="1.4em" height="1.4em" /> -->
				</a>
				<a
					href={`https://steamdb.info/app/${appInfo.app_id}/`}
					class="button"
					target="_blank"
					rel="external"
					{@attach tooltip({
						content: 'View on SteamDB',
						outDelay: 200,
						arrow: true,
						placement: 'bottom',
						autoPlacement: false,
						arrowFollowCursor: true
					})}>
					<IcoSteamDB style="width: 1.4em; height: 1.4em;" />
				</a>
				<a
					href={`https://www.protondb.com/app/${appInfo.app_id}`}
					class="button"
					target="_blank"
					rel="external"
					{@attach tooltip({
						content: 'View on ProtonDB',
						outDelay: 200,
						arrow: true,
						placement: 'bottom',
						autoPlacement: false,

						arrowFollowCursor: true
					})}>
					<IcoProtonDB style="width: 1.4em; height: 1.4em;" />
				</a>
				<a
					href={`https://www.pcgamingwiki.com/api/appid.php?appid=${appInfo.app_id}`}
					class="button"
					target="_blank"
					rel="external"
					{@attach tooltip({
						content: 'View on PCGaming Wiki',
						outDelay: 200,
						arrow: true,
						placement: 'bottom',
						autoPlacement: false,

						arrowFollowCursor: true
					})}>
					<IcoPCGW style="width: 1.4em; height: 1.4em;" />
				</a>
			{/if}
		</div>
	</section>
{/snippet}

<style lang="postcss">
:global(section.app-header) {
	display: grid;
	width: 100%;
	gap: 1em;
	grid-template-columns: repeat(auto-fit, minmax(min(100%, 32em), auto));
	max-width: calc(100dvw - 2em);
	container-type: inline-size;
	justify-self: center;

	.stacked {
		display: grid;
		grid-template-areas: 'stack';
		grid-template-columns: 1fr;
		grid-template-rows: 1fr;
		place-items: center;
		& > :global(*) {
			grid-area: stack;
		}
	}
	& > :first-child {
		position: relative;
		display: grid;
		align-items: center;
		width: 100%;
		height: 100%;
		min-height: 12em;
		padding: 1em 1.6em;
		color: var(--text-color-dark);
		@container (width > 1200px) {
			max-width: 40cqw;
		}

		& .capsule {
			position: absolute;
			inset: 0;
			height: 100%;
			width: 100%;
			object-fit: cover;
			object-position: center;
			z-index: -1;
			border-radius: 1em;
			box-shadow: 0 0.25em 0.5em black;

			object-fit: cover;
			object-position: center;
			width: 100%;
			box-shadow: 0 0.2em 0.7em 0em var(--shadow-color);

			& img {
				width: 100%;
				height: 100%;
			}
		}

		& picture,
		& img {
			object-fit: cover;
			object-position: center;
			overflow: hidden;
			width: fit-content;
		}

		& > :last-child {
			display: grid;
			grid-template-columns: minmax(min-content, 2em) auto;
			width: 100%;
			margin-right: auto;
			align-items: center;
			grid-column-gap: 1em;
			grid-row-gap: 0.25em;

			filter: drop-shadow(2px 2px 2px black) drop-shadow(0px 0px 8px black);

			& > i {
				grid-row: 2 / span 1;
				grid-column: 1 / span 2;
			}

			& h1 {
				text-overflow: ellipsis;
				width: 100%;
				overflow: hidden;
				margin-right: auto;
			}

			& > :last-child {
				grid-column: 1 / -1;
				display: flex;
				flex-direction: row wrap;
				gap: 1ch;
				align-items: center;
				/* background: rebeccapurple; */
			}
		}
	}
	& > :last-child {
		display: flex;
		flex-flow: row wrap-reverse;
		align-items: center;
		justify-content: end;
		width: 100%;
		margin: auto;
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
		}
		& .button {
			&:hover,
			&:focus-visible {
				color: var(--text-color-dark);
				background-color: var(--color-primary);
			}
		}
	}

	.official-configs {
		display: flex;
		flex-flow: row wrap;
		gap: 1ch;
		& > :first-child {
			width: 100%;
			flex: 1 0 100%;
		}
		& > a {
			display: grid;
			place-items: center;
			font-weight: bold;
			display: grid;
			align-items: center;
			justify-content: center;
			padding: 0.5em 1em;
			gap: 0.5ch;
			background: linear-gradient(
				215deg,
				color-mix(in srgb, var(--card-color), transparent 35%) 0%,
				color-mix(in srgb, var(--card-color), transparent 60%) 70%
			);
			& > :global(svg) {
				width: 2em;
				height: 2em;
			}
		}
	}
}
</style>
