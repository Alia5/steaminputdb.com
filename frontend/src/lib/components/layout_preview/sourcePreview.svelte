<script lang="ts">
import BindingPreview from './bindingPreview.svelte';
import type { ControllerGroup, ControllerPreset, LayoutFile } from './config.types';
import { niceInputMap } from './const';
import { glyphFor } from './glyphs/glyphFor.svelte';
import GyroPreview from './gyroPreview.svelte';

const {
	source,
	cls,
	selectedController,
	parsedVdf,
	getSourceGroups,
	findGyroButton,
	findModeShiftTriggers
}: {
	source: string;
	cls: string;
	selectedController: string;
	parsedVdf?: LayoutFile;
	getSourceGroups: (
		source: string,
		preset?: ControllerPreset
	) => { group: ControllerGroup; binding: string }[];
	findGyroButton: () => string | undefined;
	findModeShiftTriggers: (source: string) => string[];
} = $props();

const groups = $derived(getSourceGroups(source));
</script>

{#if groups?.length}
	<div class={cls} data-gamepadnav-focusable>
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
					>{group.settings.output_trigger == '1' ? 'Left' : 'Right'} Analog Trigger</span
				>
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
						`Joystick ${group.settings.output_joystick}`}</span
				>
			{/if}
			{#if source === 'gyro' && !isModeshift}
				<GyroPreview
					group={group}
					gyroBtn={findGyroButton()}
					selectedController={selectedController}
				/>
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
								{#each Object.entries(input.activators ?? {}) as [name, activators] (name)}
									{#each Array.isArray(activators) ? activators : [activators] as activator (activator)}
										{#each Object.entries(activator?.bindings ?? {}) as [bindingName, bindings] (bindingName)}
											{#each Array.isArray(bindings) ? bindings : [bindings].filter((v) => !!v) as b, bIdx (`${bindingName}_${bIdx}_${b}`)}
												<BindingPreview name={name} b={b} parsedVdf={parsedVdf} />
											{/each}
										{/each}
									{/each}
								{/each}
							</div>
						</div>
					{/if}
				{/each}
			{/if}
		{/each}
	</div>
{/if}

<style lang="postcss">
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
</style>
