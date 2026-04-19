<script lang="ts">
import { beforeNavigate } from '$app/navigation';
import { resolve } from '$app/paths';
import { page } from '$app/state';
import Icon from '@iconify/svelte';
import UserMenu from './UserMenu.svelte';
let steamId = $derived(page.data.steamId);
let userInfo = $derived(page.data.userInfo);
let loginRoute = $state(false);
beforeNavigate(({ from, to }) => {
	if (from?.route.id?.includes('/login') || to?.route.id?.includes('/login')) {
		loginRoute = true;
		return;
	}
	loginRoute = false;
});
</script>

{#if page.url?.pathname === `/user/${steamId}`}
	<a class={loginRoute ? 'login-view-transition' : ''} href={resolve('/logout')}>
		<Icon icon="mdi:logout" width="1.2em" height="1.2em" />
		<span>Logout</span>
	</a>
{:else if !page.route.id?.startsWith('/login')}
	{#if !userInfo}
		<a class={loginRoute ? 'login-view-transition' : ''} href={resolve('/login')}>
			<Icon icon="mdi:steam" width="1.2em" height="1.2em" />
			<span>Sign In</span>
		</a>
	{:else}
		<div class={loginRoute ? 'login-view-transition' : ''}>
			<UserMenu />
		</div>
	{/if}
{/if}

<style lang="postcss">
.login-view-transition {
	view-transition-name: steamlogin;
}

::view-transition-old(steamlogin) {
	animation: var(--transition-duration) var(--ease-vanish) fade-out;
	animation-fill-mode: forwards;
}

::view-transition-new(steamlogin) {
	animation: var(--transition-duration) var(--ease-appear) slide-down-fade;
	animation-fill-mode: forwards;
}

@keyframes slide-down-fade {
	from {
		transform: translateY(-100%);
		opacity: 0;
	}
	to {
		transform: translateY(0);
		opacity: 1;
	}
}

span {
	font-weight: bold;
	overflow: hidden;
	text-overflow: ellipsis;
	width: 100%;
}

a {
	font-size: 1.4em;
	white-space: nowrap;
	display: grid;
	place-items: center;
	grid-template-columns: min-content auto;
	width: 100%;

	gap: 0.25em;
	color: var(--text-color);
	&:hover,
	&:focus,
	&:focus-within {
		color: var(--color-primary);
		outline: 1px solid transparent;
	}
}
</style>
