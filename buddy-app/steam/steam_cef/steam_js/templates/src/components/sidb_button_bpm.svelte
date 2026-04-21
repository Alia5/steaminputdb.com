<script lang="ts" module>
export const ATTR = 'data-sidb-button-bpm';
</script>

<script lang="ts">
import Logo from '$lib/assets/logo.svg?component';

const { ...props }: Record<string, unknown> = $props();

let container = $state<HTMLElement>()!;

const onclick = () => {};
</script>

<div bind:this={container} {...props} data-sidb-button-bpm aria-label="SteamInputDB" {onclick}>
	<Logo />
</div>
