<script lang="ts" module>
export { configurationFeatureList };
</script>

<script lang="ts">
import type { components } from '$lib/api/openapi';
import Icon from '@iconify/svelte';
</script>

{#snippet configurationFeatureList({
	fileInfo
}: {
	fileInfo: components['schemas']['ConfigItem' | 'ConfigDetailResponse'];
})}
	{#each fileInfo.tags?.filter((t) => t.startsWith('feature_') && !t.endsWith('_activator')) as tag (tag)}
		<span class="feature">
			{#if tag === 'feature_gamepad'}
				<Icon icon="mdi:controller" width="1.2em" /> Gamepad Inputs
			{:else if tag === 'feature_keyboard' || tag === 'feature_keboard'}
				<Icon icon="mdi:keyboard" width="1.2em" /> Keyboard Inputs
			{:else if tag === 'feature_mouse'}
				<Icon icon="mdi:mouse" width="1.2em" /> Mouse Inputs
			{:else if tag === 'feature_gyro'}
				<Icon icon="game-icons:gyroscope" width="1.2em" /> Gyro Inputs
			{:else if tag === 'feature_touchmenu'}
				<Icon icon="mdi:gesture-touch" width="1.2em" /> Touch Menus
			{:else if tag === 'feature_radialmenu'}
				<Icon icon="material-symbols:joystick" width="1.2em" /> Radial Menus
			{:else if tag === 'feature_modeshift'}
				<Icon icon="mdi:set-right" width="1.2em" /> Mode Shifts
			{:else if tag === 'feature_mouseregion'}
				<Icon icon="fluent:cursor-hover-16-filled" width="1.2em" /> Mouse Regions
			{:else if tag === 'feature_actionset'}
				<Icon icon="material-symbols:layers-rounded" width="1.2em" /> Action Sets
			{:else}
				{tag}
			{/if}
		</span>
	{/each}
{/snippet}

<style lang="postcss">
:global(span.feature) {
	display: flex;
	align-items: center;
	gap: 0.2ch;
	padding: 0.4em 0.6em;
	border-radius: var(--border-radius);
	box-shadow: var(--card-shadow);
	background: var(--card-background-noise);
	font-size: 0.9em;
	font-weight: bold;
}
</style>
