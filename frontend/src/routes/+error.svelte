<script lang="ts">
import { page } from '$app/state';
import SC2 from '$lib/assets/SC2_Googley.svg.svelte';
import { statusCodeNames } from '$lib/statuscodes';
import { onMount } from 'svelte';

let eyes = $state<{
	left: HTMLElement;
	right: HTMLElement;
}>({
	left: undefined!,
	right: undefined!
})!;

onMount(() => {
	const group = document.querySelector('#sc2>*>svg>g>:last-child')!;
	eyes.left = group.children[0] as HTMLElement;
	eyes.right = group.children[1] as HTMLElement;
});
</script>

<svelte:window
	onmousemove={(e) => {
		if (!eyes || !eyes.left || !eyes.right) {
			return;
		}

		Object.values(eyes).forEach((eye, idx) => {
			const rect = eye.getBoundingClientRect();
			const eyeX = rect.left + rect.width / 2;
			const eyeY = rect.top + rect.height / 2;
			const deltaX = e.clientX - eyeX;
			const deltaY = e.clientY - eyeY;
			const angle = Math.atan2(deltaY, deltaX);
			const distance = Math.min(16, Math.hypot(deltaX, deltaY) / 2);
			const translateX = Math.cos(angle) * distance * (idx === 0 ? 1 : -1);
			const translateY = Math.sin(angle) * distance * (idx === 0 ? -1 : 1);
			eye.style.transform = `translate(${translateX}px, ${translateY}px)`;
		});
	}}
/>

<main>
	<div>
		<h1>{page.status}</h1>
		{#if statusCodeNames[page.status] && statusCodeNames[page.status] !== page.error?.message}
			<h3>{statusCodeNames[page.status]}</h3>
		{/if}
		<h2>{page.error?.message}</h2>
		<div id="sc2">
			<SC2
				height="100%"
				--eyes-color="black"
				--eyes-white-color="var(--text-color-dark)"
				--eyes-border-color="light-dark(var(--text-color-light), transparent)"
			/>
		</div>
	</div>
</main>

<style lang="postcss">
main {
	display: grid;
	place-items: center;
}

div {
	display: grid;
	place-items: center;
	width: 100%;
}

h1 {
	font-size: 6em;
	font-weight: bold;
}

h2 {
	color: var(--highlight-color);
	font-size: 2em;
}

h3 {
	translate: 0 -1em;
}

#sc2 {
	filter: drop-shadow(0px -0.25em 0.2em var(--shadow-color))
		drop-shadow(0px -0.25em 0.2em var(--shadow-color));
	max-height: 16em;
	height: 100%;
	width: 100%;
	display: grid;
	place-items: center;
	padding: 1em;
	@media (any-pointer: coarse) {
		:global(ellipse) {
			transition: transform calc(var(--transition-duration) * 2) var(--default-ease);
		}
	}

	transform: rotate(180deg);

	:global(#qam),
	:global(#pad_r),
	:global(#y),
	:global(#x) {
		display: none;
	}

	:global(#logo_overlay) {
		fill: var(--background-color);
	}

	:global(#dpad) {
		transform: rotate(3deg);
		transform-origin: 25% 25%;
		translate: -0.1em 0.1em;
	}
}
</style>
