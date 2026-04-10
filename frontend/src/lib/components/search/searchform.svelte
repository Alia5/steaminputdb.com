<script lang="ts">
import { enhance } from '$app/forms';
import Searchbar from '$lib/components/search/searchbar.svelte';
import Icon from '@iconify/svelte';
import { tick } from 'svelte';
import { cubicInOut } from 'svelte/easing';
import type { HTMLFormAttributes } from 'svelte/elements';
import { fade, slide } from 'svelte/transition';

import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { page } from '$app/state';
import { BuddyState } from '$lib/buddy-app/buddyState.svelte';
import BPMSelect from '$lib/components/BPM_Select/BPM_Select.svelte';
import BPMOption from '$lib/components/BPM_Select/BPM_option.svelte';
import IcoDropdown from '~icons/mdi/chevron-down';
import { CONTROLLER_LIST } from './controllerlist.svelte';

let {
	showAdvancedFilters = true,
	showFeatureFilter = false,
	showExcludedFeatureFilter = false,
	form = $bindable(),
	// eslint-disable-next-line no-useless-assignment
	searchtext = $bindable(),
	disabled = false,
	method = 'GET',
	values = $bindable({}),
	submitOnChange = false,
	showTotalCount = true,
	enhanceParams,
	...rest
}: {
	showAdvancedFilters?: boolean;
	showFeatureFilter?: boolean;
	showExcludedFeatureFilter?: boolean;
	form?: HTMLFormElement;
	searchtext?: string;
	disabled?: boolean;
	method?: string;
	values?: Record<string, unknown>;
	submitOnChange?: boolean;
	showTotalCount?: boolean | number;
	enhanceParams?: Parameters<typeof enhance>[1];
} & HTMLFormAttributes = $props();

const changeSubmitHandler = () => {
	if (submitOnChange) {
		form!.requestSubmit();
	}
};
$effect(() => {
	if (!browser) {
		return;
	}
	if (
		BuddyState.connectedControllersChanged &&
		!values['controller_type'] &&
		BuddyState.connectedControllers?.length
	) {
		tick().then(() => {
			if (values['controller_type']) {
				return;
			}
			values['controller_type'] = BuddyState.connectedControllers?.[0]?.type;
			tick().then(() => {
				if (values['controller_type']) {
					page.url.searchParams.set('controller_type', values['controller_type'] as string);
					// form?.requestSubmit();
					// eslint-disable-next-line svelte/no-navigation-without-resolve
					goto(page.url, {
						replaceState: true,
						noScroll: true,
						invalidate: [page.url]
					});
				}
			});
		});
	}
});
let showMoreControllers = $state(false);
let isSteamBigPicture = $derived(browser ? navigator?.userAgent?.includes('Steam Gamepad') : false);
</script>

