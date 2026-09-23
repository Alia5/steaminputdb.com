<script lang="ts">
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { onMount } from 'svelte';
import { SvelteSet } from 'svelte/reactivity';
import { parse } from 'vdf-parser';
import Spinner from '../Spinner.svelte';
import type { ControllerGroup, ControllerPreset, LayoutFile } from './config.types';
import { defaultGyroButton } from './const';
import { mergeLayerInputs } from './helper';
import LayersAndSets from './layersAndSets.svelte';
import PresetPreview from './presetPreview.svelte';

import { browser } from '$app/environment';
import IconHelp from '~icons/material-symbols/help-outline';
import IcoDropDown from '~icons/mdi/chevron-down';
import BPMSelect from '../BPM_Select/BPM_Select.svelte';
import BPMOption from '../BPM_Select/BPM_option.svelte';

const {
	vdfLink
}: {
	vdfLink?: string;
} = $props();

let vdfPromise = $state<Promise<LayoutFile | undefined>>();
let parsedVdf = $derived(await vdfPromise);
onMount(() => {
	vdfPromise = (async () => {
		// await new Promise((resolve) => setTimeout(resolve, 5000));
		if (!vdfLink) {
			return;
		}
		const vdf = await fetch(vdfLink).then((resp) => resp.text());
		const res = parse<LayoutFile>(vdf);
		selectedSetName =
			Object.keys(res.controller_mappings.actions ?? {})[0] ??
			(Array.isArray(res.controller_mappings.preset)
				? res.controller_mappings.preset[0]
				: res.controller_mappings.preset
			)?.name;
		const controllerTypeAlias: Record<string, string> = {
			controller_switch2_pro: 'controller_switch_pro',
			controller_8bitdo: 'controller_switch_pro',
			controller_ps5_edge: 'controller_ps5',
			controller_xboxelite: 'controller_xboxone'
		};
		const rawType = res.controller_mappings.controller_type || 'controller_generic';
		selectedController = controllerTypeAlias[rawType] ?? rawType;
		return res;
	})();
});

let selectedSetName = $state<string | undefined>('Default');
let selectedPreset = $derived.by(() => {
	if (!selectedSetName || !parsedVdf) {
		return;
	}
	return Array.isArray(parsedVdf.controller_mappings.preset)
		? parsedVdf.controller_mappings.preset.find((preset) => preset.name === selectedSetName)
		: parsedVdf.controller_mappings.preset;
});
let selectedController = $state<string>('controller_generic');

let isSteamBigPicture = $derived(browser ? navigator?.userAgent?.includes('Steam Gamepad') : false);

const getAllGroupsForSource = (source: string, preset: ControllerPreset | undefined) => {
	if (!preset || !parsedVdf) {
		return [];
	}
	const allGroups = Array.isArray(parsedVdf.controller_mappings.group)
		? parsedVdf.controller_mappings.group
		: [parsedVdf.controller_mappings.group];
	return Object.entries(preset.group_source_bindings)
		.filter(([, binding]) => binding.startsWith(source))
		.map(([groupId, binding]) => {
			let group = allGroups.find((g) => g.id == groupId);
			if (group?.mode === 'reference' && group.settings?.referenced_mode) {
				group = allGroups.find((g) => g.id == group!.settings!.referenced_mode);
			}
			return { group, binding };
		})
		.filter((g): g is { group: ControllerGroup; binding: string } => !!g.group);
};

const resolveGroups = (source: string, preset: ControllerPreset | undefined) => {
	return getAllGroupsForSource(source, preset).filter(({ binding }) => !binding.includes('inactive'));
};

const findGyroButton = (): string | undefined => {
	const allPresets = !parsedVdf
		? []
		: Array.isArray(parsedVdf.controller_mappings.preset)
			? parsedVdf.controller_mappings.preset
			: [parsedVdf.controller_mappings.preset];
	const match = allPresets
		.flatMap((p) => getAllGroupsForSource('gyro', p))
		.map(({ group }) => group.settings?.gyro_ratchet_button_mask ?? group.settings?.gyro_button)
		.find((val) => val != null && val != undefined);
	if (match != null && match != undefined) {
		return String(match) === '0' ? defaultGyroButton[selectedController] : String(match);
	}
	const hasInactiveGyro = allPresets
		.flatMap((p) => getAllGroupsForSource('gyro', p))
		.some(({ binding }) => binding.includes('inactive'));
	return hasInactiveGyro ? undefined : defaultGyroButton[selectedController];
};

