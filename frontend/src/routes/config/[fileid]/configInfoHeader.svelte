<script lang="ts">
import { browser } from '$app/environment';

import { resolve } from '$app/paths';

import type { components } from '$lib/api/openapi';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { client } from '$lib/buddy-api/client';
import type { components as buddyComponents } from '$lib/buddy-api/openapi';
import { BuddyState } from '$lib/buddy-app/buddyState.svelte';
import { CONTROLLER_LIST } from '$lib/components/search/controllerlist.svelte';
import { configRating } from '$lib/snippets/configRating.svelte';
import { assetUrlBase, communityUrlBase } from '$lib/steamapi/const';
import { format, formatDistance, formatDistanceToNow } from 'date-fns';
import { cubicOut } from 'svelte/easing';
import { fade, slide } from 'svelte/transition';
import IconDownload from '~icons/mdi/download';
import IcoLink from '~icons/mdi/link-variant';
import { default as IconSteam, default as IcoSteam } from '~icons/mdi/steam';
import BuddyApplyButton from './BuddyApplyButton.svelte';

let {
	fileInfo,
	appInfo,
	isMobileBrowser,
	creatorInfo
}: {
	fileInfo: components['schemas']['ConfigDetailResponse'];
	appInfo?: components['schemas']['AppItem'];
	isMobileBrowser?: boolean;
	creatorInfo?: components['schemas']['UserInfoResponse'];
} = $props();

const controller = $derived.by(() => {
	// eslint-disable-next-line prettier/prettier
	return CONTROLLER_LIST.find(
        (c) => c.type == fileInfo.controller_type
    ) || CONTROLLER_LIST[0];
});
$inspect(controller);
</script>

