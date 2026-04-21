<script lang="ts">
import { onMount, type Snippet } from 'svelte';

let {
	value,
	selected = $bindable(),
	showButton = false,
	popovertarget = 'bpm-select-',
	onselected,
	children
}: {
	selected?: unknown;
	value?: unknown;
	showButton?: boolean;
	popovertarget?: string;
	onselected?: (value: unknown) => void;
	children?: Snippet;
} = $props();

const scrollIntoView = () => {
	if (selected == value && showButton && buttonEl) {
		buttonEl!.scrollIntoView({ behavior: 'instant', block: 'center' });
	}
};

let buttonEl = $state<HTMLButtonElement>()!;

onMount(() => {
	const popover = buttonEl?.closest('[popover]');
	if (!popover) return;
	const onToggle = (e: Event) => {
		if ((e as ToggleEvent).newState === 'open') {
			scrollIntoView();
		}
	};
	popover.addEventListener('toggle', onToggle);
	return () => popover.removeEventListener('toggle', onToggle);
});
</script>

{#if selected === value && selected !== undefined && !showButton}
	<div class="bpm-option">
		{@render children?.()}
	</div>
{/if}
{#if showButton}
	<button
		type="button"
		class={'bpm-option' + (selected === value ? ' selected' : '')}
		autofocus={selected === value}
		onclick={() => {
			selected = value;
			onselected?.(value);
		}}
		popovertarget={popovertarget}
		bind:this={buttonEl}
		popovertargetaction="hide">
		{@render children?.()}
	</button>
{/if}

<style lang="postcss">
button {
	border-radius: 0;
}
.selected {
	font-weight: bold;
	background-color: var(--color-primary);
	&:hover,
	&:focus-visible {
		background-color: color-mix(in srgb, var(--color-primary), var(--text-color) 75%);
	}
}
</style>
