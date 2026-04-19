<script lang="ts" module>
export { sectionHead };
</script>

<script lang="ts">
import type { components } from '$lib/api/openapi';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { assetUrlBase, communityUrlBase, steamStoreUrlBase } from '$lib/steamapi/const';
import Icon from '@iconify/svelte';
import { cubicOut } from 'svelte/easing';
import { fade } from 'svelte/transition';
import IcoDesktop from '~icons/mdi/monitor';
</script>

{#snippet sectionHead({
	appInfo,
	fallbackName
}: {
	appInfo?: components['schemas']['AppItem'];
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
		</div>
		<div>
			<!-- 
            TODO: create buddy-app that interacts with steam via cef-remote-debug
            If you are reading this and think this works without - Nope CORS policy. and that's a good thing!
			<a href="#" class="button">
				<Icon icon="mdi:steam" width="1.4em" height="1.4em" />
				<span>Show Controller Config</span>
			</a> -->
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
					<Icon icon="mdi:steam" width="1.4em" height="1.4em" />
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
					<Icon icon="simple-icons:steamdb" width="1.4em" height="1.4em" />
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
					<Icon icon="simple-icons:protondb" width="1.4em" height="1.4em" />
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
	max-width: calc(100dvw -2em);
	container-type: inline-size;
	justify-self: center;

	& > :first-child {
		position: relative;
		display: grid;
		grid-template-columns: minmax(min-content, 2em) auto;
		align-items: center;
		width: 100%;
		height: fit-content;
		min-height: 12em;
		grid-column-gap: 1em;
		grid-row-gap: 0.25em;
		padding: 1em 1.6em;
		margin-right: auto;
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
			box-shadow: 0 0.2em 0.7em 0em var(--shadow-color);
		}

		& > i {
			grid-row: 2 / span 1;
			grid-column: 1 / span 2;
		}

		& > :last-child {
			margin-right: auto;
			filter: drop-shadow(2px 2px 2px black) drop-shadow(0px 0px 8px black);
			& h1 {
				text-overflow: ellipsis;
				width: 100%;
				overflow: hidden;
			}
		}
	}
	& > :last-child {
		display: flex;
		flex-flow: row wrap;
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
