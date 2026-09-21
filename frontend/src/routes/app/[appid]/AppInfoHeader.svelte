<script lang="ts" module>
const CONTROLLER_SUPPORT_LEVEL_FULL = 2;
const CONTROLLER_SUPPORT_LEVEL_PARTIAL = 1;
</script>

<script lang="ts">
import type { components } from '$lib/api/openapi';
import { ControllerSupportRating } from '$lib/api/steaminputdbEnums';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { assetUrlBase, communityUrlBase, steamStoreUrlBase } from '$lib/steamapi/const';
import Icon from '@iconify/svelte';
import { cubicOut } from 'svelte/easing';
import { fade } from 'svelte/transition';
import IcoDesktop from '~icons/mdi/monitor';

import RatingBronze from '$lib/assets/ratings/RatingBronze.svg?component';
import RatingGold from '$lib/assets/ratings/RatingGold.svg?component';
import RatingNoSupport from '$lib/assets/ratings/RatingNoSupport.svg?component';
import RatingPlatinum from '$lib/assets/ratings/RatingPlatinum.svg?component';
import RatingSilver from '$lib/assets/ratings/RatingSilver.svg?component';
import RatingUnknown from '$lib/assets/ratings/RatingUnknown.svg?component';
import IcoFullController from '$lib/assets/steam_controller_type_svgs/xbox.svg?component';
import IcoPartialController from '$lib/assets/steam_controller_type_svgs/xbox_partial.svg?component';
import IconHelp from '~icons/material-symbols/help-outline';

import IcoDs4Full from '$lib/assets/steam_controller_type_svgs/ps4.svg?component';
import IcoDs4Partial from '$lib/assets/steam_controller_type_svgs/ps4_partial.svg?component';
import IcoDs5Full from '$lib/assets/steam_controller_type_svgs/ps5.svg?component';
import IcoDs5Partial from '$lib/assets/steam_controller_type_svgs/ps5_partial.svg?component';
import IcoSIAPI from '$lib/assets/steam_controller_type_svgs/siapi.svg?component';

import IcoSteam from '~icons/mdi/steam';
import IcoPCGW from '~icons/simple-icons/pcgamingwiki';
import IcoProtonDB from '~icons/simple-icons/protondb';
import IcoSteamDB from '~icons/simple-icons/steamdb';

import IcoCross from '~icons/mdi/close';
import IcoPencil from '~icons/mdi/pencil';

let {
	appInfo,
	fallbackName,
	isAllowedEditInfo,
	editControllerSupport = $bindable<boolean>(false)
}: {
	appInfo?: components['schemas']['AppInfoItem'];
	fallbackName?: string;
	isAllowedEditInfo?: boolean;
	editControllerSupport?: boolean;
} = $props();
</script>