const getParentPreset = (preset: ControllerPreset) => {
	if (!parsedVdf) {
		return;
	}
	const parentName = parsedVdf.controller_mappings.action_layers?.[preset.name]?.parent_set_name;
	if (!parentName) {
		return;
	}
	return (
		Array.isArray(parsedVdf.controller_mappings.preset)
			? parsedVdf.controller_mappings.preset
			: [parsedVdf.controller_mappings.preset]
	).find((p) => p.name === parentName);
};

const getInputs = (group: string, device: string, preset = selectedPreset) => {
	if (!preset) {
		return [];
	}
	const groups = resolveGroups(group, preset);
	const result = groups.map((g) => g.group?.inputs?.[device]).filter(Boolean);
	if (result.length) {
		return result;
	}
	const parentPreset = getParentPreset(preset);
	if (!parentPreset) {
		return result;
	}
	return getInputs(group, device, parentPreset);
};

const getSourceGroups = (source: string, preset = selectedPreset) => {
	if (!preset) {
		return [];
	}
	const result = resolveGroups(source, preset);
	const hasInputs = result.some((g) => Object.keys(g.group.inputs ?? {}).length > 0);
	if (hasInputs) {
		const parentPreset = getParentPreset(preset);
		const parentGroups = parentPreset ? resolveGroups(source, parentPreset) : [];
		const merged = parentGroups.length ? mergeLayerInputs(result, parentGroups) : result;
		return merged.toSorted((a, b) => {
			const aShift = a.binding.includes('modeshift') ? 1 : 0;
			const bShift = b.binding.includes('modeshift') ? 1 : 0;
			return aShift - bShift;
		});
	}
	const parentPreset = getParentPreset(preset);
	if (!parentPreset) {
		return result;
	}
	const parentGroups = resolveGroups(source, parentPreset);
	if (!parentGroups.length) {
		return result.toSorted((a, b) => {
			const aShift = a.binding.includes('modeshift') ? 1 : 0;
			const bShift = b.binding.includes('modeshift') ? 1 : 0;
			return aShift - bShift;
		});
	}
	const layerActiveMode = result.find((g) => !g.binding.includes('modeshift'))?.group.mode;
	const parentActiveMode = parentGroups.find((g) => !g.binding.includes('modeshift'))?.group.mode;
	if (layerActiveMode && parentActiveMode && layerActiveMode !== parentActiveMode) {
		return result.toSorted((a, b) => {
			const aShift = a.binding.includes('modeshift') ? 1 : 0;
			const bShift = b.binding.includes('modeshift') ? 1 : 0;
			return aShift - bShift;
		});
	}
	return getSourceGroups(source, parentPreset);
};

const findModeShiftTriggers = (source: string): string[] => {
	if (!parsedVdf || !selectedPreset) {
		return [];
	}
	const switchGroups = getSourceGroups('switch');
	const sourcePattern = `mode_shift ${source} `;

	const triggers = new SvelteSet<string>();
	for (const { group } of switchGroups) {
		for (const [inputName, input] of Object.entries(group.inputs ?? {})) {
			for (const acts of Object.values(input?.activators ?? {})) {
				for (const activator of Array.isArray(acts) ? acts : [acts]) {
					for (const b of Object.values(activator?.bindings ?? {})
						.flatMap((b) => (Array.isArray(b) ? b : [b]))
						.filter((b): b is string => !!b)) {
						if (b.startsWith(sourcePattern)) {
							triggers.add(inputName);
						}
					}
				}
			}
		}
	}
	return [...triggers];
};
</script>

