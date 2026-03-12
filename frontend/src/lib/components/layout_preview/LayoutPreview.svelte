<script lang="ts">
import { tooltip } from '$lib/attachments/tooltip.svelte';
import Icon from '@iconify/svelte';
import { onMount } from 'svelte';
import { SvelteSet } from 'svelte/reactivity';
import { parse } from 'vdf-parser';
import Spinner from '../Spinner.svelte';
import type {
	Action,
	ActionLayer,
	ActivatorMap,
	ControllerGroup,
	ControllerPreset,
	LayoutFile
} from './config.types';
import { defaultGyroButton, niceInputMap } from './const';
import { controllerForType } from './controllers/controllerForType.svelte';
import { glyphFor } from './glyphs/glyphFor.svelte';
import { decodeGyroButtons, mergeLayerInputs, niceInputName } from './helper';

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
		selectedController = res.controller_mappings.controller_type || 'controller_generic';
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

const localizeGameAction = (input: string): string | undefined => {
	const parts = input.split(' ');
	const actionName = parts[2];
	if (!actionName) return undefined;

	const english = parsedVdf?.controller_mappings?.localization?.english;
	if (english) {
		const localized = english[actionName];
		if (typeof localized === 'string' && localized) {
			return localized;
		}
	}
	return undefined;
};