<section class="cfg-head">
	<div class="config-lead">
		<div class="config-banner">
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
							src={`${assetUrlBase}${assets.asset_url_format?.replace('${FILENAME}', assetChosen)}`}
							alt="Capsule"
							height="100%"
						></enhanced:img>
					</picture>
				{/if}
			{/if}
			<div class="info-grid">
				<h1>{fileInfo.title}</h1>

				<span>by</span>
				<div>
					{#if creatorInfo?.steamid}
						<a href={resolve(`/user/${creatorInfo?.steamid}`)}>
							{creatorInfo?.personaname || 'Unknown'}
						</a>
					{:else}
						<span>{creatorInfo?.personaname || 'Unknown'}</span>
					{/if}
				</div>
				<controller.icon style="width: 1.4em; height: 1.4em;" />
				<span>{fileInfo.controller_type_nice}</span>
				{#if appInfo}
					{#if appInfo.assets?.community_icon}
						<enhanced:img
							src={`${communityUrlBase}${appInfo.app_id}/${appInfo.assets.community_icon}.jpg`}
							alt="Icon"
							style="min-width: 1.2em; height: 1.2em; margin-right: 0.1em;"
						></enhanced:img>
					{:else}
						<IcoSteam style="width: 1.2em; height: 1.2em;" />
					{/if}
				{:else}
					<IcoLink style="width: 1.2em; height: 1.2em;" />
				{/if}
				<h2>
					<a href={resolve(`/app/${appInfo?.app_id || fileInfo.app_id_string}`)}>
						{appInfo?.name || fileInfo.app_id_string}
					</a>
				</h2>
				<div class="created-date">
					<span>{format(new Date(fileInfo.time_created), 'PPpp')}</span>
					<i>({formatDistanceToNow(new Date(fileInfo.time_created))} ago)</i>
				</div>
			</div>
		</div>
		<div class="rating-container">
			{#if fileInfo.votes}
				<div>
					<div class="cfgrating">
						{@render configRating({ item: fileInfo })}
					</div>
				</div>
			{/if}
			{#if fileInfo.playtime_seconds || fileInfo.lifetime_playtime_seconds}
				<div>
					<div class="playtime">
						{#if fileInfo.lifetime_playtime_sessions}
							<span
								>{formatDistance(
									new Date(fileInfo.lifetime_playtime_seconds * 1000),
									new Date(0)
								)}</span
							>
							<span>combined playtime in</span>
							<span>{(fileInfo.lifetime_playtime_sessions ?? 0).toLocaleString()}</span>
							<span>sessions</span>
							<i>(all users - since upload)</i>
							{#if fileInfo.playtime_seconds}
								<span class="last14-head">Last 14 days:</span>
								<span class="last14">
									{formatDistance(
										new Date((fileInfo.playtime_seconds ?? 0) * 1000),
										new Date(0)
									)} in {(fileInfo.playtime_sessions ?? 0).toLocaleString()} sessions
								</span>
							{/if}
						{:else if fileInfo.playtime_seconds}
							<span
								>{formatDistance(
									new Date((fileInfo.playtime_seconds ?? 0) * 1000),
									new Date(0)
								)}</span
							>
							<span>combined playtime in</span>
							<span>{(fileInfo.playtime_sessions ?? 0).toLocaleString()}</span>
							<span>sessions</span>
							<i>(all users - last 14 days)</i>
						{/if}
					</div>
				</div>
			{/if}
		</div>
	</div>
	<div>
		{#if !isMobileBrowser}
			{#snippet tooltipContent()}
				<div
					style="display: grid; place-items: center;"
					in:slide={{ duration: 196, easing: cubicOut }}
					out:fade={{ duration: 196, easing: cubicOut }}
				>
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
						for better and direct Steam integration</strong
					>
					<br />
					<code
						>steam://controllerconfig/{encodeURI(
							fileInfo.app_id_string ?? ''
						)}/{fileInfo.file_id}</code
					>
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
					})}
				>
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
						apps={apps}
					/>
				</svelte:boundary>
			{:else}
				{@render defaultPreviewLinkButton()}
			{/if}
		{/if}

		{#if fileInfo.file_url && !(browser && navigator?.userAgent?.includes('Steam Gamepad'))}
			{#snippet downloadTooltip()}
				<div
					style="display: grid; place-items: center;"
					in:slide={{ duration: 196, easing: cubicOut }}
					out:fade={{ duration: 196, easing: cubicOut }}
				>
					<b>File Size</b>
					<span>{(fileInfo.file_size / 1000).toFixed(0)} kB</span>
				</div>
			{/snippet}
			<a
				href={fileInfo.file_url}
				class="button"
				rel="external"
				{@attach tooltip({
					snippet: fileInfo.file_size ? downloadTooltip : undefined,
					snippetInDefaultBackground: true,
					outDelay: 200,
					arrow: true,
					autoPlacement: false,
					placement: 'bottom',
					arrowFollowCursor: false
				})}
			>
				<IconDownload style="width: 1.4em; height: 1.4em;" />
				<span>Download .vdf</span>
			</a>
		{/if}
	</div>
</section>

<style lang="postcss">
section.cfg-head {
	display: flex;
	flex-flow: row wrap;
	width: 100%;
	gap: 1em;
	container-type: inline-size;
	justify-self: center;
	justify-content: center;
	padding: 0 1em;

	& > :last-child {
		display: grid;
		place-items: center;
		margin: auto;
		gap: 1em;
		width: 100%;
		grid-template-columns: repeat(auto-fit, minmax(19ch, auto));
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

.config-lead {
	display: flex;
	flex-flow: row wrap;
	gap: 1em;
	margin-right: auto;
	flex-grow: 1;
	width: 100%;
}

.config-banner {
	position: relative;
	display: grid;
	align-items: center;
	height: fit-content;
	min-height: 12em;
	flex-grow: 9999999;

	grid-column-gap: 1em;
	grid-row-gap: 0.25em;
	margin-right: auto;
	color: var(--text-color-dark);
	border-radius: 1em;
	overflow: clip;

	& .capsule {
		position: absolute;
		inset: 0;
		height: 100%;
		width: 100%;
		object-fit: cover;
		object-position: center;
		z-index: -1;
		border-radius: inherit;
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

	& > .info-grid {
		margin-right: auto;
		display: grid;
		height: 100%;
		grid-template-columns: min-content 1fr;
		align-items: center;
		grid-column-gap: 1ch;
		grid-row-gap: 0;
		margin: 0;
		padding: 1em 1.6em;
		width: fit-content;

		background: linear-gradient(90deg, rgba(0, 0, 0, 0.623), rgb(0 0 0 / 0));

		& :global(> *) {
			filter: drop-shadow(2px 2px 2px black) drop-shadow(0px 0px 4px rgba(0, 0, 0, 0.801))
				drop-shadow(0px 0px 24px black);
		}
		& > :first-child {
			grid-column: 1 / span 2;
			text-align: start;
			width: 100%;
			color: var(--text-color-dark);
			text-overflow: ellipsis;
			white-space: nowrap;
			overflow: hidden;
		}
		& :global(> *) {
			font-size: 1.4em;
		}
		& h1 {
			font-size: 1.5em;
			margin-bottom: 0.125em;
		}
		& > :nth-child(2) {
			margin-left: 0.5em;
		}
		& > :nth-child(2),
		& > :nth-child(3) {
			font-size: 1.1em;
		}

		& > :nth-child(1n + 2) {
			margin-right: auto;
		}
		& > :nth-last-child(2),
		& :global(> :nth-last-child(3)) {
			margin-top: 1em;
		}
		& > :nth-child(-n + 3):not(:first-child) {
			margin-bottom: 0.5em;
		}
		/* & > :first-child {
			margin-bottom: 1em;
		} */
		& a {
			color: color-mix(in srgb, var(--color-primary), var(--text-color-dark) 36%);
			font-weight: 500;
			display: flex;
			align-items: center;
			gap: 0.5ch;
		}
		& h2 {
			font-size: 1.4em;
			& a {
				font-weight: bold;
				font-weight: 600;
				text-decoration: none;
			}
		}
	}
}

.rating-container {
	position: relative;
	display: flex;
	flex-flow: row wrap;
	align-items: stretch;
	color: var(--text-color);
	min-width: 222px;
	flex-grow: 1;
	gap: 1em;

	& > * {
		flex: 1 1 auto;
		padding: 1em 1.5em;
		display: grid;
		place-items: center;
		& :global(> :first-child) {
			height: 100%;
		}
	}

	.cfgrating {
		padding: 1.5rem;
	}

	:global(> *) {
		position: relative;
		isolation: isolate;
		border-radius: var(--border-radius);
		box-shadow: var(--card-shadow);
		backdrop-filter: blur(8px);

		&::before {
			content: '';
			position: absolute;
			inset: 0;
			background: var(--card-glass);
			opacity: 0.5;
			border-radius: var(--border-radius);
			z-index: -1;
		}
		&::after {
			content: '';
			position: absolute;
			inset: 0;
			border-radius: inherit;
			border: 1px solid transparent;
			background: var(--card-border-pseudo-gradient) border-box;
			mask:
				linear-gradient(black, black) border-box,
				linear-gradient(black, black) padding-box;
			mask-composite: subtract;
			z-index: -1;
		}
	}
}

.playtime {
	display: grid;
	place-items: center;
	height: 100%;

	& span {
		font-size: 1.2em;
		font-weight: 500;
		color: var(--text-color-dark);
		filter: drop-shadow(1px 1px 1px black);
	}

	& > :nth-child(2n) {
		font-size: 1em;
		font-weight: normal;
		text-align: center;
	}

	.last14-head {
		color: var(--text-color-dark);
		opacity: 0.9;
		font-size: 1em;
		filter: drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.897)) drop-shadow(0px 0px 2px rgba(0, 0, 0, 0.671));
		padding-top: 0.5em;
	}
	.last14 {
		opacity: 0.8;
		font-size: 0.9em;
		filter: drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.87)) drop-shadow(0px 0px 2px rgba(0, 0, 0, 0.527));
	}
	& i {
		color: var(--text-color-dark);
		opacity: 0.75;
		font-size: 0.75em;
		filter: drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.76)) drop-shadow(0px 0px 2px rgba(0, 0, 0, 0.651));
	}
	& :global(> div) {
		& > :first-child {
			font-size: 1.8em;
			filter: drop-shadow(1px 2px 3px black);
			transform: translate(0, 0.5em);
		}
		& > :last-child {
			font-size: 1.6em;
		}
	}
}

.created-date {
	grid-column: 1 / -1;
	display: grid;
	place-items: end;
	margin-left: 1.6em;
	margin-top: 0.5em !important;
	opacity: 0.9;
	& span {
		font-size: 0.7em;
		font-weight: 300;
	}
	& i {
		color: var(--text-color-dark);
		opacity: 0.8;
		font-size: 0.65em;
	}
}
</style>
