<script lang="ts">
import Icon from '@iconify/svelte';
import type { Action, ActionLayer, ControllerPreset } from './config.types';

let {
	presets,
	actionLayers,
	actions,
	selectedSetName = $bindable()
}: {
	presets: ControllerPreset[];
	actionLayers?: Record<string, ActionLayer>;
	actions?: Record<string, Action>;
	selectedSetName?: string;
} = $props();
</script>

<div class="layers-and-sets">
	<div>
		{#if actions}
			{#each Object.entries(actions) as [key, action] (key)}
				{@const title = action?.title?.startsWith('#') ? key : action?.title}
				<button data-selected={key === selectedSetName} onclick={() => (selectedSetName = key)}
					>{title}</button
				>
				{#each Object.entries(actionLayers || {}).filter(([, v]) => {
					return v.parent_set_name === key;
				}) as [layer_key, layer] (layer_key)}
					<button
						data-selected={layer_key === selectedSetName}
						onclick={() => (selectedSetName = layer_key)}
						><Icon icon="mdi:layers-triple" width="1.2em" />
						{title}:
						{layer.title?.startsWith('#') ? key : layer.title}</button
					>
				{/each}
			{/each}
		{:else}
			{#each presets as preset (preset.name)}
				<button
					data-selected={preset.name === selectedSetName}
					onclick={() => (selectedSetName = preset.name)}
				>
					{actionLayers?.[preset.name]?.title ?? preset.name}
				</button>
			{/each}
		{/if}
	</div>
</div>

<style lang="postcss">
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
				background-color: color-mix(in srgb, var(--color-primary), transparent 10%);
			}
		}
	}
}
</style>
