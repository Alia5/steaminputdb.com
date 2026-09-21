<script lang="ts">
import { resolve } from '$app/paths';
import type { components } from '$lib/api/openapi';
import {
	AppGlyphTagType,
	HWFeature,
	HWFeatureControllerType,
	MixedInputSupportType,
	SteamInputAPISupportType,
	SteamInputCameraSupport
} from '$lib/api/steaminputdbEnums';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { CONTROLLER_LIST } from '$lib/components/search/controllerlist.svelte';

import IcoHaptics from '$lib/assets/icohaptics.svg?component';
import IcoMixedInput from '$lib/assets/mixedinput.svg?component';
import IcoDs5 from '$lib/assets/steam_controller_type_svgs/ps5.svg?component';
import IcoSIAPI from '$lib/assets/steam_controller_type_svgs/siapi.svg?component';
import IcoSwitch from '$lib/assets/steam_controller_type_svgs/switchpro.svg?component';
import IcoXbox from '$lib/assets/steam_controller_type_svgs/xbox.svg?component';
import IcoDotsPer360 from '$lib/components/layout_preview/glyphs/cmnd_dots_per_360_calibration_spin.svg?component';
import IcoTrigger from '$lib/components/layout_preview/glyphs/shared_trigger.svg?component';
import MarkdownSSR from '$lib/components/markdown/markdownSSR.svelte';
import IcoTouchpads from '~icons/fluent/cursor-hover-16-filled';
import IcoRumble from '~icons/fluent/haptic-strong-16-filled';
import IcoGyro from '~icons/game-icons/gyroscope';
import IcoActionSets from '~icons/material-symbols/layers-rounded';
import IconMDIChecked from '~icons/mdi/check-circle-outline';
import IcoMdiCross from '~icons/mdi/close-circle-outline';
import IconMdiGamepad from '~icons/mdi/controller';
import IcoDpad from '~icons/mdi/gamepad';
import IconMdiGamepadCircle from '~icons/mdi/gamepad-circle';
import IcoLightbar from '~icons/mdi/lightbulb-on';
import IcoMdiDash from '~icons/mdi/minus-circle-outline';
import IcoMdiSetRight from '~icons/mdi/set-right';
import IcoSteam from '~icons/mdi/steam';
import IcoAudioHaptics from '~icons/mdi/volume-vibrate';
import getModUrlName from './modsUrlNames';

const {
	appInfo
}: {
	appInfo: components['schemas']['AppInfoItem'];
} = $props();
let steaminputdbInfo = $derived(appInfo?.steaminputdb_info);
</script>

{#snippet controllerIcon(type: string | undefined, glyphNotes?: string)}
	{#each CONTROLLER_LIST.filter((c) => {
		return c.type === type;
	}) as controller (controller.type)}
		<controller.icon
			style="width: 2.4em; height: 2.4em;"
			{@attach tooltip({
				content: glyphNotes ?? '',
				outDelay: 200,
				arrow: true,
				autoPlacement: true,
				arrowFollowCursor: false
			})}
		/>
	{/each}
{/snippet}

