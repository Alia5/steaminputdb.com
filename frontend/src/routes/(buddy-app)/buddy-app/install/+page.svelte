<script lang="ts">
import { selectAllHandler } from '$lib/attachments/selectAllHandler.svelte';
import Spinner from '$lib/components/Spinner.svelte';
import { fade } from 'svelte/transition';
import '../../../../css/md.pcss';

import AboutBuddy from './about-buddy.svx';

import { browser } from '$app/environment';
import IcoAlert from '~icons/mdi/alert';
import IcoGitHub from '~icons/mdi/github';
import IcoLinux from '~icons/mdi/linux';
import IcoSteam from '~icons/mdi/steam';
import IcoWindows from '~icons/mdi/windows';

interface GitHubRelease {
	tag_name: string;
	html_url: string;
	prerelease: boolean;
	draft: boolean;
	assets: { name: string; browser_download_url: string }[];
}

const ua = $derived<string>(browser ? navigator.userAgent : '');
// const ua = $derived<string>('linux');
const isWindows = $derived(ua.toLowerCase().includes('windows'));
const isLinux = $derived(
	ua.toLowerCase().includes('linux') ||
		ua.toLowerCase().includes('x11') ||
		ua.toLowerCase().includes('wayland')
);
const platform = $derived(isWindows ? 'windows' : isLinux ? 'linux' : 'unknown');
const isMobile = $derived(
	ua.toLowerCase().includes('mobile') ||
		(navigator as unknown as { userAgentData?: { mobile?: boolean } })?.userAgentData?.mobile
);

const fetchReleases = async () => {
	if (!browser) {
		throw new Error('fetchReleases can only be called in the browser');
	}
	const headers = { Accept: 'application/vnd.github+json' };
	const [stableRes, allRes] = await Promise.all([
		fetch('https://api.github.com/repos/Alia5/steaminputdb.com/releases/latest', { headers }),
		fetch('https://api.github.com/repos/Alia5/steaminputdb.com/releases?per_page=20', {
			headers
		})
	]);
	const stable: GitHubRelease | null = stableRes.ok ? await stableRes.json() : null;
	let prerelease: GitHubRelease | null = null;
	if (allRes.ok) {
		const all: GitHubRelease[] = await allRes.json();
		prerelease = all.find((r) => !r.draft && r.prerelease) ?? null;
	}
	return { stable, prerelease };
};
</script>

<svelte:head>
	<title>SteamInputDB | Buddy App Downloads and Installation</title>
	<meta
		name="description"
		content="Community-driven database of Steam Input configurations using the Steam API." />
	<meta
		name="keywords"
		content="Steam Input DB, Buddy, Buddy-App, Buddy App, Steam DB, DB, Steam Deck, Steam Input, Steam controller configs, controller layouts, community database, Steam API, gamepad configurations, controller presets" />
	<meta
		name="robots"
		content="index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1" />

	<meta property="og:site_name" content="SteamInputDB - BuddyApp" />
	<meta property="og:type" content="website" />
	<meta property="og:title" content="SteamInputDB - BuddyApp" />
	<meta
		property="og:description"
		content="SteamInputDB-Buddy is an official buddy-app for SteamInputDB that provides direct integration into the Steam client " />
	<meta property="og:url" content="https://www.steaminputdb.com/buddy-app/install" />
	<meta property="og:image" content="https://www.steaminputdb.com/ogimage.png" />
	<meta property="og:image:alt" content="SteamInputDB preview image" />

	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="SteamInputDB - BuddyApp" />
	<meta
		name="twitter:description"
		content="SteamInputDB-Buddy is an official buddy-app for SteamInputDB that provides direct integration into the Steam client " />
	<meta name="twitter:image" content="https://www.steaminputdb.com/ogimage.png" />
	<meta name="twitter:image:alt" content="SteamInputDB preview image" />
</svelte:head>

