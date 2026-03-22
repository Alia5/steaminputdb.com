<script lang="ts">
import { resolve } from '$app/paths';
import IconDiscord from '~icons/simple-icons/discord';
import IconGithub from '~icons/simple-icons/github';
</script>

<footer>
	<div class="legal disclaimer">
		<span>SteamInputDB is not affiliated with Steam, Valve Software or any of their partners.</span>
		<span> Steam and the Steam logo are trademarks of Valve Corporation.</span>
		<span> All other trademarks are property of their respective owners.</span>
	</div>
	<div id="links">
		<a href={resolve('/faq')}>FAQ</a>
		<a href="https://discord.gg/hs34MtcHJY" target="_blank" rel="external">
			<IconDiscord style="font-size: 2.4em;" />
		</a>
		<a href="https://github.com/Alia5/steaminputdb.com" target="_blank" rel="external">
			<IconGithub style="font-size: 2.4em;" />
		</a>
	</div>
	<div>
		<span class="c"
			>Copyright © {new Date().getFullYear() > 2026 ? `2026-${new Date().getFullYear()}` : 2026} SteamInputDB
			-
			<a href="https://github.com/Alia5" target="_blank" rel="external">Peter Repukat</a></span>

		<div class="legal links">
			<a href={resolve('/privacy')}>Privacy</a>
			<a href={resolve('/imprint')}>Imprint</a>
		</div>
	</div>
</footer>

<style lang="postcss">
footer {
	padding: 1em;
	padding-bottom: 0.5em;
	background: var(--card-background-noise);
	box-shadow: 0 0px 4px var(--shadow-color);
	display: grid;
	width: 100%;
	grid-template-columns: 1fr 1fr;
	grid-template-rows: auto auto;
	isolation: isolate;

	& > :first-child {
		display: grid;
		grid-template-columns: 1fr;
		grid-row: 1 / span 1;
	}
	& > :last-child {
		grid-column: 1 / span all;
		grid-row: 2 / span 1;
		display: grid;
		grid-template-columns: 1fr auto 1fr;
		place-items: center;
	}
}

#links {
	display: flex;
	flex-flow: row wrap-reverse;
	gap: 2em;
	justify-content: flex-end;
	padding: 0.25em;
	margin-right: 1em;
	align-items: center;

	& a {
		font-weight: bold;
		color: var(--text-color);
		&:has(:global(svg))::after {
			display: none;
		}
		&:hover,
		&:focus-visible {
			color: var(--color-primary);
		}
	}
}

.c {
	color: var(--text-muted);
	opacity: 0.9;
	font-size: 0.9em;
	align-self: end;
	justify-self: center;
	grid-column: 2 / span 1;
	& a {
		color: inherit;
		opacity: 0.9;
	}
}

.disclaimer {
	display: grid;
	grid-template-columns: 1fr;
	color: var(--text-muted);
	align-items: center;
	height: fit-content;
	align-self: center;
	margin-left: 1em;
}

.legal {
	font-size: 0.8em;
	color: var(--text-muted);
	& > a {
		color: inherit;
		height: fit-content;
	}
	&.links {
		display: flex;
		flex-flow: row wrap;
		column-gap: 2em;
		width: 100%;
		justify-content: flex-end;
		align-self: end;

		& > *:not(:last-child) {
			position: relative;
			&::before {
				content: '|';
				position: absolute;
				right: -1em;
				translate: 50% 0;
			}
		}
	}
}
</style>
