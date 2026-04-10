<script lang="ts">
import { browser } from '$app/environment';
import type { components } from '$lib/api/openapi';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { client } from '$lib/buddy-api/client';
import type { components as buddyComponents } from '$lib/buddy-api/openapi';
import BPMOption from '$lib/components/BPM_Select/BPM_option.svelte';
import BPMSelect from '$lib/components/BPM_Select/BPM_Select.svelte';
import Modal from '$lib/components/Modal.svelte';
import Spinner from '$lib/components/Spinner.svelte';
import { toast } from '$lib/toaster/toaster.svelte';
import { fade, slide } from 'svelte/transition';

import IconChevronDown from '~icons/mdi/chevron-down';
import IconSteam from '~icons/mdi/steam';

const {
	fileInfo,
	appInfo,
	controllers,
	apps
}: {
	fileInfo: components['schemas']['ConfigItem'];
	appInfo?: components['schemas']['AppItem'];
	controllers: buddyComponents['schemas']['ControllerResponse'][];
	apps: buddyComponents['schemas']['AppResponse'][];
} = $props();

const defaultAppID = $derived(
	apps?.find((a) => {
		if (a.appid === appInfo?.app_id) {
			return true;
		}
		if (a.isNonSteam && a.name.toLowerCase() === fileInfo.app_id_string?.toLowerCase()) {
			return true;
		}
		return false;
	})?.appid
);

let selectedAppID = $derived(defaultAppID ?? -1);
let selectedController = $derived(
	controllers.find((c) => c.type == fileInfo.controller_type)?.index ?? controllers?.[0]?.index ?? -1
);
let showConfigurator = $state(true);
let isApplyingConfig = $state(false);

let isSteamBigPicture = $derived(browser ? navigator?.userAgent?.includes('Steam Gamepad') : false);

const applyConfig = async (appId: string | number, controllerIdx: number, openConfigurator: boolean) => {
	if (!fileInfo.file_id) {
		// cannot happen
		throw new Error('No fileId');
	}
	await client.POST('/v1/steam/apply_config', {
		body: {
			appId: typeof appId === 'string' ? parseInt(appId, 10) : appId,
			controllerIndex: controllerIdx,
			workshopItemId: `${fileInfo.file_id}`
		}
	});
	if (openConfigurator) {
		await client.POST('/v1/steam/open_configurator', {
			body: {
				appId: typeof appId === 'string' ? parseInt(appId, 10) : appId
			}
		});
	}
};

const applyButtonHandler = () => {
	if ((controllers?.length ?? 0) <= 0 || selectedController < 0) {
		toast({
			message:
				'SteamInputDB-Buddy needs a controller connected in order to apply configs\nPlease connect a controller',
			color: 'firebrick'
		});
		return;
	}
	if (!defaultAppID) {
		toast({
			message: 'You must own the game to apply configs!\nYou may apply this config to any game you own',
			color: 'firebrick'
		});
		dialogOpen = true;
		return;
	}
	isApplyingConfig = true;
	applyConfig(defaultAppID, controllers?.[0]?.index ?? selectedController, false)
		.then(() => {
			toast({
				message: 'Config applied successfully!',
				color: 'seagreen'
			});
		})
		.catch((err) => {
			toast({
				message: `Failed to apply config\n${err.message}`,
				color: 'firebrick'
			});
		})
		.finally(() => {
			isApplyingConfig = false;
		});
};
let dialogOpen = $state(false);
</script>