<section>
	<div>
		<div class="info">
			<p
				{@attach tooltip({
					content:
						'This Layout preview is not throughly tested\nNo guarantees for correctness or completeness!',
					outDelay: 200,
					arrow: true,

					arrowFollowCursor: true
				})}
			>
				<em>Beta</em>
				<IconHelp style="width: 1.6em; height: 1.6em;" />
			</p>
			{#if isSteamBigPicture}
				<BPMSelect name="Controller-Type" bind:value={selectedController}>
					{#snippet children({ ...rest })}
						<span>Controller-Type:</span>
						<BPMOption value="controller_neptune" {...rest}>Steam Deck</BPMOption>
						<BPMOption value="controller_triton" {...rest}>Steam Controller</BPMOption>
						<BPMOption value="controller_steamcontroller_gordon" {...rest}
							>Steam Controller (2015)</BPMOption
						>
						<BPMOption value="controller_ps5" {...rest}>DualSense / DualSense Edge</BPMOption>
						<BPMOption value="controller_ps4" {...rest}>DualShock 4</BPMOption>
						<BPMOption value="controller_switch_pro" {...rest}
							>Switch (1/2) Pro / 8BitDo</BPMOption
						>
						<BPMOption value="controller_xboxone" {...rest}>XBox One / Elite</BPMOption>
						<BPMOption value="controller_generic" {...rest}>Other</BPMOption>
						<IcoDropDown />
					{/snippet}
				</BPMSelect>
			{:else}
				<label for="controller-type">
					<span>Controller-Type:</span>
					<select id="controller-type" name="controller-type" bind:value={selectedController}>
						<option value="controller_neptune">Steam Deck</option>
						<option value="controller_triton">Steam Controller</option>
						<option value="controller_steamcontroller_gordon">Steam Controller (2015)</option>
						<option value="controller_ps5">DualSense / DualSense Edge</option>
						<option value="controller_ps4">DualShock 4</option>
						<option value="controller_switch_pro">Switch (1/2) Pro / 8BitDo</option>
						<option value="controller_xboxone">XBox One / Elite</option>
						<option value="controller_generic">Other</option>
					</select>
					<IcoDropDown />
				</label>
			{/if}
			<span>Best viewed on Desktop (Landscape)</span>
		</div>
		<svelte:boundary>
			{#if parsedVdf}
				<LayersAndSets
					presets={Array.isArray(parsedVdf.controller_mappings.preset)
						? parsedVdf.controller_mappings.preset
						: [parsedVdf.controller_mappings.preset]}
					actionLayers={parsedVdf.controller_mappings.action_layers}
					actions={parsedVdf.controller_mappings.actions}
					bind:selectedSetName={selectedSetName}
				/>
				{#if selectedPreset}
					<PresetPreview
						selectedController={selectedController}
						parsedVdf={parsedVdf}
						getInputs={getInputs}
						getSourceGroups={getSourceGroups}
						findGyroButton={findGyroButton}
						findModeShiftTriggers={findModeShiftTriggers}
					/>
				{:else}
					<div class="no-preview">
						<p>No Layout selected or found</p>
					</div>
				{/if}
			{:else}
				<div class="no-preview">
					<Spinner size="min(75dvw, 12em)" />
				</div>
			{/if}
			{#snippet pending()}
				<div class="no-preview">
					<Spinner size="min(75dvw, 12em)" />
				</div>
			{/snippet}
			{#snippet failed()}
				<div class="no-preview">
					<p>Whoops, seems we were not able to parse the config file...</p>
					<!-- <p>{error}</p> -->
				</div>
			{/snippet}
		</svelte:boundary>
	</div>
	<div class="card glass"></div>
</section>

<style lang="postcss">
section {
	width: 100%;
	padding: 1em;
	padding-top: 0;

	@media (orientation: landscape) and (max-width: 1300px) {
		font-size: 0.78em;
		& :global(svg:not(.ctrl > svg)) {
			width: 1.6em;
			height: 1.6em;
		}
	}

	display: grid;
	place-items: center;
	position: relative;
	isolation: isolate;
	grid-template-columns: 1fr;
	& > :first-child {
		width: 100%;
		display: grid;
		padding: 0;
		overflow: auto;
	}
	& > .card {
		position: absolute;
		inset: 0 1em 0 1em;
		z-index: -1;
	}

	:global(.bpm-option) {
		margin-left: auto;
		font-size: 1.2em;
	}
}

.info {
	display: grid;
	grid-template-columns: auto auto;
	justify-content: space-between;
	align-items: center;
	columns-gap: 1em;
	row-gap: 0.5em;
	padding: 1em;
	& > :first-child {
		display: flex;
		align-items: center;
		gap: 0.5em;
	}

	& > :last-child:is(span) {
		grid-column: 1 / -1;
		grid-row: 2;
		@media (orientation: landscape) {
			display: none;
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

label[for='controller-type'] {
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

.no-preview {
	display: grid;
	padding: 2em 1em;
	place-items: center;
	width: 100%;
	overflow: hidden;
	grid-row: 2;
}
</style>
