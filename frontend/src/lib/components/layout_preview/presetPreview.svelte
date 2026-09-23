<script lang="ts">
import type { ControllerGroup, ControllerPreset, GroupInput, LayoutFile } from './config.types';
import { controllerForType } from './controllers/controllerForType.svelte';
import MappingPreview from './mappingPreview.svelte';
import SourcePreview from './sourcePreview.svelte';

const {
	selectedController,
	parsedVdf,
	getInputs,
	getSourceGroups,
	findGyroButton,
	findModeShiftTriggers
}: {
	selectedController: string;
	parsedVdf?: LayoutFile;
	getInputs: (group: string, device: string, preset?: ControllerPreset) => (GroupInput | undefined)[];
	getSourceGroups: (
		source: string,
		preset?: ControllerPreset
	) => { group: ControllerGroup; binding: string }[];
	findGyroButton: () => string | undefined;
	findModeShiftTriggers: (source: string) => string[];
} = $props();

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
		trackpads: true,
		capsense: true
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

<div class="preset-preview">
	<div class="ctrl">
		{@render controllerForType(selectedController)}
	</div>
	<MappingPreview
		prefix="switch"
		device="left_bumper"
		cls="lb"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getInputs={getInputs}
	/>
	<MappingPreview
		prefix="switch"
		device="right_bumper"
		cls="rb"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getInputs={getInputs}
	/>
	<SourcePreview
		source="left_trigger"
		cls="lt"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getSourceGroups={getSourceGroups}
		findGyroButton={findGyroButton}
		findModeShiftTriggers={findModeShiftTriggers}
	/>
	<SourcePreview
		source="right_trigger"
		cls="rt"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getSourceGroups={getSourceGroups}
		findGyroButton={findGyroButton}
		findModeShiftTriggers={findModeShiftTriggers}
	/>
	<MappingPreview
		prefix="switch"
		device="button_back_left_upper"
		cls="l4"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getInputs={getInputs}
	/>
	<MappingPreview
		prefix="switch"
		device="button_back_left"
		cls="l5"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getInputs={getInputs}
	/>
	<MappingPreview
		prefix="switch"
		device="button_back_right_upper"
		cls="r4"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getInputs={getInputs}
	/>
	<MappingPreview
		prefix="switch"
		device="button_back_right"
		cls="r5"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getInputs={getInputs}
	/>
	{#if optionalDevices[selectedController]?.capsense === true}
		<MappingPreview
			prefix="switch"
			device="button_leftauxcapsense"
			cls="capsense_left"
			selectedController={selectedController}
			parsedVdf={parsedVdf}
			getInputs={getInputs}
		/>
		<MappingPreview
			prefix="switch"
			device="button_rightauxcapsense"
			cls="capsense_right"
			selectedController={selectedController}
			parsedVdf={parsedVdf}
			getInputs={getInputs}
		/>
	{/if}
	<MappingPreview
		prefix="switch"
		device="button_menu"
		cls="select"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getInputs={getInputs}
	/>
	{#if optionalDevices[selectedController]?.trackpads === true}
		<SourcePreview
			source="left_trackpad"
			cls="lpad"
			selectedController={selectedController}
			parsedVdf={parsedVdf}
			getSourceGroups={getSourceGroups}
			findGyroButton={findGyroButton}
			findModeShiftTriggers={findModeShiftTriggers}
		/>
		<SourcePreview
			source="right_trackpad"
			cls="rpad"
			selectedController={selectedController}
			parsedVdf={parsedVdf}
			getSourceGroups={getSourceGroups}
			findGyroButton={findGyroButton}
			findModeShiftTriggers={findModeShiftTriggers}
		/>
	{:else}
		<div class="lpad"></div>
		<div class="rpad"></div>
	{/if}
	<MappingPreview
		prefix="switch"
		device="button_escape"
		cls="start"
		selectedController={selectedController}
		parsedVdf={parsedVdf}
		getInputs={getInputs}
	/>
	<div>
		<div class="group">
			{#if optionalDevices[selectedController]?.dpad !== false}
				<SourcePreview
					source="dpad"
					cls="dpad"
					selectedController={selectedController}
					parsedVdf={parsedVdf}
					getSourceGroups={getSourceGroups}
					findGyroButton={findGyroButton}
					findModeShiftTriggers={findModeShiftTriggers}
				/>
			{/if}
			<SourcePreview
				source="joystick"
				cls="lstick"
				selectedController={selectedController}
				parsedVdf={parsedVdf}
				getSourceGroups={getSourceGroups}
				findGyroButton={findGyroButton}
				findModeShiftTriggers={findModeShiftTriggers}
			/>
		</div>
		{#if optionalDevices[selectedController]?.gyro !== false}
			<div class="group">
				<SourcePreview
					source="gyro"
					cls="gyro"
					selectedController={selectedController}
					parsedVdf={parsedVdf}
					getSourceGroups={getSourceGroups}
					findGyroButton={findGyroButton}
					findModeShiftTriggers={findModeShiftTriggers}
				/>
			</div>
		{/if}
		<div class="group">
			{#if optionalDevices[selectedController]?.['right-stick'] !== false}
				<SourcePreview
					source="right_joystick"
					cls="rstick"
					selectedController={selectedController}
					parsedVdf={parsedVdf}
					getSourceGroups={getSourceGroups}
					findGyroButton={findGyroButton}
					findModeShiftTriggers={findModeShiftTriggers}
				/>
			{/if}
			<SourcePreview
				source="button_diamond"
				cls="buttons"
				selectedController={selectedController}
				parsedVdf={parsedVdf}
				getSourceGroups={getSourceGroups}
				findGyroButton={findGyroButton}
				findModeShiftTriggers={findModeShiftTriggers}
			/>
		</div>
	</div>
</div>

<style lang="postcss">
.preset-preview {
	width: 100%;
	display: grid;
	padding: 1em;
	padding-top: 2em;
	grid-template-areas:
		'lb                 ctrl    rb'
		'lt                 ctrl    rt'
		'l4                 ctrl    r4'
		'l5                 ctrl    r5'
		'capsense_left      ctrl    capsense_right'
		'select             ctrl    start'
		'lpad               ctrl    rpad'
		'last               last    last';
	grid-template-rows:
		repeat(5, min-content)
		minmax(6em, 1fr)
		minmax(8em, auto);
	grid-template-columns: 1fr minmax(15.5%, 46%) 1fr;
	/* --vertical-space: 0.4em; */
	--vertical-space: 0;
	min-width: 650px;
	max-height: 90dvh;
	@media (orientation: portrait) {
		grid-template-columns: 1fr minmax(10%, 20%) 1fr;
		min-width: 542px;
		grid-template-rows:
			repeat(5, min-content)
			minmax(6em, 1fr)
			fit-content;
		max-height: unset;
	}
	@media (orientation: landscape) and (max-width: 1300px) {
		grid-template-columns: auto minmax(15.5%, 20%) auto;
	}

	& > :global(*) {
		padding: 0.2em 0;
	}
	& > :first-child {
		padding-top: 0;
	}
	& > :last-child {
		padding-bottom: 0;
	}

	place-items: center;
	& > :global(*:not(:last-child)) {
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
		overflow: clip;
		overflow-clip-margin: 2em;
		& :global(svg) {
			max-height: 95%;
			max-width: 95%;
			opacity: 0.6;
			overflow: clip;
			overflow-clip-margin: 2em;
		}
	}

	/* placeholders rendered in place of the optional trackpads */
	.lpad {
		grid-area: lpad;
		margin-bottom: auto;
		width: 100%;
		justify-content: center;
		grid-auto-rows: min-content;
		gap: 0.1em;
	}
	.rpad {
		grid-area: rpad;
		margin-bottom: auto;
		width: 100%;
		justify-content: center;
		grid-auto-rows: min-content;
		gap: 0.1em;
	}

	:global(.lb),
	:global(.rb),
	:global(.lt),
	:global(.rt),
	:global(.l4),
	:global(.r4),
	:global(.l5),
	:global(.r5),
	:global(.capsense_left),
	:global(.capsense_right),
	:global(.select),
	:global(.start) {
		padding-bottom: var(--vertical-space);
	}

	& > :last-child {
		grid-area: last;
		display: grid;
		grid-auto-flow: column;
		grid-auto-columns: 1fr;
		place-items: center;
		width: 100%;
		height: 100%;
		overflow: auto;
		max-height: 18em;

		/* @media (orientation: landscape) and (max-width: 1300px) {
			display: flex;
			gap: 0.5em;
			flex-flow: row nowrap;
			margin: 0;
			justify-content: space-between;
		} */
		display: flex;
		gap: 1.5em;
		flex-flow: row nowrap;
		margin: 0;
		justify-content: space-evenly;
		& > :global(*) {
			display: grid;
			height: 100%;
			overflow: auto;
			grid-auto-rows: min-content;
			gap: 0.2em;
			& :global(span:has(b)) {
				text-align: center;
			}
			& :global(* > span) {
				width: 100%;
				text-align: center;
			}
			& :global(* > div) {
				width: 100%;
				text-align: initial !important;
				& :global(span) {
					width: 100%;
					text-align: initial;
				}
			}
		}
		& .group {
			display: flex;
			flex-flow: row wrap;
			justify-content: space-between;
			flex: 1 0 auto;
			gap: 1.5em;
			& > :global(*) {
				display: grid;
				height: 100%;
				flex: 1 1 auto;
				overflow: auto;
				grid-auto-rows: min-content;
				gap: 0.16em;
			}
		}

		@media (orientation: portrait) {
			max-height: unset;
			display: flex;
			flex-flow: row wrap;
			gap: 2em;
			width: 100%;
			overflow: auto;
			align-items: baseline;
			justify-content: space-evenly;
			& > :global(*) {
				overflow: hidden;
				width: auto;
				height: fit-content;
			}
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
</style>
