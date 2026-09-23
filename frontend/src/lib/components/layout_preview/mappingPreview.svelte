<script lang="ts">
import BindingPreview from './bindingPreview.svelte';
import type { ControllerPreset, GroupInput, LayoutFile } from './config.types';
import { glyphFor } from './glyphs/glyphFor.svelte';

const {
	prefix,
	device,
	cls,
	selectedController,
	parsedVdf,
	getInputs
}: {
	prefix: string;
	device: string;
	cls: string;
	selectedController: string;
	parsedVdf?: LayoutFile;
	getInputs: (group: string, device: string, preset?: ControllerPreset) => (GroupInput | undefined)[];
} = $props();

const leftBumperInputs = $derived(getInputs(prefix, device));
</script>

{#if leftBumperInputs?.length}
	<div class={cls} data-gamepadnav-focusable>
		<div>
			{#each leftBumperInputs as input, idx (idx)}
				{#if input}
					{#each Object.entries(input.activators ?? {}) as [name, activators] (name)}
						{#each Array.isArray(activators) ? activators : [activators] as activator (activator)}
							{#each Object.entries(activator?.bindings ?? {}) as [bindingName, binding] (bindingName)}
								{#each Array.isArray(binding) ? binding : [binding].filter((v) => !!v) as b, bIdx (`${bindingName}_${bIdx}_${b}`)}
									<BindingPreview name={name} b={b} parsedVdf={parsedVdf} />
								{/each}
							{/each}
						{/each}
					{/each}
				{/if}
			{/each}
		</div>
		{@render glyphFor(selectedController, device)}
	</div>
{/if}

<style lang="postcss">
.lb {
	grid-area: lb;
}
.rb {
	grid-area: rb;
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
.capsense_left {
	grid-area: capsense_left;
}
.capsense_right {
	grid-area: capsense_right;
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
.capsense_left,
.capsense_right,
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
.capsense_left,
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
.capsense_right,
.start {
	grid-template-columns: min-content auto;
	& :global(> :last-child) {
		grid-column: 1;
		grid-row: 1 / -1;
	}
}
</style>
