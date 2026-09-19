<script lang="ts">
import type { components } from '$lib/api/openapi';
import SIAPI from '$lib/assets/steam_controller_type_svgs/siapi.svg?component';
import IconMDIChecked from '~icons/mdi/check-circle-outline';

import IcoDs4 from '$lib/assets/steam_controller_type_svgs/ps4.svg?component';
import IcoDs5 from '$lib/assets/steam_controller_type_svgs/ps5.svg?component';
import IcoSIAPI from '$lib/assets/steam_controller_type_svgs/siapi.svg?component';
import IcoSC1 from '$lib/assets/steam_controller_type_svgs/steam.svg?component';
import IcoSwitch from '$lib/assets/steam_controller_type_svgs/switchpro.svg?component';
import IcoTriton from '$lib/assets/steam_controller_type_svgs/triton.svg?component';
import IcoXbox from '$lib/assets/steam_controller_type_svgs/xbox.svg?component';
import Icon from '@iconify/svelte';
import IconMdiGamepad from '~icons/mdi/controller';

import { tooltip } from '$lib/attachments/tooltip.svelte';
import IcoMdiDash from '~icons/mdi/minus-circle-outline';

const {
	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	appInfo
}: {
	appInfo: components['schemas']['AppInfoItem'];
} = $props();
</script>

<div id="controller-support">
	<div class="card glass">
		<h3>Controller Support</h3>
		<br />
		<dl>
			<dt>Mixed Input</dt>
			<dd>
				<div>
					<IconMDIChecked style="color: green" />
					<span>Supported with Mod</span>
					<br />
					<br />
					<dl>
						<dt>Locked Glyphs</dt>
						<dd><IconMDIChecked style="color: green" /></dd>
					</dl>
				</div>
			</dd>
			<strong><hr style="opacity: 0.1" /></strong>
			<dt>Glyphs</dt>
			<dd>
				<div>
					<div>
						<div style="display: flex; gap: 1em; width: 100%;">
							<IcoXbox style="width: 2em; height: 2em;" />
							<IcoDs4 style="width: 2em; height: 2em;" />
							<IcoDs5 style="width: 2em; height: 2em;" />
							<IcoSwitch style="width: 2em; height: 2em;" />
							<IcoSC1 style="width: 2em; height: 2em;" />
							<IcoTriton style="width: 2em; height: 2em;" />
						</div>
					</div>
					<div style="display: grid;">
						<div style="display: flex; gap: 0.5em; flex-wrap: wrap;">
							<span class="feature"
								><IcoSIAPI style="width: 1.2em; height: 1.2em;" /> SIAPI Buttons</span
							>
							<span class="feature"
								><IconMdiGamepad style="width: 1.2em; height: 1.2em;" /> In Game Buttons</span
							>
							<span class="feature" style="background-color: red;"
								><IcoMdiDash style="width: 1.2em; height: 1.2em;" /> SIAPI Button Text</span
							>
						</div>
						<span>In game override/selectable <IconMDIChecked style="color: green" /> </span>
					</div>
				</div>
			</dd>
			<strong><hr style="opacity: 0.1" /></strong>
			<dt>Steam Input API</dt>
			<dd>
				<div style="display: grid; gap: 1em;">
					<div style="display: flex; gap: 0.5em; flex-wrap: wrap;">
						<span class="feature" style="background-color: green;"
							><Icon icon="mdi:set-right" width="1.2em" /> In Game Acttions</span
						>
						<span class="feature" style="background-color: yellowgreen;"
							><Icon icon="material-symbols:layers-rounded" width="1.2em" /> Action Sets</span
						>
						<span class="feature" style="background-color: orange;"
							><IcoXbox width="1.2em" /> XInput-Style actions</span
						>
						<span class="feature" style="background-color: red;"
							><Icon icon="mdi:close-circle-outline" width="1.2em" /> DISCONTINUED</span
						>
					</div>
					<dl>
						<dt>Camera Support</dt>
						<dd>
							<div style="display: grid">
								<span><IcoMdiDash style="color: orange" /> Partial</span>
								<dl>
									<dt>Notes</dt>
									<dd>
										Camera action has a built-in deadzone, rendering Gyro/Trackpad Aiming
										impossible.
									</dd>
									<dt>Pixels per 360°</dt>
									<dd><code>11800</code></dd>
								</dl>
							</div>
						</dd>
					</dl>
				</div>
			</dd>
			<dt>Hardware Features</dt>
			<dd>
				<dl>
					<dt>Common</dt>
					<dd>
						<div style="display: flex; gap: 0.5em; flex-wrap: wrap;">
							<span class="feature"
								><Icon icon="game-icons:gyroscope" width="1.2em" />Motion Inputs</span
							>
							<span class="feature"
								><Icon icon="game-icons:gyroscope" width="1.2em" />Native Gyro Camera</span
							>
							<span class="feature"
								><Icon icon="fluent:haptic-strong-16-filled" width="1.2em" /> Rumble</span
							>
							<span class="feature"
								><Icon icon="fluent:haptic-strong-16-filled" width="1.2em" /> Steam Haptics
								<em style="opacity: 0.5">(whatever that means?)</em></span
							>
						</div>
					</dd>
					<dt>Playstation Controllers</dt>
					<dd>
						<div style="display: flex; gap: 0.5em; flex-wrap: wrap;">
							<span class="feature"><Icon icon="mdi:lightbulb-on" width="1.2em" />LightBar</span
							>
							<span class="feature"
								><Icon icon="tabler:activity-heartbeat" width="1.2em" />Adaptive Triggers
								<em>(Wired only)</em></span
							>
							<span class="feature"
								><Icon icon="mdi:volume-vibrate" width="1.2em" />Audio based Haptics</span
							>
							<span
								class="feature"
								{@attach tooltip({
									content: 'Strumming the guitar makes use of the Touchpad',
									outDelay: 2000000,
									arrow: true,

									arrowFollowCursor: true,
									placement: 'bottom-end',
									autoPlacement: false
								})}><Icon icon="fluent:cursor-hover-16-filled" width="1.2em" />Touchpads</span
							>
						</div>
					</dd>
					<dt>Xbox One/Series Controllers</dt>
					<dd>
						<span class="feature"
							><Icon icon="tabler:activity-heartbeat" width="1.2em" />Impulse Triggers</span
						>
					</dd>
				</dl>
			</dd>
		</dl>
	</div>
	<div>
		<div>
			<div class="card glass">
				<span>Controller Support Rating</span>
				<div style="display: grid; place-items: center;">
					<SIAPI style="color: rgb(209 180 93); grid-column: 1; grid-row: 1;" />
					<span
						style="color: red; grid-column: 1; grid-row: 1; font-weight: bold; font-size: 2em; filter: drop-shadow(0 0 3px black) drop-shadow(0 0 3px black); transform: rotate(17deg);"
						>Placeholder</span
					>
				</div>
			</div>
			<div class="card glass" style="height: 100%; display: grid; place-items: center;">
				<h3>Separate SteamInput API rating</h3>
				<span>If applicable?</span>
				<h3>Or Mixed Input Rating?</h3>
			</div>
		</div>
		<div class="card glass">
			<h3>Additional Info</h3>
			<p>
				<strong style="color: red">PLACEHOLDER</strong>
				<br />

				DualSense's Haptic Feedback and Adaptive Triggers only works via USB.
				<br />Adaptive Triggers is not supported under Bluetooth.
				<br />Steam Haptic/Traditional Rumbles is not supported when using Dualsense controller.
				<br />
				DualSense's Haptic Feedback and Adaptive Triggers only works via USB.
				<br />
				Adaptive Triggers is not supported under Bluetooth. Steam Haptic/Traditional Rumbles is not supported
				when using Dualsense controller.
				<br />
				<br />
				<strong
					style="display: grid; grid-template-columns: auto auto; gap: 1em; place-items: center;"
					><span>Mixed Input Mod:</span>
					<span><a href="#">NexusMods</a><br /><a href="#">GitHub</a></span>
				</strong>
			</p>
		</div>
	</div>
