<script lang="ts">
import { selectAllHandler } from '$lib/attachments/selectAllHandler.svelte';
import { BuddyState } from '$lib/buddy-app/buddyState.svelte';
import Spinner from '$lib/components/Spinner.svelte';
import { fade, slide } from 'svelte/transition';
import type { PageProps } from './$types';
import { getBuddySettings, getLatestVersion, steamClientStatus } from './requests';

import { goto, invalidateAll } from '$app/navigation';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { clientWithSvelteFetch } from '$lib/buddy-api/client';
import { log } from '$lib/log';
import { toast } from '$lib/toaster/toaster.svelte';

import { resolve } from '$app/paths';
import IconHelp from '~icons/material-symbols/help-outline';
import IconMDIChecked from '~icons/mdi/check-circle-outline';
import IconMDICross from '~icons/mdi/close-circle-outline';
import IconSave from '~icons/mdi/floppy';

const { data }: PageProps = $props();

const normalizeVersion = (version: string) =>
	version?.trim()?.replace(/^v/i, '')?.split('-')?.[0]?.trim() ?? '';

const latestVersionPromise = $derived(data.latestVersion ?? getLatestVersion());
void (async () => {
	const hasBuddyCookie = document.cookie
		.split(';')
		.some((cookie) => cookie.trim().startsWith('buddy-app=enabled'));
	await BuddyState.pingBuddy(fetch);
	if (!hasBuddyCookie && BuddyState.reachable) {
		toast({
			message: 'SteamInputDB Buddy integration enabled!',
			color: 'green'
		});
	}
})();

$effect(() => {
	if (BuddyState.reachable) {
		document.cookie = 'buddy-app=enabled; path=/; max-age=31536000';
	}
});
</script>

