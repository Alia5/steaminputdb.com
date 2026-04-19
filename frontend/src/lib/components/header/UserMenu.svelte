<script lang="ts">
import { resolve } from '$app/paths';
import { page } from '$app/state';
const userInfo = $derived(page.data.userInfo);
const steamId = $derived(page.data.steamId);
</script>

<a href={resolve(`/user/${steamId}`)}>
	{#if userInfo?.avatar}
		<picture>
			<enhanced:img src={userInfo.avatar} alt="User Avatar" height="100%"></enhanced:img>
		</picture>
	{:else}
		<span>{userInfo?.personaname?.split(/(\.|\[|\]|\s)/g)?.pop()?.[0] || '¯\\_(ツ)_/¯'}</span>
	{/if}
</a>

<style lang="postcss">
a {
	display: grid;
	place-items: center;
	width: 2em;
	aspect-ratio: 1 /1;
	justify-self: center;
	position: relative;

	box-shadow: 0 0 1em -0.5em var(--color-primary);
	border: 1px solid color-mix(in srgb, var(--color-primary), transparent 60%);
	cursor: pointer;
	transition: all var(--transition-duration) var(--default-ease);
	padding: 0;
	background: var(--color-primary);

	&:hover,
	&:focus,
	&:focus-within {
		box-shadow: 0 0 1.5em 0em var(--color-primary);
		transform: scale(1.1);
		outline: 2px solid color-mix(in srgb, var(--color-primary), transparent 60%);
	}
	&::before {
		display: none;
	}
	&::after {
		display: none;
	}
	& > span {
		font-size: 1.2em;
		padding: 0 25%;
		font-weight: bold;
		color: var(--text-color-dark);
	}
}
</style>