{#if method === 'POST' || method === 'post'}
	<form bind:this={form} class="card glass" data-sveltekit-noscroll {...rest} use:enhance={enhanceParams}>
		{@render formcontents()}
	</form>
{:else}
	<form bind:this={form} class="card glass" data-sveltekit-noscroll {...rest}>
		{@render formcontents()}
	</form>
{/if}

{#snippet formcontents()}
	<div>
		<Searchbar
			name="searchtext"
			placeholder="SteamInput configuration..."
			disabled={disabled}
			bind:value={values.searchtext}
			inlineButton={false} />
		<button type="submit" disabled={disabled}>Search</button>
		{#if isSteamBigPicture}
			<BPMSelect
				bind:value={values['sort-by']}
				name="sorting"
				onchange={(v) => {
					values['sort-by'] = v;
					changeSubmitHandler();
				}}
				disabled={disabled}>
				{#snippet children({ ...rest })}
					<span>Sort by:</span>
					{#if !values['sort-by']}
						<span>Rank</span>
					{/if}
					<BPMOption value="vote" {...rest}>Rank</BPMOption>
					<BPMOption value="publication" {...rest}>Date</BPMOption>
					<BPMOption value="trend" {...rest}>Trend (30 days)</BPMOption>
					<BPMOption value="votes_asc" {...rest}>Votes (ascending)</BPMOption>
					<BPMOption value="votes_up" {...rest}>Votes (upvotes)</BPMOption>
					<BPMOption value="text_search" {...rest}>Relevance</BPMOption>
					<BPMOption value="playtime_trend" {...rest}>Playtime trend (30 days)</BPMOption>
					<BPMOption value="total_playtime" {...rest}>Total playtime</BPMOption>
					<BPMOption value="avg_playtime_trend" {...rest}>Average playtime trend</BPMOption>
					<BPMOption value="lifetime_avg_playtime" {...rest}
						>Average playtime since upload</BPMOption>
					<BPMOption value="playtime_sessions_trend" {...rest}>Sessions trend (30 days)</BPMOption>
					<BPMOption value="lifetime_playtime_sessions" {...rest}>Lifetime sessions</BPMOption>
					<IcoDropdown />
				{/snippet}
			</BPMSelect>
		{/if}
		<label for="sort-by" style={isSteamBigPicture ? 'display: none;' : ''}>
			<span>Sort by:</span>
			<select
				id="sort-by"
				name="sort-by"
				disabled={disabled}
				bind:value={values['sort-by']}
				onchange={changeSubmitHandler}>
				<option value="vote">Rank</option>
				<option value="publication">Date</option>
				<option value="trend">Trend (30 days)</option>
				<option value="votes_asc">Votes (ascending)</option>
				<option value="votes_up">Votes (upvotes)</option>
				<option value="text_search">Relevance</option>
				<option value="playtime_trend">Playtime trend (30 days)</option>
				<option value="total_playtime">Total playtime</option>
				<option value="avg_playtime_trend">Average playtime trend</option>
				<option value="lifetime_avg_playtime">Average playtime since upload</option>
				<option value="playtime_sessions_trend">Sessions trend (30 days)</option>
				<option value="lifetime_playtime_sessions">Lifetime sessions</option>
			</select>

			<IcoDropdown />
		</label>
	</div>
	{#if typeof showTotalCount === 'number'}
		<dl transition:fade={{ duration: 196, easing: cubicInOut }}>
			<dt>Total matching criteria</dt>
			<dd>{showTotalCount ?? 0}</dd>
		</dl>
	{/if}
	<button
		type="button"
		class="filter"
		onclick={() => {
			showAdvancedFilters = !showAdvancedFilters;
		}}
		>Advanced Filters {#if showAdvancedFilters}
			<Icon icon="mdi:chevron-up" height="1.8em" />
		{:else}
			<Icon icon="mdi:chevron-down" height="1.8em" />
		{/if}</button>
	{#if showAdvancedFilters}
		<fieldset
			id="controller-type"
			transition:slide={{ duration: 196, easing: cubicInOut }}
			onclickcapture={(e) => {
				const target = e.target;
				if (!(target instanceof HTMLInputElement)) {
					return;
				}
				if (target.type !== 'radio') {
					return;
				}
				if (target.name !== 'controller_type') {
					return;
				}

				if (values['controller_type'] == target.value) {
					values['controller_type'] = undefined;
					tick().then(() => {
						changeSubmitHandler();
					});
				}
			}}
			disabled={disabled}>
			<legend><span>Controller Type</span></legend>

			{#if browser && document.cookie?.includes('buddy-app=enabled') && BuddyState.reachable}
				<svelte:boundary pending={defaultControllerList} failed={defaultControllerList}>
					{@const connectedControllers = (await BuddyState.fetchConnectedControllers()) as {
						data: {
							type: string;
						}[];
					}}
					{#if connectedControllers.data?.length}
						{@render controllerList({
							filter: connectedControllers.data.map(
								(c) =>
									CONTROLLER_LIST.find((cl) => cl.type === c.type)?.type ??
									'controller_generic'
							),
							type: 'include'
						})}
					{:else}
						{@render defaultControllerList()}
					{/if}
				</svelte:boundary>
			{:else}
				{@render defaultControllerList()}
			{/if}

			{#snippet defaultControllerList()}
				{@render controllerList({
					filter: [
						'controller_mobile_touch',
						'controller_rog_ally',
						'controller_legion_go_s',
						'controller_steamcontroller_headcrab',
						'controller_hori_steam',
						'controller_xboxelite',
						'controller_ps5_edge',
						'controller_ps3'
					],
					type: 'exclude'
				})}
			{/snippet}

			<div class="show-more">
				<button type="button" onclick={() => (showMoreControllers = !showMoreControllers)}>
					<span class="show-more-text">{showMoreControllers ? 'Show Less' : 'Show More'}</span>
					{#if showMoreControllers}
						<Icon icon="mdi:chevron-up" height="1.6em" />
					{:else}
						<Icon icon="mdi:chevron-down" height="1.6em" />
					{/if}
				</button>
			</div>
			{#if showMoreControllers}
				{#if browser && document.cookie?.includes('buddy-app=enabled') && BuddyState.reachable}
					<svelte:boundary pending={defaultControllerList} failed={defaultControllerList}>
						{@const connectedControllers = (await BuddyState.fetchConnectedControllers()) as {
							data: {
								type: string;
							}[];
						}}
						{#if connectedControllers.data?.length}
							{@render controllerList({
								filter: connectedControllers.data.map(
									(c) =>
										CONTROLLER_LIST.find((cl) => cl.type === c.type)?.type ??
										'controller_generic'
								),
								type: 'exclude'
							})}
						{:else}
							{@render defaultMoreControllerList()}
						{/if}
					</svelte:boundary>
				{:else}
					{@render defaultMoreControllerList()}
				{/if}

				{#snippet defaultMoreControllerList()}
					{@render controllerList({
						filter: [
							'controller_mobile_touch',
							'controller_rog_ally',
							'controller_legion_go_s',
							'controller_steamcontroller_headcrab',
							'controller_hori_steam',
							'controller_xboxelite',
							'controller_ps5_edge',
							'controller_ps3'
						],
						type: 'include'
					})}
				{/snippet}
			{/if}
		</fieldset>
	{/if}
	{#if showAdvancedFilters}
		<fieldset id="features" transition:slide={{ duration: 196, easing: cubicInOut }} disabled={disabled}>
			<legend
				><button
					type="button"
					class="plain"
					onclick={() => {
						showFeatureFilter = !showFeatureFilter;
					}}
					>Must have features
					{#if showFeatureFilter}
						<Icon icon="mdi:chevron-up" height="1.6em" />
					{:else}
						<Icon icon="mdi:chevron-down" height="1.6em" />
					{/if}
				</button></legend>
			{#if showFeatureFilter}
				{@render featurefilters(values)}
			{/if}
		</fieldset>
	{/if}
	{#if showAdvancedFilters}
		<fieldset
			id="excluded-features"
			transition:slide={{ duration: 196, easing: cubicInOut }}
			disabled={disabled}>
			<legend>
				<button
					type="button"
					class="plain"
					onclick={() => {
						showExcludedFeatureFilter = !showExcludedFeatureFilter;
					}}
					>Must <strong>not</strong> have features
					{#if showExcludedFeatureFilter}
						<Icon icon="mdi:chevron-up" height="1.6em" />
					{:else}
						<Icon icon="mdi:chevron-down" height="1.6em" />
					{/if}
				</button></legend>
			{#if showExcludedFeatureFilter}
				{@render featurefilters(values, 'exclude_')}
			{/if}
		</fieldset>
	{/if}
{/snippet}

{#snippet controllerList({
	filter,
	type
}: {
	filter: `${(typeof CONTROLLER_LIST)[number]['type']}`[];
	type: 'include' | 'exclude';
})}
	{#each CONTROLLER_LIST.filter((c) => {
		return type === 'include' ? filter.includes(c.type) : !filter.includes(c.type);
	}) as controller (controller.type)}
		<label for={controller.type} transition:slide|global={{ duration: 196, easing: cubicInOut }}>
			<input
				type="radio"
				id={controller.type}
				name="controller_type"
				value={controller.type}
				bind:group={values['controller_type'] as string}
				onchange={changeSubmitHandler} />
			<controller.icon style="width: 1.2em; height: 1.2em;" />
			<span> {controller.niceName} </span>
		</label>
	{/each}
{/snippet}

{#snippet featurefilters(bindMap: Record<string, unknown>, prefix = '')}
	<label for={`${prefix}feature_gamepad`} transition:slide|global={{ duration: 196, easing: cubicInOut }}>
		<input
			type="checkbox"
			id={`${prefix}feature_gamepad`}
			name={`${prefix}feature_gamepad`}
			bind:checked={bindMap[`${prefix}feature_gamepad`] as boolean}
			onchange={changeSubmitHandler} />
		<Icon icon="mdi:controller" width="1.2em" />
		<span>Gamepad Inputs</span>
	</label>
	<label for={`${prefix}feature_keyboard`} transition:slide|global={{ duration: 196, easing: cubicInOut }}>
		<!-- actually typo in valves data: feature_keboard -->
		<input
			type="checkbox"
			id={`${prefix}feature_keyboard`}
			name={`${prefix}feature_keboard`}
			bind:checked={bindMap[`${prefix}feature_keboard`] as boolean}
			onchange={changeSubmitHandler} />
		<Icon icon="mdi:keyboard" width="1.2em" />
		<span>Keyboard Inputs</span>
	</label>
	<label for={`${prefix}feature_mouse`} transition:slide|global={{ duration: 196, easing: cubicInOut }}>
		<input
			type="checkbox"
			id={`${prefix}feature_mouse`}
			name={`${prefix}feature_mouse`}
			bind:checked={bindMap[`${prefix}feature_mouse`] as boolean}
			onchange={changeSubmitHandler} />
		<Icon icon="mdi:mouse" width="1.2em" />
		<span>Mouse Inputs</span>
	</label>
	<label for={`${prefix}feature_gyro`} transition:slide|global={{ duration: 196, easing: cubicInOut }}>
		<input
			type="checkbox"
			id={`${prefix}feature_gyro`}
			name={`${prefix}feature_gyro`}
			bind:checked={bindMap[`${prefix}feature_gyro`] as boolean}
			onchange={changeSubmitHandler} />
		<Icon icon="game-icons:gyroscope" width="1.2em" />
		<span>Gyro Inputs</span>
	</label>
	<label for={`${prefix}feature_touchmenu`} transition:slide|global={{ duration: 196, easing: cubicInOut }}>
		<input
			type="checkbox"
			id={`${prefix}feature_touchmenu`}
			name={`${prefix}feature_touchmenu`}
			bind:checked={bindMap[`${prefix}feature_touchmenu`] as boolean}
			onchange={changeSubmitHandler} />
		<Icon icon="mdi:gesture-touch" width="1.2em" />
		<span>Touch Menus</span>
	</label>
	<label
		for={`${prefix}feature_radialmenu`}
		transition:slide|global={{ duration: 196, easing: cubicInOut }}>
		<input
			type="checkbox"
			id={`${prefix}feature_radialmenu`}
			name={`${prefix}feature_radialmenu`}
			bind:checked={bindMap[`${prefix}feature_radialmenu`] as boolean}
			onchange={changeSubmitHandler} />
		<Icon icon="material-symbols:joystick" width="1.2em" />
		<span>Radial Menus</span>
	</label>
	<label for={`${prefix}feature_modeshift`} transition:slide|global={{ duration: 196, easing: cubicInOut }}>
		<input
			type="checkbox"
			id={`${prefix}feature_modeshift`}
			name={`${prefix}feature_modeshift`}
			bind:checked={bindMap[`${prefix}feature_modeshift`] as boolean}
			onchange={changeSubmitHandler} />
		<Icon icon="material-symbols:layers-rounded" width="1.2em" />
		<span>Mode Shifts</span>
	</label>
	<label
		for={`${prefix}feature_mouseregion`}
		transition:slide|global={{ duration: 196, easing: cubicInOut }}>
		<input
			type="checkbox"
			id={`${prefix}feature_mouseregion`}
			name={`${prefix}feature_mouseregion`}
			bind:checked={bindMap[`${prefix}feature_mouseregion`] as boolean}
			onchange={changeSubmitHandler} />
		<Icon icon="fluent:cursor-hover-16-filled" width="1.2em" />
		<span>Mouse Regions</span>
	</label>
	<label for={`${prefix}feature_actionset`} transition:slide|global={{ duration: 196, easing: cubicInOut }}>
		<input
			type="checkbox"
			id={`${prefix}feature_actionset`}
			name={`${prefix}feature_actionset`}
			bind:checked={bindMap[`${prefix}feature_actionset`] as boolean}
			onchange={changeSubmitHandler} />
		<Icon icon="mdi:set-right" width="1.2em" />
		<span>Action Sets</span>
	</label>
{/snippet}

<style lang="postcss">
form {
	display: flex;
	flex-flow: row wrap;
	width: 100%;
	gap: 1em;
	backdrop-filter: blur(12px);

	max-width: calc(100dvw -2em);

	& > :first-child {
		width: 100%;
		flex-grow: 1;
		gap: 1em;
		width: 100%;
		display: flex;
		flex-flow: row wrap-reverse;
		position: relative;
		margin-bottom: 1em;
		:global(> :first-child) {
			flex-grow: 1;
			max-width: max(52ch, 25dvw);
		}
	}

	label {
		display: flex;
		gap: 0.5em;
		align-items: center;
	}

	:global(.bpm-select) {
		margin-left: auto;
		font-size: 1.2em;
	}
	label[for='sort-by'] {
		margin-left: auto;
		display: grid;
		grid-template-columns: auto auto;
		gap: 0.5em;
		align-items: center;
		font-size: 1.2em;
		position: relative;
		isolation: isolate;
		border: 1px solid color-mix(in srgb, var(--text-color), transparent 90%);
		padding: 0.5em 1em;
		box-shadow: 0 1px 4px 0 rgb(0 0 0 / 0.25);
		border-radius: 100vw;
		transition: all var(--transition-duration) var(--default-ease);

		& :global([disabled]) {
			opacity: 0.5;
		}

		&:has([disabled]) {
			opacity: 0.5;
		}

		&:hover,
		&:focus-within {
			outline: 0.1em solid var(--color-primary);
			box-shadow: 0 0 1.3em -0.4em var(--color-primary);
		}

		& > :first-child {
			white-space: nowrap;
		}
		:global(> :last-child) {
			content: '';
			color: var(--text-color);
			position: absolute;
			z-index: 1;
			height: 100%;
			width: 1.4em;
			top: 50%;
			translate: 0 -50%;
			right: 0.5em;
			background-size: contain;
			pointer-events: none;
		}
	}
}

select {
	font-style: inherit;
	background: transparent;
	border: 1px solid transparent;
	outline: none;
	color: var(--text-color);
	cursor: pointer;
	appearance: none;
	padding-right: 2em;
	position: relative;
	width: 100%;

	& option {
		color: var(--text-color);
		background: var(--card-color);
	}
}

fieldset {
	border-radius: 0.5em;
	padding: 0 1em;
	background: var(--card-background-noise);
	border: 1px solid color-mix(in srgb, var(--text-color), transparent 90%);
	position: relative;
	box-shadow: inset 0.1em 0.2em 0.5em 0 light-dark(#0f0f0f27, #0e0e0e7e);

	width: 100%;
	transition: all var(--transition-duration);
	&:has(label) {
		padding: 1em;
	}

	&[disabled] {
		opacity: 0.5;
	}

	& legend {
		font-size: 1.1em;
		border-radius: 0.5em;
		background: var(--card-background-noise);
		isolation: isolate;
		position: relative;
		& > :first-child {
			padding: 0.25em 0.5em;
			display: block;
		}
		&::before {
			content: '';
			position: absolute;
			inset: 0;
			z-index: -1;
			border-radius: inherit;
			border: 1px solid transparent;
			background: linear-gradient(0deg, transparent 0%, transparent 40%, var(--text-color) 100%)
				border-box;
			mask:
				linear-gradient(black, black) border-box,
				linear-gradient(black, black) padding-box;
			mask-composite: subtract;
			opacity: 0.2;
		}
	}
}

#controller-type {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(25ch, 1fr));
	gap: 1em;
	position: relative;
	inset: 0;

	& input {
		min-width: 1.4em;
		min-height: 1.4em;
	}

	& label {
		display: grid;
		grid-template-columns: min-content min-content auto;
		align-items: center;
		gap: 0.5em;
	}
}
#features,
#excluded-features {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(22ch, 1fr));
	gap: 1em;
	position: relative;
	inset: 0;

	& input {
		min-width: 1.4em;
		min-height: 1.4em;
	}

	& label {
		display: grid;
		grid-template-columns: min-content min-content auto;
		align-items: center;
		gap: 0.5em;
	}
}
button.plain {
	background: none;
	border: none;
	padding: 0;
	box-shadow: none;
	display: flex !important;
	gap: 1ch;
	&:hover,
	&:focus-visible {
		background-color: color-mix(in srgb, var(--color-primary), transparent 50%);
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

	&:is(.filter) {
		width: min(100%, 25ch);
		justify-content: center;
		align-items: center;
		margin-left: auto;
	}
}

.show-more {
	width: 100%;
	grid-column: 1 / -1;
	display: flex;
	justify-content: end;
	& button {
		display: flex;
		gap: 1ch;
		background: none;
		padding: 0.5em 1em;
		color: var(--text-color);
		&:hover,
		&:focus-visible {
			background-color: color-mix(in srgb, var(--color-primary), transparent 75%);
			color: var(--text-color) !important;
		}
	}
}

dl {
	display: flex;
	gap: 0.5em;
	align-items: center;
	font-size: 1.2em;
	color: var(--text-color);
	& dt {
		font-weight: bold;
	}
	& dd {
		opacity: 0.8;
	}
}
</style>