<section class="app-header">
	<div class="app-lead">
		<div class="app-banner">
			{#if appInfo?.assets?.library_hero || appInfo?.assets?.package_header || appInfo?.assets?.main_capsule || appInfo?.assets?.header}
				{@const srcChosen = appInfo?.assets?.asset_url_format
					? `${assetUrlBase}${appInfo.assets?.asset_url_format?.replace(
							'${FILENAME}',
							appInfo?.assets?.library_hero ??
								appInfo.assets.package_header ??
								appInfo.assets.main_capsule ??
								appInfo.assets.header ??
								'undefined'
						)}`
					: undefined}
				{#if srcChosen}
					<picture class="capsule" transition:fade={{ duration: 196, easing: cubicOut }}>
						<enhanced:img src={srcChosen} alt="Thumbnail" height="100%"></enhanced:img>
					</picture>
				{/if}
			{/if}
			<div class="info">
				{#if appInfo?.assets?.community_icon}
					<picture transition:fade={{ duration: 196, easing: cubicOut }}>
						<enhanced:img
							src={`${communityUrlBase}${appInfo.app_id}/${appInfo.assets?.community_icon}.jpg`}
							alt="Icon"
						></enhanced:img>
					</picture>
				{:else}
					<!-- KEEP! -->
					{#if appInfo?.app_id == 413080 || appInfo?.app_id == 769}
						<IcoDesktop />
					{:else}
						<Icon icon="mdi:link-variant" width="2.5em" height="2.5em" />
					{/if}
				{/if}
				{#if !appInfo && fallbackName}
					<i>(Non Steam Shortcut)</i>
				{/if}
				<h1>{appInfo?.name ?? fallbackName}</h1>
				<div class="official-controller-support">
					{#if appInfo?.controller_support?.support_level === CONTROLLER_SUPPORT_LEVEL_FULL}
						<div
							{@attach tooltip({
								content: 'Full controller support',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}
						>
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
							})}
						>
							<IcoPartialController width="2em" />
						</div>
					{:else if !appInfo?.store_url_path}
						<!-- <div
								{@attach tooltip({
									content: 'Unknown controller support',
									outDelay: 200,
									arrow: true,
									placement: 'bottom',
									autoPlacement: false
								})}
								class="stacked"
							>
								<IcoFullController width="2em" opacity="0.5" />
								<IcoUnknown width="2em" height="2em" style="z-index: 1;" />
							</div> -->
					{/if}
					{#if appInfo?.controller_support?.steaminputapi_support}
						<div
							{@attach tooltip({
								content: 'Steam Input API support',
								outDelay: 200,
								arrow: true,
								placement: 'bottom',
								autoPlacement: false
							})}
						>
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
							})}
						>
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
							})}
						>
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
							})}
						>
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
							})}
						>
							<IcoDs5Full width="2em" />
						</div>
					{/if}
				</div>
			</div>

			{#if isAllowedEditInfo}
				{#if editControllerSupport}
					<button
						class="edit-info-button"
						type="button"
						onclick={() => (editControllerSupport = false)}
					>
						<IcoCross style="width: 1.2em; height: 1.2em; " />Cancel Edit
					</button>
				{:else}
					<button
						class="edit-info-button"
						type="button"
						onclick={() => (editControllerSupport = true)}
					>
						<IcoPencil style="width: 1.2em; height: 1.2em; " />Edit
					</button>
				{/if}
			{/if}
		</div>

		<div class="rating-container card glass">
			{#snippet controllerSupportRatingHelp()}
				<div
					style="display: grid; place-items: center; max-width: 80svw; filter: drop-shadow(0 0 0.3rem rgba(0, 0, 0, 0.5));"
				>
					<p style="text-align: center; font-size: 1.2em; font-weight: bold; margin-bottom: 0.5em;">
						Controller support rating
					</p>
					<p style="text-align: center;">
						The controller support rating system is a work and progress and to be determined
					</p>
					<p style="text-align: center;">
						Ratings will be based on a set of criteria including native controller support,
						SteamInputAPI, glyphs, and mixed input support.
					</p>
				</div>
			{/snippet}
			<div
				class="rating-help"
				{@attach tooltip({
					snippet: controllerSupportRatingHelp,
					snippetInDefaultBackground: true,
					autoPlacement: true,
					outDelay: 200,
					arrow: true,
					arrowFollowCursor: true
				})}
			>
				<IconHelp style="width: 1.6em; height: 1.6em; color: var(--text-muted);" />
			</div>
			{#if !appInfo?.store_url_path}
				<RatingUnknown />
				<div>
					<span>Not yet rated</span>
				</div>
			{:else}
				{#if appInfo?.steaminputdb_info?.controller_support_rating === ControllerSupportRating.Platinum}
					<RatingPlatinum />
					<div>
						<span style="color: #afd5e0;">Platinum</span>
					</div>
				{:else if appInfo?.steaminputdb_info?.controller_support_rating === ControllerSupportRating.Gold}
					<RatingGold />
					<div>
						<span style="color: #d4a74f;">Gold</span>
					</div>
				{:else if appInfo?.steaminputdb_info?.controller_support_rating === ControllerSupportRating.Silver}
					<RatingSilver />
					<div>
						<span style="color: silver;">Silver</span>
					</div>
				{:else if appInfo?.steaminputdb_info?.controller_support_rating === ControllerSupportRating.Bronze}
					<RatingBronze />
					<div>
						<span style="color: #b1865b;">Bronze</span>
					</div>
				{:else if appInfo?.steaminputdb_info?.controller_support_rating === ControllerSupportRating.Wood}
					<RatingNoSupport style="color: firebrick;" />
					<div>
						<span>No controller support</span>
					</div>
				{:else}
					<RatingUnknown />
					<div>
						<span>Not yet rated</span>
					</div>
				{/if}
			{/if}
		</div>
	</div>

	<div class="link-container">
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
				})}
			>
				<IcoSteam style="width: 1.4em; height: 1.4em;" />
				<!-- <Icon icon="mdi:local-grocery-store" width="1.4em" height="1.4em" /> -->
			</a>
		{/if}
		{#if appInfo?.app_id && typeof appInfo?.app_id === 'number'}
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
				})}
			>
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
				})}
			>
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
				})}
			>
				<IcoPCGW style="width: 1.4em; height: 1.4em;" />
			</a>
		{/if}
	</div>
</section>

<style lang="postcss">
section.app-header {
	display: flex;
	flex-flow: row wrap;
	width: 100%;
	gap: 1em;
	max-width: calc(100svw - 2em);
	container-type: inline-size;
	justify-self: center;
	justify-content: center;

	.app-lead {
		display: flex;
		flex-flow: row wrap;
		gap: 1em;
		margin-right: auto;
		flex-grow: 1;
		width: 100%;
	}

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

	.rating-container {
		position: relative;
		display: grid;
		place-items: center;
		color: var(--text-color);
		min-height: 192px;
		min-width: 222px;
		height: fit-content;
		& > :global(svg) {
			max-height: 6.66em;
		}
		& span {
			white-space: nowrap;
			font-size: 1.6em;
			font-weight: bold;
		}

		& > :last-child {
			display: grid;
			grid-auto-flow: column;
			align-items: center;
			gap: 0.5ch;
			color: var(--text-color-dark);
			& span {
				filter: drop-shadow(2px 2px 2px rgba(0, 0, 0, 0.842))
					drop-shadow(0px 0px 8px rgba(0, 0, 0, 0.671));
			}
		}
		.rating-help {
			position: absolute;
			top: 1em;
			right: 1em;
		}

		flex-grow: 1;
	}

	.app-banner {
		position: relative;
		display: grid;
		align-items: center;
		min-height: 12em;
		padding: 1em 1.6em;
		flex-grow: 9999999;
		min-width: min(100cqi, 32em);
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

		& > .info {
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

			& .official-controller-support {
				grid-column: 1 / -1;
				display: flex;
				flex-direction: row wrap;
				gap: 1ch;
				align-items: center;
				/* background: rebeccapurple; */
			}
		}
	}

	& .link-container {
		display: flex;
		flex-flow: row wrap-reverse;
		align-items: start;
		justify-content: center;
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
}

.edit-info-button {
	position: absolute;
	top: 1em;
	right: 1em;
	z-index: 10;
	font-weight: bold;
	display: flex;
	gap: 0.5em;
	&:hover,
	&:focus-visible {
		color: var(--text-color);
	}
}
</style>