const presetTitle = (presetId: string): string | undefined => {
	if (!parsedVdf) {
		return;
	}
	const presets = Array.isArray(parsedVdf.controller_mappings.preset)
		? parsedVdf.controller_mappings.preset
		: [parsedVdf.controller_mappings.preset];
	const preset = presets.find((p) => p.id == presetId);
	if (!preset) {
		return;
	}
	return (
		parsedVdf.controller_mappings.action_layers?.[preset.name]?.title ??
		parsedVdf.controller_mappings.actions?.[preset.name]?.title ??
		undefined
	);
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

const optionalDevices: Record<string, Record<string, boolean | undefined>> = {
	controller_steamcontroller_gordon: {
		trackpads: true,
		'right-stick': false,
		dpad: false
	},
	controller_ps4: {
		trackpads: true
	},
	controller_ps5: {
		trackpads: true
	},
	controller_triton: {
		trackpads: true
	},
	controller_neptune: {
		trackpads: true
	},
	controller_xboxone: {
		gyro: false
	},
	controller_xboxelite: {
		gyro: false
	},
	controller_xbox360: {
		gyro: false
	}
};
</script>

{#snippet LayersAndSets(
	presets: ControllerPreset[],
	actionLayers?: Record<string, ActionLayer>,
	actions?: Record<string, Action>
)}
	<div class="layers-and-sets">
		<div>
			{#if actions}
				{#each Object.entries(actions) as [key, action] (key)}
					{@const title = action?.title?.startsWith('#') ? key : action?.title}
					<button data-selected={key === selectedSetName} onclick={() => (selectedSetName = key)}
						>{title}</button>
					{#each Object.entries(actionLayers || {}).filter(([, v]) => {
						return v.parent_set_name === key;
					}) as [layer_key, layer] (layer_key)}
						<button
							data-selected={layer_key === selectedSetName}
							onclick={() => (selectedSetName = layer_key)}
							><Icon icon="mdi:layers-triple" width="1.2em" />
							{title}:
							{layer.title?.startsWith('#') ? key : layer.title}</button>
					{/each}
				{/each}
			{:else}
				{#each presets as preset (preset.name)}
					<button
						data-selected={preset.name === selectedSetName}
						onclick={() => (selectedSetName = preset.name)}>
						{actionLayers?.[preset.name]?.title ?? preset.name}
					</button>
				{/each}
			{/if}
		</div>
	</div>
{/snippet}

{#snippet gyroPreview(group: ControllerGroup, gyroBtn?: string)}
	{#if gyroBtn}
		<span>Choose Gyro Button(s)</span>
		{#each decodeGyroButtons(gyroBtn) as name (name)}
			<div class="gyro-glyph">
				{@render glyphFor(selectedController, name)}
			</div>
		{/each}
		<span
			>{group.settings?.gyro_button_invert == '1'
				? 'Hold to disable Gyro'
				: 'Hold to enable Gyro'}</span>
	{:else if !group.settings?.gyro_ratchet_button_mask}
		<span>Always On</span>
	{/if}
{/snippet}

{#snippet bindingPreview(name: string, b?: string)}
	{#if b}
		{@const [input, description] = b.split(',').map((s: string) => s.trim())}
		{@const isEmpty = input?.includes('empty_binding')}
		{@const isModeShiftBinding = input?.startsWith('mode_shift ')}
		{@const isHoldLayerBinding = input?.includes('hold_layer')}
		{#if isEmpty}
			<span>
				{#if description}
					{description}
					(Cleared from Parent)
				{:else}
					--
				{/if}
			</span>
		{:else if isModeShiftBinding}
			{@const parts = input?.split(' ') ?? []}
			{@const target = (niceInputMap as Record<string, string>)[parts[1] ?? ''] ?? parts[1]}
			<span>Mode Shift ({target})</span>
		{:else if isHoldLayerBinding}
			{@const parts = input?.split(' ') ?? []}
			{@const layerNum = parts[2]}
			{@const layerTitle = layerNum ? presetTitle(String(Number(layerNum) - 1)) : undefined}
			<span
				>{niceInputName(input)}
				{#if layerTitle}
					({layerTitle})
				{/if}
			</span>
		{:else}
			{@const isGameAction = input?.toLowerCase().includes('game_action')}
			{@const gameActionLabel = isGameAction ? localizeGameAction(input!) : undefined}
			<span>
				{#if isGameAction}
					{gameActionLabel ?? description ?? input?.split(' ')[2] ?? niceInputName(input)}
				{:else if description}
					{description} ({niceInputName(input)})
				{:else if name == 'Full_Press'}
					{niceInputName(input)}
				{:else}
					({name.replace(/_/g, ' ')}) {niceInputName(input)}
				{/if}
			</span>
		{/if}
	{/if}
{/snippet}

{#snippet activatorPreview(activatorMap?: ActivatorMap)}
	{#each Object.entries(activatorMap ?? {}) as [name, activators] (name)}
		{#each Array.isArray(activators) ? activators : [activators] as activator (activator)}
			{#each Object.entries(activator?.bindings ?? {}) as [bindingName, binding] (bindingName)}
				{#each Array.isArray(binding) ? binding : [binding].filter((v) => !!v) as b, bIdx (`${bindingName}_${bIdx}_${b}`)}
					{@render bindingPreview(name, b)}
				{/each}
			{/each}
		{/each}
	{/each}
{/snippet}

{#snippet mappingPreview(prefix: string, device: string, cls: string)}
	{@const leftBumperInputs = getInputs(prefix, device)}
	{#if leftBumperInputs?.length}
		<div class={cls}>
			<div>
				{#each leftBumperInputs as input, idx (idx)}
					{#if input}
						{@render activatorPreview(input.activators)}
					{/if}
				{/each}
			</div>
			{@render glyphFor(selectedController, device)}
		</div>
	{/if}
{/snippet}

{#snippet sourcePreview(source: string, cls: string)}
	{@const groups = getSourceGroups(source)}
	{#if groups?.length}
		<div class={cls}>
			{#if source.includes('trackpad')}
				<span class="device-label">{source === 'left_trackpad' ? 'Left' : 'Right'} Trackpad</span>
			{:else if source === 'dpad'}
				<span class="device-label">D-Pad</span>
			{:else if source === 'joystick'}
				<span class="device-label">Left Joystick</span>
			{:else if source === 'right_joystick'}
				<span class="device-label">Right Joystick</span>
			{:else if source === 'gyro'}
				<span class="device-label">Gyro</span>
			{:else if source === 'button_diamond'}
				<span class="device-label">{(niceInputMap as Record<string, string>)[source] ?? source}</span>
			{/if}

			{#each groups as { group, binding }, idx (idx)}
				{@const isModeshift = binding.includes('modeshift')}
				{@const modeLabel = (niceInputMap as Record<string, string>)[group.mode]}
				{#if isModeshift}
					{@const sourceLabel = (niceInputMap as Record<string, string>)[source] ?? source}
					{@const shiftTriggers = findModeShiftTriggers(source)}
					<span><b>Mode Shift - {modeLabel ?? sourceLabel}</b></span>
					<span>Using:</span>
					<div class="mode-trigger">
						{#if shiftTriggers.length}
							{#each shiftTriggers as trigger (trigger)}
								{@render glyphFor(selectedController, trigger)}
							{/each}
						{:else}
							<span>No Buttons Selected</span>
						{/if}
					</div>
				{/if}
				{#if !isModeshift && modeLabel && group.mode !== 'switches' && group.mode !== 'trigger'}
					<span><b>{modeLabel}</b></span>
				{/if}
				{#if group.settings?.output_trigger}
					<span class="analog"
						>{group.settings.output_trigger == '1' ? 'Left' : 'Right'} Analog Trigger</span>
				{/if}
				{#if group.settings?.output_joystick}
					{@const joystickLabels: Record<string, string> = {
						'0': 'Relative Mouse',
						'1': 'Left Joystick',
						'2': 'Right Joystick',
						'3': 'Relative Mouse',
						'4': 'Mouse Joystick'
					}}
					<span
						>Output: {joystickLabels[group.settings.output_joystick] ??
							`Joystick ${group.settings.output_joystick}`}</span>
				{/if}
				{#if source === 'gyro' && !isModeshift}
					{@render gyroPreview(group, findGyroButton())}
				{/if}
				{#if !isModeshift}
					{#each Object.entries(group.inputs ?? {}) as [device, input] (device)}
						{@const isTriggerSource = source.includes('trigger')}
						{@const deviceLabel =
							device === 'edge'
								? isTriggerSource
									? 'Soft Pull'
									: 'Soft Press'
								: device === 'click'
									? isTriggerSource
										? 'Full Pull'
										: 'Click'
									: ((niceInputMap as Record<string, string>)[device] ?? device)}
						{#if input}
							{@const isMenu = deviceLabel.includes('_menu_button')}
							<div class={'device-group ' + (isMenu ? 'full ' : '')}>
								{#if !isMenu}
									{@render glyphFor(selectedController, `${source}_${device}`)}
								{/if}
								<div class="device-bindings">
									{@render activatorPreview(input.activators)}
								</div>
							</div>
						{/if}
					{/each}
				{/if}
			{/each}
		</div>
	{/if}
{/snippet}

{#snippet presetPreview()}
	<div class="preset-preview">
		<div class="ctrl">
			{@render controllerForType(selectedController)}
		</div>
		{@render mappingPreview('switch', 'left_bumper', 'lb')}
		{@render mappingPreview('switch', 'right_bumper', 'rb')}
		{@render sourcePreview('left_trigger', 'lt')}
		{@render sourcePreview('right_trigger', 'rt')}
		{@render mappingPreview('switch', 'button_back_left_upper', 'l4')}
		{@render mappingPreview('switch', 'button_back_left', 'l5')}
		{@render mappingPreview('switch', 'button_back_right_upper', 'r4')}
		{@render mappingPreview('switch', 'button_back_right', 'r5')}
		{@render mappingPreview('switch', 'button_menu', 'select')}
		{#if optionalDevices[selectedController]?.trackpads === true}
			{@render sourcePreview('left_trackpad', 'lpad')}
			{@render sourcePreview('right_trackpad', 'rpad')}
		{:else}
			<div class="lpad"></div>
			<div class="rpad"></div>
		{/if}
		{@render mappingPreview('switch', 'button_escape', 'start')}
		<div>
			{#if optionalDevices[selectedController]?.dpad !== false}
				{@render sourcePreview('dpad', 'dpad')}
			{/if}
			{@render sourcePreview('joystick', 'lstick')}
			{#if optionalDevices[selectedController]?.gyro !== false}
				{@render sourcePreview('gyro', 'gyro')}
			{/if}
			{#if optionalDevices[selectedController]?.['right-stick'] !== false}
				{@render sourcePreview('right_joystick', 'rstick')}
			{/if}
			{@render sourcePreview('button_diamond', 'buttons')}
		</div>
	</div>
{/snippet}

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
				})}>
				<em>Beta</em>
				<Icon icon="material-symbols:help-outline" width="1.6em" />
			</p>
			<label for="controller-type">
				<span>Controller-Type:</span>
				<select id="controller-type" name="controller-type" bind:value={selectedController}>
					<option value="controller_neptune">Steam Deck</option>
					<option value="controller_triton">Steam Controller</option>
					<option value="controller_steamcontroller_gordon">Steam Controller (2015)</option>
					<option value="controller_ps5">DualSense / DualSense Edge</option>
					<option value="controller_ps4">DualShock 4</option>
					<option value="controller_switch_pro">Switch Pro / 8BitDo</option>
					<option value="controller_xboxone">XBox One / Elite</option>
					<option value="controller_generic">Other</option>
				</select>
				<Icon icon="mdi:chevron-down" />
			</label>
		</div>
		<svelte:boundary>
			{#if parsedVdf}
				{@render LayersAndSets(
					Array.isArray(parsedVdf.controller_mappings.preset)
						? parsedVdf.controller_mappings.preset
						: [parsedVdf.controller_mappings.preset],
					parsedVdf.controller_mappings.action_layers,
					parsedVdf.controller_mappings.actions
				)}
				{#if selectedPreset}
					{@render presetPreview()}
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

	display: grid;
	place-items: center;
	position: relative;
	isolation: isolate;
	& > :first-child {
		width: 100%;
		display: grid;
		padding: 0;
		overflow: auto;
	}
	& > .card {
		position: absolute;
		inset: 1em;
		z-index: -1;
	}
}

.layers-and-sets {
	width: 100%;
	display: flex;
	justify-content: space-evenly;
	overflow-x: auto;
	padding: 1em;

	& div {
		gap: 1em;
		display: flex;
		justify-content: center;
		margin: auto;

		flex-flow: row nowrap;
		& button {
			display: flex;
			gap: 1ch;
			align-items: center;
			white-space: nowrap;
			border-radius: 100vw;
			font-weight: bold;
			padding: 0.5em 1em;
			&:hover,
			&:focus-visible {
				background-color: color-mix(in srgb, var(--color-primary), transparent 50%);
				color: var(--text-color);
			}
			&[data-selected='true'] {
				background-color: color-mix(in srgb, var(--text-color), transparent 85%);
			}
		}
	}
}
.preset-preview {
	width: 100%;
	display: grid;
	padding: 0.5em 1em;
	padding-top: 0em;
	padding-bottom: 1em;
	grid-template-areas:
		'lb      ctrl   rb'
		'lt      ctrl   rt'
		'l4      ctrl   r4'
		'l5      ctrl   r5'
		'select  ctrl   start'
		'lpad    ctrl   rpad'
		'last    last   last';
	grid-template-rows:
		repeat(5, min-content)
		minmax(6em, 1fr)
		minmax(8em, auto);
	grid-template-columns: 1fr minmax(12.5%, 46%) 1fr;
	/* --vertical-space: 0.4em; */
	--vertical-space: 0;
	min-width: 850px;
	@media (orientation: portrait) {
		grid-template-columns: 1fr minmax(10%, 20%) 1fr;
		min-width: 542px;
	}

	& > * {
		padding: 0.2em 0;
	}
	& > :first-child {
		padding-top: 0;
	}
	& > :last-child {
		padding-bottom: 0;
	}

	max-height: 90dvh;

	place-items: center;
	& > *:not(:last-child) {
		display: grid;
		overflow: auto;
		width: 100%;
		height: 100%;
	}
	.ctrl {
		grid-area: ctrl;
		height: 100%;
		width: 100%;
		place-items: center;
		padding: 0 1em;
		& :global(svg) {
			max-height: 100%;
			max-width: 100%;
			opacity: 0.6;
		}
	}
	.lb {
		grid-area: lb;
	}
	.rb {
		grid-area: rb;
	}
	.lt {
		grid-area: lt;
		justify-content: right;
		& > * {
			display: grid;
			padding-right: 1em;
		}
		& > .device-group {
			grid-template-columns: auto min-content;
			gap: 0.5em;
			padding-right: 0;
			align-items: center;
			padding-bottom: var(--vertical-space);
			justify-content: right;
			& :global(> :first-child) {
				grid-column: 2;
				grid-row: 1 / -1;
			}
		}
	}
	.rt {
		grid-area: rt;
		& > * {
			display: grid;
			padding-left: 2em;
		}
		& > .device-group {
			display: grid;
			grid-template-columns: min-content auto;
			gap: 0.5em;
			padding-left: 0;
			align-items: center;
			padding-bottom: var(--vertical-space);
		}
	}
	.l4 {
		grid-area: l4;
	}
	.r4 {
		grid-area: r4;
	}
	.l5 {
		grid-area: l5;
	}
	.r5 {
		grid-area: r5;
	}
	.select {
		grid-area: select;
	}
	.start {
		grid-area: start;
	}

	.lb,
	.rb,
	.l4,
	.r4,
	.l5,
	.r5,
	.select,
	.start {
		display: grid;
		gap: 0.5em;
		align-items: center;
		& > :first-child {
			display: grid;
			/* & span {
				white-space: nowrap;
			} */
		}
	}

	.lb,
	.l4,
	.l5,
	.select {
		justify-items: end;
		grid-template-columns: auto min-content;
		& :global(> :last-child) {
			grid-column: 2;
			grid-row: 1 / -1;
		}
	}
	.rb,
	.r4,
	.r5,
	.start {
		grid-template-columns: min-content auto;
		& :global(> :last-child) {
			grid-column: 1;
			grid-row: 1 / -1;
		}
	}

	.lb,
	.rb,
	.lt,
	.rt,
	.l4,
	.r4,
	.l5,
	.r5,
	.select,
	.start {
		padding-bottom: var(--vertical-space);
	}

	.lpad {
		grid-area: lpad;
		margin-bottom: auto;
		width: 100%;
		justify-content: center;
		grid-auto-rows: min-content;
		gap: 0.1em;

		& .device-bindings {
			justify-content: right;
			width: 100%;
		}
		& .device-group {
			grid-template-columns: auto min-content;
			& :global(> :first-child:not(div)) {
				grid-column: 2;
				grid-row: 1 / -1;
			}
		}
	}
	.rpad {
		grid-area: rpad;
		margin-bottom: auto;
		width: 100%;
		justify-content: center;
		grid-auto-rows: min-content;
		gap: 0.1em;

		& .device-bindings {
			justify-content: left;
			width: 100%;
		}
	}
	& > :last-child {
		grid-area: last;
		display: grid;
		grid-auto-flow: column;
		grid-auto-columns: 1fr;
		place-items: center;
		width: 100%;
		height: 100%;
		overflow: hidden;
		gap: 1em;
		max-height: 18em;
		@media (orientation: portrait) {
			grid-auto-flow: row;
			grid-auto-rows: min-content;
			overflow: auto;
		}
		& > * {
			display: grid;
			height: 100%;
			overflow: auto;
			grid-auto-rows: min-content;
			gap: 0.16em;
		}
	}
	:global(.glyph path),
	:global(.glyph rect),
	:global(.glyph circle),
	:global(.glyph polygon),
	:global(.glyph ellipse) {
		fill: currentColor;
	}
	:global(.color-glyph) {
		filter: drop-shadow(0 0 2px rgb(0 0 0 / 0.8)) drop-shadow(0 0 1px rgb(0 0 0 / 0.2));
	}
}
.device-group {
	display: grid;
	grid-template-columns: auto 1fr;
	gap: 0 0.5em;
	align-items: center;
	&.full {
		grid-template-columns: auto;
	}
	.device-label {
		grid-row: 1 / -1;
		display: grid;
		align-items: center;
	}
	.device-bindings {
		display: grid;
	}
}

.analog {
	text-align: center;
}

.mode-trigger {
	display: flex;
	flex-flow: row wrap;
	justify-content: center;
	align-items: center;
}

.device-label {
	font-size: 1.1em;
	opacity: 0.7;
	font-weight: bold;
	width: 100%;
	text-align: center;
}

.gyro-glyph {
	display: flex;
	flex-flow: row wrap;
	align-items: center;
	justify-content: center;
}

.info {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 1em;
	padding: 1em;
	& > :first-child {
		display: flex;
		align-items: center;
		gap: 0.5em;
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
	border-radius: 0.5em;
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
