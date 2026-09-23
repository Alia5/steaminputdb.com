<script lang="ts">
import type { ControllerGroup } from './config.types';
import { glyphFor } from './glyphs/glyphFor.svelte';
import { decodeGyroButtons } from './helper';

const {
	group,
	gyroBtn,
	selectedController
}: {
	group: ControllerGroup;
	gyroBtn?: string;
	selectedController: string;
} = $props();
</script>

{#if gyroBtn}
	<span>Choose Gyro Button(s)</span>
	{#each decodeGyroButtons(gyroBtn) as name (name)}
		<div class="gyro-glyph">
			{@render glyphFor(selectedController, name)}
		</div>
	{/each}
	<span>{group.settings?.gyro_button_invert == '1' ? 'Hold to disable Gyro' : 'Hold to enable Gyro'}</span>
{:else if !group.settings?.gyro_ratchet_button_mask}
	<span>Always On</span>
{/if}

<style lang="postcss">
.gyro-glyph {
	display: flex;
	flex-flow: row wrap;
	align-items: center;
	justify-content: center;
}
</style>