{#snippet steamInputAPIGlyphTag(tag: number)}
	{#if tag === AppGlyphTagType.SteamInputAPIButtons}
		<span class="feature" style="background-color: #00dfff;">
			<IcoSIAPI style="width: 1.4em; height: 1.4em;" /> Steam Input API Glyphs
		</span>
	{:else if tag === AppGlyphTagType.InGameButtons}
		<span class="feature" style="background-color: #12ff12;">
			<IconMdiGamepad style="width: 1.2em; height: 1.2em;" /> In-Game Glyphs
		</span>
	{:else if tag === AppGlyphTagType.SteamInputAPITexts}
		<span class="feature" style="background-color: red;">
			<IcoMdiDash style="width: 1.2em; height: 1.2em;" /> Steam Input API Button Text
		</span>
	{/if}
{/snippet}

{#snippet hwFeature(feature: number, family: number | string | undefined, notes: string | undefined)}
	<span
		class="feature"
		{@attach tooltip({
			content: notes ?? '',
			outDelay: 200,
			arrow: true,
			autoPlacement: true,
			arrowFollowCursor: false
		})}
	>
		{#if feature === HWFeature.MotionInputs}
			<IcoGyro style="width: 1.4em; height: 1.4em;" /> Motion Inputs
		{:else if feature === HWFeature.NativeGyroCamera}
			<IcoDotsPer360 style="width: 1.2em; height: 1.2em;" /> Native gyro camera
		{:else if feature === HWFeature.Rumble}
			<IcoRumble style="width: 1.2em; height: 1.2em;" /> Rumble
		{:else if feature === HWFeature.HDHaptics}
			<IcoHaptics style="width: 1.2em; height: 1.2em;" />
			{#if family == HWFeatureControllerType.Steam}
				Steam Haptics
			{:else}
				HD Haptics
			{/if}
		{:else if feature === HWFeature.Lightbar}
			<IcoLightbar style="width: 1.2em; height: 1.2em;" /> Lightbar
		{:else if feature === HWFeature.AdaptiveTriggersWired}
			<IcoTrigger style="width: 1.2em; height: 1.2em;" /> Adaptive Triggers (Wired)
		{:else if feature === HWFeature.ImpulseTriggers}
			<IcoTrigger style="width: 1.2em; height: 1.2em;" /> Impulse Triggers
		{:else if feature === HWFeature.Touchpads}
			<IcoTouchpads style="width: 1.2em; height: 1.2em;" /> Touchpads
		{:else if feature === HWFeature.AudioHaptics}
			<IcoAudioHaptics style="width: 1.2em; height: 1.2em;" /> Audio-based Haptics
		{/if}
	</span>
{/snippet}

{#snippet steamInputAPISupportTags(tag: number)}
	{#if tag === SteamInputAPISupportType.InGameActions}
		<span class="feature" style="background-color: turquoise;">
			<IcoMdiSetRight style="width: 1.2em; height: 1.2em;" /> In-Game Actions
		</span>
	{:else if tag === SteamInputAPISupportType.ActionSets}
		<span class="feature" style="background-color: yellowgreen">
			<IcoActionSets style="width: 1.2em; height: 1.2em;" /> Action Sets
		</span>
	{:else if tag === SteamInputAPISupportType.XInputActions}
		<span class="feature" style="background-color: lightgreen;">
			<IcoXbox style="width: 1.2em; height: 1.2em;" /> XInput Style Actions
		</span>
	{:else if tag === SteamInputAPISupportType.Discontinued}
		<span class="feature" style="background-color: red;">
			<IcoMdiCross style="width: 1.2em; height: 1.2em;" /> Discontinued
		</span>
	{/if}
{/snippet}

<section id="controller-support">
	{@render controllerSupportContent()}
</section>

{#snippet controllerSupportContent()}
	<div
		class={steaminputdbInfo?.glyphs ||
		steaminputdbInfo?.mixed_input ||
		steaminputdbInfo?.hw_features ||
		steaminputdbInfo?.steaminputapi_support ||
		steaminputdbInfo?.controller_support_notes
			? 'card glass'
			: ''}
	>
		{#if Object.entries(appInfo?.official_configs ?? {}).length}
			<section class="official-configs">
				<h3>Official Configs</h3>
				{#each Object.entries(appInfo?.official_configs ?? {}) as [controller_type, config_id] (config_id)}
					{@const controller_list_entry = CONTROLLER_LIST.find(
						(controller) => controller.type === controller_type
					)}
					<a href={resolve(`/config/${config_id}`)} class="button">
						{#if controller_list_entry}
							<controller_list_entry.icon width="2em" height="2em" />
						{:else}
							<IcoDpad style="width: 2em; height: 2em;" />
						{/if}
						<span>{controller_list_entry?.niceName ?? 'Generic'}</span>
					</a>
				{/each}
			</section>
		{/if}
		{#if steaminputdbInfo?.glyphs}
			<div class="rule-divider"></div>
			<section id="mixed-input-and-glyphs">
				<section id="mixed-input" class="info-group">
					<h3><IcoMixedInput style="width: 2em; height: 2em;" />Mixed Input</h3>
					{#if steaminputdbInfo?.mixed_input}
						{@const mixedinputInfo = steaminputdbInfo?.mixed_input}
						{@const mixedInputType = mixedinputInfo?.type}
						<div>
							<div>
								{#if mixedInputType === MixedInputSupportType.Unsupported}
									<IcoMdiCross style="width: 1.6em; height: 1.6em; color: firebrick" />
									<span>Unsupported</span>
								{/if}
								{#if mixedInputType === MixedInputSupportType.Supported}
									<IconMDIChecked style="width: 1.6em; height: 1.6em; color: green" />
									<span>Supported</span>
								{/if}
								{#if mixedInputType === MixedInputSupportType.WithMod}
									<IconMDIChecked style="width: 1.6em; height: 1.6em; color: yellowgreen" />
									<span>Requires Mod</span>
								{/if}
								{#if mixedInputType === MixedInputSupportType.Partial}
									<IcoMdiDash style="width: 1.6em; height: 1.6em; color: orange" />
									<span>Partial</span>
								{/if}
								{#if mixedInputType === MixedInputSupportType.Unknown}
									<IcoMdiDash style="width: 1.6em; height: 1.6em; color: gray" />
									<span>Unknown</span>
								{/if}
							</div>
							<div>
								{#if mixedinputInfo?.glyph_flicker}
									<IcoMdiCross style="width: 1.6em; height: 1.6em; color: firebrick" />
									<span>Glyph Flicker</span>
								{:else}
									<IconMDIChecked style="width: 1.6em; height: 1.6em; color: green" />
									<span>Without Glyph Flicker</span>
								{/if}
							</div>
							{#if mixedinputInfo?.notes}
								<div class="notes mdcontainer">
									<MarkdownSSR content={mixedinputInfo?.notes} />
								</div>
							{/if}
						</div>
					{:else}
						<div>
							<IcoMdiDash style="width: 1.6em; height: 1.6em; color: gray" />
							<span>Unknown</span>
						</div>
					{/if}
				</section>
				<section id="glyphs" class="info-group">
					<h3><IconMdiGamepadCircle style="width: 1.6em; height: 1.6em;" /> Glyphs</h3>
					{#if steaminputdbInfo?.glyphs}
						{@const glyphsInfo = steaminputdbInfo?.glyphs}
						<div>
							<div class="ctrl-glyphs">
								{#each glyphsInfo.controllers as controller (controller.controller_type)}
									{@render controllerIcon(controller.controller_type, controller.notes)}
								{/each}
							</div>

							<div>
								{#if glyphsInfo.autodetect}
									<IconMDIChecked style="width: 1.6em; height: 1.6em; color: green" />
									<span>Autodetect</span>
								{/if}
								{#if glyphsInfo.manual_select}
									<IconMDIChecked style="width: 1.6em; height: 1.6em; color: green" />
									<span>Manual Override / Lock</span>
								{/if}
							</div>
							{#if steaminputdbInfo?.steaminputapi_support?.glyphs}
								{@const siapiGlphTags = steaminputdbInfo?.steaminputapi_support?.glyphs}
								<div>
									{#each siapiGlphTags as tag (tag)}
										{@render steamInputAPIGlyphTag(tag)}
									{/each}
								</div>
							{/if}
							<div>
								{#if glyphsInfo?.notes}
									<div class="notes mdcontainer">
										<MarkdownSSR content={glyphsInfo?.notes} />
									</div>
								{/if}
							</div>
						</div>
					{:else}
						<div>
							<IcoMdiDash style="width: 1.6em; height: 1.6em; color: gray" />
							<span>Unknown</span>
						</div>
					{/if}
				</section>
			</section>
		{/if}
		{#if steaminputdbInfo?.steaminputapi_support || (steaminputdbInfo?.hw_features || []).length}
			<div class="horizontal-group">
				{#if (steaminputdbInfo?.hw_features || []).length}
					{@const hwFeatures = steaminputdbInfo?.hw_features}
					<section id="hw-features" class="info-group">
						<h3><IconMdiGamepad style="width: 1.6em; height: 1.6em;" /> Hardware Features</h3>
						<div>
							{#each Object.entries((hwFeatures || []).reduce((acc, feature) => {
										if (!acc[feature.controller_family]) {
											acc[feature.controller_family] = [];
										}
										acc[feature.controller_family]!.push(feature);
										return acc;
									}, {} as Record<number, typeof hwFeatures>)) as [ctrl_family, features] (ctrl_family)}
								{#if ctrl_family == `${HWFeatureControllerType.Common}`}
									<div>
										<IcoDpad style="width: 1.4em;" />
										<span>Generic</span>
										{#each features as f (f.feature)}
											{@render hwFeature(f.feature, ctrl_family, f.notes)}
										{/each}
									</div>
								{:else if ctrl_family == `${HWFeatureControllerType.Steam}`}
									<div>
										<IcoSteam style="width: 1.4em;" />
										<span>Steam</span>
										{#each features as f (f.feature)}
											{@render hwFeature(f.feature, ctrl_family, f.notes)}
										{/each}
									</div>
								{:else if ctrl_family == `${HWFeatureControllerType.PlayStation}`}
									<div>
										<IcoDs5 style="width: 1.2em" />
										<span>PlayStation</span>
										{#each features as f (f.feature)}
											{@render hwFeature(f.feature, ctrl_family, f.notes)}
										{/each}
									</div>
								{:else if ctrl_family == `${HWFeatureControllerType.Xbox}`}
									<div>
										<IcoXbox style="width: 1.2em;" />
										<span>Xbox</span>
										{#each features as f (f.feature)}
											{@render hwFeature(f.feature, ctrl_family, f.notes)}
										{/each}
									</div>
								{:else if ctrl_family == `${HWFeatureControllerType.Nintendo}`}
									<div>
										<IcoSwitch style="width:1.2em;" />
										<span>Nintendo</span>
										{#each features as f (f.feature)}
											{@render hwFeature(f.feature, ctrl_family, f.notes)}
										{/each}
									</div>
								{/if}
							{/each}
						</div>
					</section>
				{/if}
				{#if steaminputdbInfo?.steaminputapi_support}
					{@const steaminputapi = steaminputdbInfo?.steaminputapi_support}
					<section id="steaminputapi" class="info-group">
						<h3><IcoSIAPI style="width: 2em; height: 2em;" /> Steam Input API Support</h3>
						{#if (steaminputapi?.support_tags || []).length}
							{@const siapiSuppTags = steaminputdbInfo?.steaminputapi_support?.support_tags}
							<div>
								<div>
									{#each siapiSuppTags as tag (tag)}
										{@render steamInputAPISupportTags(tag)}
									{/each}
								</div>
							</div>
						{/if}
						<div>
							{#if steaminputapi?.camera_support ?? 0 > 0}
								<div>
									{#if steaminputapi?.camera_support === SteamInputCameraSupport.None}
										<IcoMdiCross style="width: 1.6em; height: 1.6em; color: red" />
										<strong>No</strong>
									{:else if steaminputapi?.camera_support === SteamInputCameraSupport.Partial}
										<IcoMdiDash style="width: 1.6em; height: 1.6em; color: orange" />
										<strong>Partial</strong>
									{:else if steaminputapi?.camera_support === SteamInputCameraSupport.Full}
										<IconMDIChecked style="width: 1.6em; height: 1.6em; color: green" />
										<strong>Full</strong>
									{/if}
									<span>SIAPI Camera Support</span>
								</div>
								{#if steaminputapi?.pixels_per_360}
									<div>
										<IcoDotsPer360 style="width: 1.6em; height: 1.6em;" />
										<span>Pixels Per 360°</span>
										<code>{steaminputapi.pixels_per_360}</code>
									</div>
								{/if}
							{/if}
						</div>

						{#if steaminputapi?.notes}
							<div>
								<div class="notes mdcontainer">
									<MarkdownSSR content={steaminputapi?.notes} />
								</div>
							</div>
						{/if}
					</section>
				{/if}
			</div>
		{/if}
	</div>
	{#if steaminputdbInfo?.controller_support_notes || (steaminputdbInfo?.mixed_input?.mixed_input_mod_urls || []).length}
		<aside id="controller-support-notes" class="card glass">
			<h3>Additional Info</h3>
			{#if steaminputdbInfo?.controller_support_notes}
				<div class="mdcontainer scrollable">
					<div class="content">
						<MarkdownSSR content={steaminputdbInfo?.controller_support_notes} />
						<p></p>
					</div>
				</div>
			{/if}
			{#if (steaminputdbInfo?.mixed_input?.mixed_input_mod_urls || []).length}
				<div class="mod-link-list">
					<strong style="padding-top: 1.5em;"
						>Mixed Input / Controller support Mod{#if (steaminputdbInfo?.mixed_input?.mixed_input_mod_urls || []).length > 1}s{/if}
					</strong>
					{#each steaminputdbInfo?.mixed_input?.mixed_input_mod_urls as url, idx (idx)}
						<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
						<a href={url} target="_blank" rel="noopener noreferrer">{getModUrlName(url)}</a>
					{/each}
				</div>
			{/if}
		</aside>
	{/if}
{/snippet}

<style lang="postcss">
#controller-support {
	display: flex;
	flex-flow: row wrap;
	overflow: clip;
	overflow-clip-margin: 2em;

	--gap: 1em;
	--info-min-width: 58ch;
	width: 100%;
	align-items: stretch;
	height: 100%;

	gap: 1em;
	& > :first-child {
		flex: 3 1 var(--info-min-width);
		min-width: 0;
		display: grid;
		align-content: start;
		gap: 1em;
	}
	& > aside {
		flex: 1 1 24ch;
		min-width: 0;
		display: flex;
		flex-flow: column nowrap;
		align-items: stretch;
		gap: 1em;
		& .mod-link-list {
			display: grid;
			gap: 0.5em;
		}

		& .scrollable {
			flex: 1 1 auto;
			min-height: 0;
			overflow: auto;
			overflow: auto;
			& > .content {
				max-height: 24em;
			}
		}
	}
}

h3 {
	font-size: 1.1em;
	align-items: center;
	display: flex;
	gap: 0.5ch;
}

.rule-divider {
	opacity: 0.1;
	height: 1px;
	background: var(--text-color);
}

#mixed-input-and-glyphs {
	display: flex;
	flex-flow: row wrap;
	gap: 1em;
	justify-content: space-evenly;
}

.horizontal-group {
	display: flex;
	flex-flow: row wrap;
	gap: 1em;
	justify-content: space-evenly;
	align-items: center;
	& > * {
		max-width: 50%;
		width: fit-content;
	}
}

.info-group {
	display: flex;
	flex-flow: column wrap;
	gap: 1em;
	align-items: center;
	& > div {
		display: grid;
		align-items: center;
		justify-items: center;
		gap: 0.5em;
		& > div {
			display: flex;
			flex-flow: row wrap;
			gap: 0.5ch;
			align-items: center;
			justify-content: center;
		}
	}

	.ctrl-glyphs {
		display: flex;
		flex-flow: row wrap;
		gap: 1em;
		align-items: center;
	}
}

#hw-features {
	& > div {
		gap: 1em;
		& > div {
			position: relative;
			display: flex;
			flex-flow: row wrap;
			align-items: center;
			gap: 0.5em;
			width: 100%;
			justify-content: center;
		}
		& > div:not(:first-child)::after {
			content: '';
			position: absolute;
			height: 1px;
			top: -0.41em;
			left: 0;
			width: 100%;
			background: var(--text-color);
			opacity: 0.2;
		}
	}
}

#mixed-input {
	display: flex;
	flex-flow: column wrap;
	gap: 1em;
	align-items: center;
	& > div {
		display: grid;
		align-items: center;
		justify-items: center;
		gap: 0.5em;
		& > div {
			display: flex;
			flex-flow: row wrap;
			gap: 0.5ch;
			align-items: center;
		}
	}
}

.notes {
	max-width: 42ch;
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
	& > :global(:first-child) {
		margin-right: 0.5ch;
	}
}

.official-configs {
	display: flex;
	flex-flow: row wrap;
	gap: 1ch;
	margin-left: auto;
	justify-content: center;
	width: 100%;
	align-items: baseline;
	height: fit-content;
	& > :first-child {
		width: 100%;
		padding-bottom: 0.5ch;
		justify-content: center;
		/* filter: drop-shadow(2px 2px 2px rgba(0, 0, 0, 0.452)) drop-shadow(0px 0px 8px rgba(0, 0, 0, 0.63)); */
	}
	& > a {
		display: grid;
		place-items: center;
		font-weight: bold;
		display: grid;
		align-items: center;
		justify-content: center;
		height: fit-content;
		padding: 0.5em 1em;
		gap: 0.5ch;
		background: linear-gradient(
			215deg,
			color-mix(in srgb, var(--card-color), transparent 35%) 0%,
			color-mix(in srgb, var(--card-color), transparent 60%) 70%
		);
		& > :global(svg) {
			width: 2em;
			height: 2em;
		}
	}
}
</style>