<main>
	<h1><span>SteamInputDB-<strong>Buddy</strong></span> <em>(Beta)</em></h1>
	<section id="about-buddy">
		<AboutBuddy />
	</section>
	<p class="box alert card glass">
		<strong class="card glass"><IcoSteam /> Steam<em>Beta</em>Update required!</strong>
		<span>
			The <em>Steam <strong>Beta</strong> Update</em> is required in order for SteamInputDB-Buddy to be
			functional.
			<br />
			You can opt-into the beta in your Steam Client settings under <em>Interface</em>.
		</span>
	</p>
	<section id="install">
		<h2>Downloads | Installation</h2>
		{#if isLinux}
			<p class="box alert card glass">
				<strong class="card glass"><IcoAlert /> Decky users beware</strong>
				<span>
					We are currently investigating a Decky-plugin as installation method for the Buddy-App
				</span>
			</p>
		{/if}
		{#if !isMobile}
			<h3>Automatic installation</h3>
			<p>
				The easiest way to get everything up and running is to run this script in {isWindows
					? 'Powershell'
					: 'your terminal'}.
				<br />
				<!-- eslint-disable prettier/prettier -->
			This will download and install the latest release of SteamInputDB-Buddy and apply the neccessary modification to your Steam client.
            <!-- eslint-enable prettier/prettier -->
				<br />
			</p>
			<div class="box card glass">
				<strong class="card glass">
					{#if isWindows}
						<IcoWindows />
						Powershell
					{:else if isLinux}
						<IcoLinux />
						Shell
					{/if}
				</strong>
				<code
					{@attach selectAllHandler()}
					style="border: none; box-shadow: none; outline: none; user-select: all; background: transparent; width: 100%; display: block; font-size: 1.2em; padding: 0.5em 1em;">
					{#if isWindows}
						irm https://www.steaminputdb.com/buddy-app/install.ps1 | iex
					{/if}
					{#if isLinux}
						curl -L https://www.steaminputdb.com/buddy-app/install.sh | sh
					{/if}
				</code>
			</div>
			<h3>Direct Downloads</h3>

			<div class="stacked">
				<svelte:boundary pending={spinner}>
					{@const releases = await fetchReleases()}
					{@const stable = releases.stable}
					{@const prerelease = releases.prerelease}
					{@const stableAsset = stable?.assets?.find((a) =>
						a.name.toLowerCase().includes(platform)
					)}
					{@const prereleaseAsset = prerelease?.assets?.find((a) =>
						a.name.toLowerCase().includes(platform)
					)}
					<div transition:fade={{ duration: 196 }} style="display: grid; gap: 1em;">
						{#if stableAsset}
							<a
								href={stableAsset.browser_download_url}
								target="_blank"
								rel="external"
								style="background-color: color(from var(--color-primary) srgb r g b / 0.7);"
								class="button">
								{#if isWindows}
									<IcoWindows />
								{:else if isLinux}
									<IcoLinux />
								{/if}
								Download <em>({stable?.tag_name})</em>
							</a>
						{:else if stable}
							<a href={stable.html_url} target="_blank" rel="external" class="button">
								<IcoGitHub />
								View on GitHub
							</a>
						{/if}

						{#if prereleaseAsset}
							<a
								href={prereleaseAsset.browser_download_url}
								target="_blank"
								rel="external"
								style="background-color: color(from var(--highlight-color) srgb r g b / 0.33); font-size: 1.2em;"
								class="button">
								{#if isWindows}
									<IcoWindows />
								{:else if isLinux}
									<IcoLinux />
								{/if}
								Download pre-release
							</a>
						{:else if prerelease}
							<a href={prerelease.html_url} target="_blank" rel="external" class="button">
								<IcoGitHub />
								View on GitHub
							</a>
						{/if}
					</div>
					{#snippet failed()}
						<!-- blank -->
					{/snippet}
				</svelte:boundary>
			</div>
			<h3>Looking for other Platforms?</h3>
		{:else}
			<p>
				It seems you are accessing this page from a mobile device.
				<br />
				<!-- eslint-disable prettier/prettier -->
				SteamInputDB-Buddy provides direct integration of SteamInputDB with a running Steam Client and must therefore run on a desktop System or Steam Deck / Steam Machine.
				<br />
				Don't worry, you cann still browse the latest releases on GitHub
                <!-- eslint-enable prettier/prettier -->
			</p>
		{/if}
		{@render fallbackContents()}
	</section>
</main>

{#snippet fallbackContents()}
	<div transition:fade|global={{ duration: 196 }}>
		<a
			href="https://github.com/Alia5/steaminputdb.com/releases"
			target="_blank"
			rel="external"
			style="font-size: 1.1em; font-weight: normal; margin-top: 1em;"
			class="button"><IcoGitHub /> View all releases on GitHub</a>
	</div>
{/snippet}

{#snippet spinner()}
	<div transition:fade|global={{ duration: 196 }}>
		<Spinner size="12em" />
	</div>
{/snippet}

<style lang="postcss">
main {
	position: relative;
	display: grid;
	padding: 1em;

	gap: 1em;
	grid-template-rows: min-content;
	grid-template-columns: minmax(min(100%, auto), 50%);
	width: 100%;
	max-width: 1440px;
	justify-self: center;
	height: fit-content;

	text-align: justify;
	height: fit-content;

	:global(h1),
	:global(h2),
	:global(h3),
	:global(h4),
	:global(h5),
	:global(h6) {
		margin-top: 0.5em;
		margin-bottom: 0.2em;
	}
	:global(li) {
		margin-bottom: 0.33em;
	}
	:global(ol),
	:global(ul) {
		overflow: hidden;
		padding-left: 1.5em;
		margin-left: 0;
		margin-bottom: 1em;
	}
	:global(p) {
		margin-bottom: 0.5em;
		text-align: justify;
		word-break: normal;
	}
}

#install {
	display: grid;
	place-items: center;
	gap: 1em;
	width: 100%;
	& p {
		margin-bottom: 0.25em;
	}
	& code {
		font-size: 1.2em;
		padding: 0.5em 1em;
	}
}

.box {
	display: grid;
	border-radius: 0.5em;
	padding: 0;
	--box-color: color(from var(--text-color) srgb r g b / 0.33);
	--heading-color: color(from var(--color-primary) srgb r g b / 0.33);
	border: 1px solid var(--box-color);
	overflow: hidden;
	white-space: normal;
	width: min(100%, 80ch);
	justify-self: center;

	& > :first-child {
		background-color: var(--heading-color);
		font-size: 1.1em;
		padding: 0.5em;
		border-radius: 0;
		display: flex;
		align-items: center;
		gap: 1ch;
		& :global(svg) {
			color: var(--box-color);
		}
		width: 100%;
	}
	display: grid;
	& code {
		text-align: center;
	}
	& > :last-child {
		padding: 1em;
	}
}

p.alert {
	--box-color: firebrick;
	--heading-color: color(from firebrick srgb r g b / 0.33);
}

a.button {
	display: inline-flex;
	font-size: 1.6em;
	align-items: center;
	gap: 0.5em;
	padding: 0.5em 1.5em;
	background: var(--card-color);
	color: var(--text-color);
	border-radius: 0.5em;
	font-weight: bold;
	&:hover,
	&:focus-visible {
		background: var(--inverse-card-color) !important;
		color: var(--color-primary) !important;
	}
	& em {
		font-size: 0.8em;
		opacity: 0.75;
	}
}

#about-buddy {
	font-size: 1.2em;
	& :global(strong) {
		color: var(--color-primary);
	}
}

.stacked {
	display: grid;
	grid-template-areas: 'stack';
	& > * {
		grid-area: stack;
	}
}

h1 {
	justify-self: center;
	& strong {
		color: var(--color-primary);
	}
	& span {
		color: var(--highlight-color);
	}
	& em {
		font-size: 0.8em;
		opacity: 0.75;
	}
}

h2 {
	font-size: 1.6em;
}
h3 {
	font-size: 1.4em;
}
</style>