<div class="buddy-preview-button button blue">
	<button
		class="button"
		{@attach tooltip({
			content: 'Apply this config to your primary controller using SteamInputDB Buddy',
			outDelay: 200,
			arrow: true,
			arrowFollowCursor: true
		})}
		disabled={isApplyingConfig}
		onclick={applyButtonHandler}>
		{#if isApplyingConfig}
			<Spinner size="1.4em" thickness="2px" />
		{:else}
			<IconSteam style="width: 1.4em; height: 1.4em;" />
		{/if}
		<span>Apply</span>
	</button>
	<button
		id="advanced-apply-button"
		class="button"
		disabled={isApplyingConfig}
		onclick={() => (dialogOpen = true)}
		style="anchor-name: --advanced-apply-button">
		<IconChevronDown style="width: 1.4em; height: 1.4em;" />
	</button>
</div>

<Modal bind:open={dialogOpen}>
	<div class="card modal-card">
		<form
			transition:slide|global
			onsubmit={(e) => {
				e.preventDefault();
				isApplyingConfig = true;
				applyConfig(selectedAppID, selectedController, showConfigurator)
					.then(() => {
						toast({
							message: 'Config applied successfully!',
							color: 'seagreen'
						});
					})
					.catch((err) => {
						toast({
							message: `Failed to apply config\n${err.message}`,
							color: 'firebrick'
						});
					})
					.finally(() => {
						isApplyingConfig = false;
						dialogOpen = false;
					});
			}}>
			<strong>Apply To</strong>
			{#if isSteamBigPicture}
				<BPMSelect name="App" bind:value={selectedAppID}>
					{#snippet children({ ...rest })}
						<span>App:</span>
						{#each apps as app (app.appid)}
							<BPMOption value={app.appid} {...rest}
								>{app.name} {app.isNonSteam ? '(Non-Steam)' : ''}</BPMOption>
						{/each}
						<IconChevronDown style="width: 1.6em; height: 1.6em;" />
					{/snippet}
				</BPMSelect>
				<BPMSelect name="Controller" bind:value={selectedController}>
					{#snippet children({ ...rest })}
						<span>Controller:</span>
						{#each controllers as controller (controller.index)}
							<BPMOption value={controller.index} {...rest}>{controller.name}</BPMOption>
						{/each}
						<IconChevronDown style="width: 1.6em; height: 1.6em;" />
					{/snippet}
				</BPMSelect>
			{:else}
				<label for="app">
					<span>App</span>
					<select id="app" name="app" bind:value={selectedAppID}>
						{#each apps as app (app.appid)}
							<option value={app.appid}
								>{app.name} {app.isNonSteam ? '(Non-Steam)' : ''}</option>
						{/each}
					</select>
					<IconChevronDown style="width: 1.6em; height: 1.6em;" />
				</label>
				<label for="controller">
					<span>Controller</span>
					<select id="controller" name="controller" bind:value={selectedController}>
						{#each controllers as controller (controller.index)}
							<option value={controller.index}>
								{controller.name}
							</option>
						{/each}
					</select>
					<IconChevronDown style="width: 1.6em; height: 1.6em;" />
				</label>
			{/if}
			<label for="showConfigurator">
				<span>Open Steam Input configurator</span>
				<input
					type="checkbox"
					id="showConfigurator"
					name="showConfigurator"
					bind:checked={showConfigurator} />
			</label>
			<div>
				{#if isApplyingConfig}
					<div class="spinner" transition:fade>
						<Spinner size="3em" thickness="2px" />
					</div>
				{/if}
				<button
					type="submit"
					disabled={selectedController < 0 || selectedAppID < 0 || isApplyingConfig}>
					<IconSteam style="width: 1.4em; height: 1.4em;" />
					<span>Apply</span>
				</button>
			</div>
		</form>
	</div>
</Modal>

<style lang="postcss">
.buddy-preview-button {
	display: grid;
	grid-template-columns: 1fr auto;
	width: 100%;
	padding: 0;

	background:
		linear-gradient(
			215deg,
			color-mix(in srgb, var(--card-color), transparent 75%) 0%,
			color-mix(in srgb, var(--card-color), transparent 90%) 70%
		),
		var(--bg-noise-transparent);
	backdrop-filter: blur(6px);
	background-color: #1a9fff;

	& > button {
		background: transparent;
		background-color: color(from #1a9fff srgb r g b / 0.5);
		&:hover,
		&:focus-within {
			background-color: color(from #1a9fff srgb r g b / 0.7);
		}
	}

	& > :first-child {
		border-top-right-radius: 0;
		border-bottom-right-radius: 0;
		display: grid;
		align-items: center;
		justify-content: center;
		gap: 0.5ch;
		font-weight: bold;
	}
	& > * {
	}
	& > :last-child {
		border-top-left-radius: 0;
		border-bottom-left-radius: 0;
		padding: 1em 0.5em;
	}
}

button {
	display: flex;
	font-weight: bold;
	gap: 1ch;
	&:hover,
	&:focus-visible {
		color: var(--text-color-dark) !important;
		background-color: color(from var(--color-primary) srgb r g b / 0.5);
	}
	&[disabled] {
		cursor: not-allowed;
		opacity: 0.5;
	}
}

.modal-card {
	position: fixed;
	left: 50%;
	top: 50%;
	translate: -50% -50%;

	padding: 0;
}

form {
	padding: 1em;
	display: grid;
	gap: 1em;
	grid-template-columns: auto auto;
	gap: 1em;
	place-items: center;

	& > :first-child {
		grid-column: 1/-1;
		margin-right: auto;
	}
	& > :last-child {
		grid-column: 1/-1;
		display: grid;
		margin-left: auto;
		grid-auto-flow: column;
		gap: 2em;
		& > :last-child {
			justify-self: end;
			background-color: #1a9fff;
			&:hover,
			&:focus-visible {
				background-color: color-mix(in srgb, #1a9fff, var(--color-primary) 50%);
			}
		}
	}
	& :global(.bpm-select) {
		width: 100%;
		grid-column: 1/-1;
	}
}

select {
	font-style: inherit;
	color: var(--text-color);
	cursor: pointer;
	appearance: none;
	padding-right: 2em;
	position: relative;
	width: 100%;
	background: var(--card-background-noise);
	background: transparent;
	border: 1px solid color-mix(in srgb, var(--text-color), transparent 90%);
	padding: 0.5em 1em;
	color: var(--text-color);
	border-radius: 0.5rem;
	box-shadow: 0 1px 4px 0 rgb(0 0 0 / 0.25);

	outline: 1px solid transparent;

	transition: all var(--transition-duration) var(--default-ease);
	&:hover,
	&:focus-visible {
		outline: 0.16em solid var(--color-primary);
	}

	& > option {
		color: var(--text-color);
		background-color: var(--card-color);
		&:checked {
			background-color: var(--color-primary);
			color: var(--text-color-dark);
		}
		&:hover {
			background-color: var(--color-primary);
			color: var(--text-color-dark);
		}
	}
}

label {
	display: grid;
	grid-column: 1/-1;
	grid-template-columns: subgrid;
	margin-left: 1em;
}

label:has(select) {
	position: relative;

	& > :global(svg) {
		position: absolute;
		right: 0.5em;
		bottom: 0;
		translate: 0 -25%;
	}
	align-items: center;
}
</style>