</div>

<style lang="postcss">
#controller-support {
	display: flex;
	flex-flow: row wrap-reverse;
	overflow: clip;
	overflow-clip-margin: 2em;

	--gap: 1em;
	--info-min-width: 58ch;
	width: 100%;

	gap: 1em;

	& > :first-child {
		padding: 1em;
		overflow-x: hidden;
		position: relative;
		backdrop-filter: blur(8px);

		flex-basis: calc(66% - var(--gap));
		min-width: min(100%, var(--info-min-width));
	}
	& > :last-child {
		display: grid;
		height: fit-content;
		align-self: flex-end;
		width: 100%;
		flex-shrink: 1;
		gap: 1em;

		--width: 20ch;

		flex-basis: calc(34%);
		min-width: min(100%, var(--width));
		& > :first-child {
			display: grid;
			place-items: center;
			margin: auto;
			gap: 1em;
			width: 100%;
			grid-template-columns: repeat(auto-fit, minmax(calc(var(--width) - 1em), auto));
			:global(> *) {
				width: 100%;
				padding: 1em;
			}
		}
	}
}

dl {
	display: grid;
	place-items: center;
	grid-template-columns: auto 1fr;
	grid-column-gap: 1em;
	grid-row-gap: 0.5em;
	width: 100%;

	& > strong {
		grid-column: 1 / -1;
		justify-self: start;
		font-size: 1.2em;
		margin-bottom: 0.25em;
		width: 100%;
		&:not(:nth-of-type(1)) {
			margin-top: 1em;
			position: relative;
			&::before {
				content: '';
				position: absolute;
				inset: 0;
				border-top: 1px solid var(--text-color);
				top: -0.75em;
				opacity: 0.33;
			}
		}
	}

	&:has(> strong) {
		dt {
			margin-left: 1em;
		}
	}

	& dt {
		font-weight: bold;
		justify-self: start;
		display: flex;
		align-items: center;
		gap: 1ch;
	}
	dd {
		max-width: 100%;
		overflow: clip;
		overflow-clip-margin: 2em;
		display: flex;
		align-items: center;
		gap: 1ch;
		& a {
			font-weight: bold;
		}
	}

	& > :nth-child(2n) {
		justify-self: start;
		& em {
			opacity: 0.75;
			& a {
				font-weight: bold;
			}
		}
	}
}

.feature {
	display: flex;
	align-items: center;
	gap: 0.2ch;
	padding: 0.4em 0.6em;
	border-radius: var(--border-radius);
	box-shadow: var(--card-shadow);
	background: var(--card-background-noise);
	font-size: 0.9em;
	font-weight: bold;
	white-space: nowrap;
}
</style>