<main>
	<h1>SteamInputDB Buddy</h1>
	<section id="buddy-status" class="card glass">
		<div>
			<dl>
				<dt>Status:</dt>
				<dd class="stacked">
					{#if BuddyState.pending}
						<span style="opacity: 0.8;" transition:fade={{ duration: 196 }}>Loading...</span>
					{:else if BuddyState.reachable}
						<span class="center" style="color: green;">Running <IconMDIChecked /></span>
					{:else}
						<span style="color: firebrick;">Could not reach Buddy-App</span>
					{/if}
				</dd>
				{#if BuddyState.pingResponse}
					{#if BuddyState.pingResponse.version}
						<dt>Version:</dt>
						<dd>
							{BuddyState.pingResponse.version}
							{#snippet pingfailed()}
								<em transition:fade={{ duration: 196 }}>(Failed to check for updates)</em>
							{/snippet}
							<svelte:boundary failed={pingfailed}>
								{@const latest = normalizeVersion(await latestVersionPromise)}
								{@const buddyVer = normalizeVersion(BuddyState.pingResponse.version)}
								{#if latest && latest != buddyVer}
									<em transition:fade
										>(Update available: <a
											href="https://github.com/Alia5/steaminputdb.com/releases/latest"
											target="_blank"
											rel="noopener noreferrer">v{latest}</a
										>)</em>
								{:else}
									<em transition:fade={{ duration: 196 }}
										>(You are running the latest version)</em>
								{/if}
							</svelte:boundary>
						</dd>
					{/if}
				{/if}
			</dl>
		</div>
		{#if BuddyState.pending}
			{@render spinner()}
		{:else if !BuddyState.reachable}
			<p transition:fade|global={{ duration: 196 }}>
				Make sure SteamInputDB Buddy is running, and allow SteamInputDB.com to access localhost in
				your browser settings.
			</p>
			<h3>Want to get rid of SteamInputDB Buddy?</h3>
			<button
				class="disable-buddy"
				onclick={() => {
					document.cookie = 'buddy-app=disabled; path=/; max-age=0';
					toast({
						message: 'SteamInputDB Buddy integration disabled',
						color: 'orange'
					});
					goto(resolve('/?buddy-app=disabled'));
				}}>Disable SteamInputDB Buddy integration</button>
		{/if}
	</section>
	{#if BuddyState.reachable}
		<section class="card glass" transition:fade|global={{ duration: 196 }}>
			<h2>Steam connection</h2>
			<br />
			<div class="stacked">
				<!-- eslint-disable-next-line @typescript-eslint/no-explicit-any -->
				{#snippet steamStatusFailed(e: any)}
					<div>
						<p transition:fade={{ duration: 196 }}>
							SteamInputDB-Buddy encountered an issue communicating with your Steam client.
							<br />
							Please make sure Steam is running.
						</p>
						{#if e?.body?.err?.errors}
							<div class="errinfo">
								<strong>Error info:</strong>
								{#each e?.body?.err?.errors as err, idx (idx)}
									<span>{err.message}</span>
								{/each}
							</div>
						{/if}
					</div>
				{/snippet}
				<svelte:boundary pending={spinner} failed={steamStatusFailed}>
					{@const steamStatus = await steamClientStatus()}
					<dl transition:fade|global={{ duration: 196 }}>
						<dt>Connected to Steam:</dt>
						<dd>
							{#if steamStatus.cefRemoteDebugReachable}
								<IconMDIChecked style="color: green; width: 1.6em; height: 1.6em;" />
							{:else}
								<IconMDICross style="color: firebrick; width: 1.6em; height: 1.6em;" />
							{/if}
						</dd>
						<dt>Steam Directory:</dt>
						<dd>
							<code
								{@attach selectAllHandler(
									`outline: 1px solid transparent;
                                        background: rgb(128 128 128 / 0.10);`
								)}>{steamStatus.steamPath}</code>
						</dd>
						<dt>CEF Remote Debug enabled:</dt>
						<dd>
							{#if steamStatus.cefDebugEnableFilePresent}
								<IconMDIChecked style="color: green; width: 1.6em; height: 1.6em;" />
							{:else}
								<IconMDICross style="color: firebrick; width: 1.6em; height: 1.6em;" />
							{/if}
						</dd>
						<dt>Steam running:</dt>
						<dd>
							{#if steamStatus.steamRunning}
								<IconMDIChecked style="color: green; width: 1.6em; height: 1.6em;" />
							{:else}
								<IconMDICross style="color: firebrick; width: 1.6em; height: 1.6em;" />
							{/if}
						</dd>
					</dl>
				</svelte:boundary>
			</div>
		</section>
		<section class="card glass" transition:fade|global={{ duration: 196 }}>
			<h2>Settings</h2>
			<br />
			<div class="stacked">
				{#snippet buddySettingsFailed()}
					<div>
						<p transition:fade={{ duration: 196 }}>
							Could not reach SteamInputDB-Buddy
							<br />
							Please make sure SteamInputDB-Buddy is running and your browser is allowed to make requests
							to "localhost".
						</p>
					</div>
				{/snippet}
				<svelte:boundary pending={spinner} failed={buddySettingsFailed}>
					{@const settings = await getBuddySettings()}
					<form
						transition:fade|global={{ duration: 196 }}
						class="settings"
						onsubmit={async (e) => {
							e.preventDefault();
							const client = clientWithSvelteFetch(fetch);
							try {
								const { data, error: err } = await client.PUT('/v1/settings', {
									body: settings
								});
								if (data) {
									log.info('Config updated successfully');
									toast({
										message: 'Config updated successfully',
										color: 'green'
									});
									void invalidateAll();
								}
								if (err) {
									throw err;
								}
							} catch (e) {
								log.error('Failed to update config', e);
								toast({
									message: 'Failed to update config',
									color: 'firebrick'
								});
							}
						}}>
						<dl>
							<strong>System</strong>
							<dt
								{@attach tooltip({
									content: 'Run SteamInputDB-Buddy upon system Startup'
								})}>
								<label for="autoStart">Run on system startup</label>
								<IconHelp style="width: 1.6em; height: 1.6em;" />
							</dt>
							<dd>
								<input type="checkbox" id="autoStart" bind:checked={settings.autoStart} />
							</dd>
							<strong>Steam UI Integration</strong>
							<dt
								{@attach tooltip({
									snippet: desktopButtonTooltip,
									snippetInDefaultBackground: true
								})}>
								<label for="addDesktopUIEntries">Add Steam Desktop UI Buttons </label>
								<IconHelp style="width: 1.6em; height: 1.6em;" />
							</dt>
							<dd>
								<input
									type="checkbox"
									id="addDesktopUIEntries"
									bind:checked={settings.addDesktopUIEntries} />
							</dd>
							<dt
								{@attach tooltip({
									snippet: bpmButtonTooltip,
									snippetInDefaultBackground: true
								})}>
								<label for="addBigPictureUIEntries">Add Steam Big Picture UI Buttons </label>
								<IconHelp style="width: 1.6em; height: 1.6em;" />
							</dt>
							<dd>
								<input
									type="checkbox"
									id="addBigPictureUIEntries"
									bind:checked={settings.addBigPictureUIEntries} />
							</dd>
							<dt
								{@attach tooltip({
									content:
										"Use Steam's internal browser (as opposed to your default browser) when SteamInputDB is opened from the added Buttons inside Steam"
								})}>
								<label for="desktopUseSteamBrowser"
									>Use Steam's internal browser for Desktop UI
								</label>
								<IconHelp style="width: 1.6em; height: 1.6em;" />
							</dt>
							<dd>
								<input
									type="checkbox"
									id="desktopUseSteamBrowser"
									bind:checked={settings.desktopUseSteamBrowser} />
							</dd>
						</dl>
						<button type="submit"><IconSave />Save</button>
					</form>
				</svelte:boundary>
			</div>
		</section>
	{/if}
</main>

{#snippet spinner()}
	<div class="spinner" transition:slide|global={{ duration: 196 }}>
		<Spinner size="12em" />
	</div>
{/snippet}

{#snippet desktopButtonTooltip()}
	<div style="display: grid; gap: 1em;">
		<span>Adds a SteamInputDB button to a games Library page in Steams Desktop UI</span>
		<span>The button opens the SteamInputDB page for the currently selected Game/App</span>
		<enhanced:img src="$lib/assets/buddy-app/desktop-button.png" alt="Example screenshot"></enhanced:img>
	</div>
{/snippet}

{#snippet bpmButtonTooltip()}
	<div style="display: grid; gap: 1em;">
		<span>Adds a SteamInputDB button to a games Library page in Steams Big Picture UI</span>
		<span>The button opens the SteamInputDB page for the currently selected Game/App</span>
		<enhanced:img src="$lib/assets/buddy-app/bpm-button.png" alt="Example screenshot"></enhanced:img>
	</div>
{/snippet}

<style lang="postcss">
main {
	position: relative;
	isolation: isolate;
	display: grid;
	padding: 1em;

	gap: 1em;
	grid-template-rows: min-content;
	grid-template-columns: minmax(min(100%, auto), 50%);
	width: 100%;
	max-width: 1440px;
	justify-self: center;
	height: fit-content;
}

h1 {
	margin-bottom: 0.6em;
}

.center {
	display: flex;
	flex-flow: row wrap;
	gap: 1ch;
	align-items: center;
}

.stacked {
	display: grid;
	grid-template-areas: 'a';
	& > * {
		grid-area: a;
	}
}

.errinfo {
	margin-top: 1em;
	display: grid;
	& span {
		padding-left: 1em;
	}
}

#buddy-status {
	display: grid;
	padding: 1em 2em;
	align-items: center;
	gap: 1em;
	width: 100%;

	& dl {
		& > :nth-child(-n + 2) {
			font-size: 1.4em;
		}
	}

	& > :nth-child(2) {
		grid-column: 1 / span 1;
		grid-row: 2 / span 1;
	}
}
.spinner {
	justify-self: center;
}

dl {
	display: grid;
	place-items: center;
	grid-template-columns: auto 1fr;
	grid-column-gap: 1em;
	grid-row-gap: 0.5em;
	width: 100%;

	& > strong {
		grid-column: 1 / -1;
		justify-self: start;
		font-size: 1.2em;
		margin-bottom: 0.25em;
		width: 100%;
		&:not(:nth-of-type(1)) {
			margin-top: 1em;
			position: relative;
			&::before {
				content: '';
				position: absolute;
				inset: 0;
				border-top: 1px solid var(--text-color);
				top: -0.75em;
				opacity: 0.33;
			}
		}
	}

	&:has(> strong) {
		dt {
			margin-left: 1em;
		}
	}

	& dt {
		font-weight: bold;
		justify-self: start;
		display: flex;
		align-items: center;
		gap: 1ch;
	}
	dd {
		max-width: 100%;
		overflow: clip;
		overflow-clip-margin: 2em;
		display: flex;
		align-items: center;
		gap: 1ch;
		& a {
			font-weight: bold;
		}
	}

	& > :nth-child(2n) {
		justify-self: start;
		& em {
			opacity: 0.75;
			& a {
				font-weight: bold;
			}
		}
	}
}

form.settings {
	display: flex;
	flex-flow: row wrap;
	& > :first-child {
		flex-shrink: 1;
		width: auto;
	}
	& > :last-child {
		margin-left: auto;
		height: fit-content;
		align-self: end;
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		gap: 1ch;
	}
}

button:not(.plain) {
	color: var(--text-color-dark);
	font-weight: bold;
	background:
		linear-gradient(
			215deg,
			color-mix(in srgb, var(--card-color), transparent 75%) 0%,
			color-mix(in srgb, var(--card-color), transparent 90%) 70%
		),
		var(--bg-noise-transparent);
	background-color: color-mix(in srgb, var(--color-primary), transparent 20%);

	&[disabled] {
		opacity: 0.5;
	}

	&:hover,
	&:focus-visible {
		color: var(--text-color-dark) !important;
		background-color: color-mix(in srgb, var(--color-primary), rgb(128 128 255 / 0.8) 50%);
	}
}

code {
	display: block;
}

.disable-buddy {
	margin-right: auto;
	background-color: color-mix(in srgb, firebrick, transparent 20%) !important;
	font-weight: bold;
	&:hover,
	&:focus-visible {
		color: var(--text-color-dark) !important;
	}
}
</style>
