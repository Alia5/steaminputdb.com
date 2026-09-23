<script lang="ts">
import type { HTMLInputAttributes } from 'svelte/elements';
import IcOutlineSearch from '~icons/ic/outline-search';
let {
	value = $bindable(),
	placeholder,
	inlineButton = true,
	...props
}: {
	value?: unknown;
	placeholder: string;
	inlineButton?: boolean;
	'--box-shadow'?: string;
	'--font-size'?: string;
} & HTMLInputAttributes = $props();
</script>

<label for="config-search">
	<input
		id="config-search"
		name="config-search"
		type="search"
		bind:value={value}
		{...props}
		placeholder={placeholder}
	/>
	{#if inlineButton}
		<button class="plain" type="submit">
			<IcOutlineSearch />
		</button>
	{:else}
		<IcOutlineSearch />
	{/if}
</label>

<style lang="postcss">
label {
	--font-size: 1.4em;
	--box-shadow-default: 0 1px 4px 0 var(--shadow-color);

	position: relative;
	display: grid;
	grid-template-columns: auto min-content;
	place-items: center;

	border: 1px solid color-mix(in srgb, var(--text-color), transparent 90%);
	outline: 1px solid transparent;
	color: var(--text-color);
	border-radius: 100em;
	box-shadow: var(--box-shadow, var(--box-shadow-default));
	transition: all var(--transition-duration) var(--default-ease);

	&:has([disabled]) {
		opacity: 0.5;
	}

	&:hover,
	&:focus-within {
		outline: 0.1em solid var(--color-primary);
		& :global(svg) {
			opacity: 1;
		}
		&::before {
			box-shadow: 0 0 1.3em -0.4em var(--color-primary);
		}
	}

	background: var(--card-background-noise);
	isolation: isolate;
	&::before {
		content: '';
		position: absolute;
		inset: 0;
		z-index: -1;
		background: var(--background-neutral-alpha);
		border-radius: 100em;
	}

	:global(:last-child:is(svg)) {
		width: var(--font-size);
		height: var(--font-size);
		margin-right: 1em;
		opacity: 0.5;
	}
}
input[type='search'] {
	font-size: var(--font-size);
	width: 100%;
	border: none;
	box-shadow: none;
	padding-right: 0.5em;
	padding-left: 1em;

	outline: none;
	&:is(:hover, :focus) {
		outline: none;
	}
}

button {
	padding: 0;
	outline: 1px solid transparent;
	box-shadow: none;
	border: none;
	height: 100%;
	display: grid;
	place-items: center;
	transition: none;

	& > :global(svg) {
		width: var(--font-size);
		height: var(--font-size);
		margin-right: 1em;
		opacity: 0.5;
	}
	&:hover,
	&:focus-visible {
		color: var(--highlight-color);
		filter: drop-shadow(0 0 0.2em var(--color-primary));
	}
}
</style>
