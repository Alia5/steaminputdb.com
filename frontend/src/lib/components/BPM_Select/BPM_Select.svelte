<script lang="ts">
import { type Snippet } from 'svelte';
import { slide } from 'svelte/transition';

let {
	disabled = false,
	value = $bindable(),
	name,
	onchange,
	children
}: {
	disabled?: boolean;
	value?: unknown;
	name?: string;
	onchange?: (value: unknown) => void;
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	children?: Snippet<[{ [x: string]: any }]>;
} = $props();
</script>

<button
	class={'bpm-select' + (disabled ? ' disabled' : '')}
	disabled={disabled}
	type="button"
	popovertarget={'bpm-select-' + name}>
	{@render children?.({ selected: value })}
</button>
<dialog id={'bpm-select-' + name} popover closedby="any">
	<div class="card">
		<h3>Select {name ?? 'option'}</h3>
		<div transition:slide|global class="only-options">
			{@render children?.({
				selected: value,
				showButton: true,
				popovertarget: 'bpm-select-' + name,
				onselected: (v: unknown) => {
					value = v;
					onchange?.(value);
				}
			})}
		</div>
	</div>
</dialog>

<style lang="postcss">
.only-options {
	:global(> *) {
		display: none;
	}
	:global(.bpm-option) {
		display: block !important;
	}
	display: grid;
	overflow: auto;
	min-height: 0;
}

.disabled {
	opacity: 0.5;
}

.bpm-select {
	font-style: inherit;
	background: transparent;
	border: 1px solid color-mix(in srgb, var(--text-color), transparent 90%);
	box-shadow: 0 1px 4px 0 rgb(0 0 0 / 0.25);
	border-radius: 0.5em;
	color: var(--text-color);
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
	gap: 0.5em;

	& :global(.bpm-option) {
		color: var(--text-color);
		background: var(--card-color);
	}
	&::before {
		opacity: 0;
	}

	&:hover:not(:disabled),
	&:focus-visible:not(:disabled) {
		outline: 2px solid var(--color-primary);
		box-shadow: 0 0 1.3em -0.4em var(--color-primary);
		color: var(--color-primary);
		& :global(> *) {
			color: var(--color-primary);
		}
	}
}

dialog {
	padding: 0;
	position: fixed;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);

	outline: 1px solid transparent;
	border: none;
	background: transparent;

	opacity: 0;
	transition: opacity var(--transition-duration) var(--default-ease) allow-discrete;
	animation: fade-in var(--transition-duration) var(--default-ease) forwards;

	&:popover-open {
		opacity: 1;
	}

	&::backdrop {
		opacity: 0;
		transition: opacity var(--transition-duration) var(--default-ease) allow-discrete;
		animation: fade-in var(--transition-duration) var(--default-ease) forwards;
		background: var(
			--background,
			linear-gradient(
				color-mix(in srgb, transparent, rgb(32, 25, 47) var(--background-opacity, 40%)),
				color-mix(in srgb, transparent, rgb(7, 4, 11) var(--background-opacity, 55%))
			)
		);
	}

	isolation: isolate;
	overflow: visible;

	& > :first-child {
		max-height: calc(100svh - 2em);
		overflow: hidden;
		display: flex;
		flex-direction: column;
		padding: 1em 0;
		padding-bottom: 0;
	}
	h3 {
		margin-bottom: 1em;
		flex-shrink: 0;
		padding: 0 1em;
	}
}

:global(body:has(dialog[popover]:popover-open)) {
	overflow: hidden;
}
</style>
