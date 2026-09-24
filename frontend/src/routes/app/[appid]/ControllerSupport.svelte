<script lang="ts">
import { resolve } from '$app/paths';
import type { components } from '$lib/api/openapi';
import {
	AppGlyphTagType,
	HWFeature,
	HWFeatureControllerType,
	isHWFeatureForFamily,
	MixedInputSupportType,
	SteamInputAPISupportType,
	SteamInputCameraSupport,
	SteamInputType
} from '$lib/api/steaminputdbEnums';
import { tooltip } from '$lib/attachments/tooltip.svelte';
import { CONTROLLER_LIST } from '$lib/components/search/controllerlist.svelte';

import IcoHaptics from '$lib/assets/icohaptics.svg?component';
import IcoMixedInput from '$lib/assets/mixedinput.svg?component';
import IcoCtrlGeneric from '$lib/assets/steam_controller_type_svgs/controller_generic.svg?component';
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
import IconMdiGamepadCircle from '~icons/mdi/gamepad-circle';
import IcoHelp from '~icons/mdi/help-circle-outline';
import IcoLightbar from '~icons/mdi/lightbulb-on';
import IcoMdiDash from '~icons/mdi/minus-circle-outline';
import IcoMdiSetRight from '~icons/mdi/set-right';
import IcoSteam from '~icons/mdi/steam';
import IcoSpeaker from '~icons/mdi/volume';
import IcoAudioHaptics from '~icons/mdi/volume-vibrate';
import { getModHostNameNice, getModUrlName } from './modsUrlNames';

const {
	appInfo
}: {
	appInfo: components['schemas']['AppInfoItem'];
} = $props();
let steaminputdbInfo = $derived(appInfo?.steaminputdb_info);

const showsCard = $derived(
	(steaminputdbInfo?.mixed_input &&
		Object.entries(steaminputdbInfo?.mixed_input ?? {}).filter(([k, v]) => {
			if (k == 'mixed_input_mods') {
				console.log('skipped mods');
				return false;
			}
			return !!v;
		})?.length) ||
		steaminputdbInfo?.glyphs ||
		steaminputdbInfo?.hw_features ||
		Object.values(steaminputdbInfo?.steaminputapi_support ?? {})?.filter(Boolean)?.length ||
		steaminputdbInfo?.hw_feature_notes ||
		steaminputdbInfo?.controller_support_notes
);

const controllerListOrder = (type?: string) => {
	const idx = CONTROLLER_LIST.findIndex((c) => c.type === type);
	return idx < 0 ? CONTROLLER_LIST.length : idx;
};
</script>

{#snippet controllerIcon(type: string | undefined, glyphNotes?: string)}
	{#each CONTROLLER_LIST.filter((c) => {
		return c.type === type;
	}) as controller (controller.type)}
		<div
			style="display: grid; justify-items: center;"
			{@attach tooltip({
				content: glyphNotes ?? '',
				outDelay: 100,
				arrow: true,
				autoPlacement: false,
				placement: 'top',
				arrowFollowCursor: false
			})}
		>
			{#if glyphNotes}
				<IcoHelp style="width: 1.3em; height: 1.3em; opacity: 0.7;" />
			{:else}
				<div class="spacer" style="width: 1.3em; height: 1.3em;"></div>
			{/if}
			{#if type?.includes('xbox')}
				<IcoXbox style="width: 2.4em; height: 2.4em;" />
			{:else}
				<controller.icon style="width: 2.4em; height: 2.4em;" />
			{/if}
		</div>
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
			outDelay: 100,
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
		{:else if feature === HWFeature.AdaptiveTriggersWireless}
			<IcoTrigger style="width: 1.2em; height: 1.2em;" /> Adaptive Triggers (Wireless)
		{:else if feature === HWFeature.ImpulseTriggers}
			<IcoTrigger style="width: 1.2em; height: 1.2em;" /> Impulse Triggers
		{:else if feature === HWFeature.Touchpads}
			<IcoTouchpads style="width: 1.2em; height: 1.2em;" /> Touchpads
		{:else if feature === HWFeature.AudioHaptics}
			<IcoAudioHaptics style="width: 1.2em; height: 1.2em;" /> Audio-based Haptics
		{:else if feature === HWFeature.Speaker}
			<IcoSpeaker style="width: 1.2em; height: 1.2em;" /> Speaker / Audio Output
		{/if}
		{#if notes}
			<IcoHelp style="width: 1.2em; height: 1.2em; opacity: 0.7; margin-left: 0.5em;" />
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
	{#if showsCard || (steaminputdbInfo?.mixed_input?.mixed_input_mods || []).length}
		<div
			class="ctrl-support-help"
			{@attach tooltip({
				content: `This information is community sourced and may be outdated or incomplete
                
A feature where any logged in User will be able to to propose changes will follow soon™`,
				outDelay: 100,
				arrow: true,
				autoPlacement: true,
				arrowFollowCursor: false
			})}
		>
			<IcoHelp style="width: 1.6em; height: 1.6em; opacity: 0.7;" />
		</div>
	{/if}
</section>

{#snippet controllerSupportContent()}
	<div class={showsCard ? 'card glass' : ''}>
		{#if Object.entries(appInfo?.official_configs ?? {}).length}
			<section class={showsCard ? 'official-configs' : 'official-configs no-divider'}>
				<h3>Official Configs</h3>
				{#each Object.entries(appInfo?.official_configs ?? {}) as [controller_type, config_id] (config_id)}
					{@const controller_list_entry = CONTROLLER_LIST.find(
						(controller) => controller.type === controller_type
					)}
					<a href={resolve(`/config/${config_id}`)} class="button">
						{#if controller_list_entry}
							<controller_list_entry.icon width="2em" height="2em" />
						{:else}
							<IcoCtrlGeneric style="width: 2em; height: 2em;" />
						{/if}
						<span>{controller_list_entry?.niceName ?? 'Generic'}</span>
					</a>
				{/each}
			</section>
		{/if}
		{#if steaminputdbInfo?.glyphs}
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
							<!-- eslint-disable-next-line prettier/prettier -->
                                    {#if mixedInputType !== MixedInputSupportType.Unsupported
                                        && mixedInputType !== MixedInputSupportType.Unknown 
                                     }
								{#if mixedinputInfo?.glyph_flicker === false}
									<IconMDIChecked style="width: 1.6em; height: 1.6em; color: green" />
									<span>Without Glyph Flicker</span>
								{:else if mixedinputInfo?.glyph_flicker === true}
									<IcoMdiCross style="width: 1.6em; height: 1.6em; color: firebrick" />
									<span>Glyph Flicker</span>
								{/if}
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
			<section id="glyphs" class="info-group">
				<h3><IconMdiGamepadCircle style="width: 1.4em; height: 1.4em;" /> Glyphs</h3>
				{#if steaminputdbInfo?.glyphs}
					{@const glyphsInfo = steaminputdbInfo?.glyphs}
					<div>
						<div class="ctrl-glyphs">
							{#each glyphsInfo.controllers?.sort((a, b) => {
								// eslint-disable-next-line prettier/prettier
                                return controllerListOrder(a.controller_type)
                                    - controllerListOrder(b.controller_type)
							}) as controller (controller.controller_type)}
								{@render controllerIcon(controller.controller_type, controller.notes)}
							{/each}
						</div>

						<div>
							{#if glyphsInfo.autodetect == true}
								<IconMDIChecked style="width: 1.6em; height: 1.6em; color: green" />
								<span>Autodetect</span>
							{/if}
							{#if glyphsInfo.manual_select == true}
								<IconMDIChecked style="width: 1.6em; height: 1.6em; color: green" />
								<span>Manual Override / Lock</span>
							{/if}
						</div>
						{#if steaminputdbInfo?.steaminputapi_support?.glyphs && (steaminputdbInfo.steaminputapi_support?.steam_input_type ?? 0) > 1}
							{@const siapiGlphTags = steaminputdbInfo?.steaminputapi_support?.glyphs}
							<div>
								{#each siapiGlphTags as tag (tag)}
									{@render steamInputAPIGlyphTag(tag)}
								{/each}
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
		{/if}
		{#if steaminputdbInfo?.steaminputapi_support || (steaminputdbInfo?.hw_features || []).length}
			<!-- eslint-disable prettier-prettier -->
			{#if (steaminputdbInfo?.steaminputapi_support?.support_tags || []).length || steaminputdbInfo?.steaminputapi_support?.steam_input_type || steaminputdbInfo?.steaminputapi_support?.camera_support || steaminputdbInfo?.steaminputapi_support?.notes}
				<!-- eslint-enable prettier-prettier -->
				{@const steaminputapi = steaminputdbInfo?.steaminputapi_support}
				<section id="steaminputapi" class="info-group">
					<h3><IcoSIAPI style="width: 2em; height: 2em;" /> Steam Input Support</h3>
					{#if !!steaminputapi?.steam_input_type}
						<div>
							<div>
								{#if steaminputapi?.steam_input_type === SteamInputType.None}
									<IcoMdiCross style="width: 1.6em; height: 1.6em;" />
									<span>Steam Input unaware</span>
								{:else if steaminputapi?.steam_input_type === SteamInputType.VirtualGamepad}
									<IcoMdiDash style="width: 1.6em; height: 1.6em;" />
									<span>Steam Virtual Gamepad</span>
								{:else if steaminputapi?.steam_input_type === SteamInputType.NativeAPI}
									<IconMDIChecked style="width: 1.6em; height: 1.6em; color: green" />
									<span>Native Steam Input API</span>
								{:else if steaminputapi?.steam_input_type === SteamInputType.Required}
									<IconMDIChecked style="width: 1.6em; height: 1.6em; color: orange" />
									<span>Steam Input Required</span>
								{/if}
							</div>
						</div>
					{/if}
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
							{#if (steaminputapi?.pixels_per_360 ?? '') !== '' && (steaminputapi?.camera_support ?? 0 > 0)}
								<div>
									<IcoDotsPer360 style="width: 1.6em; height: 1.6em;" />
									<span>Pixels Per 360°</span>
									{typeof steaminputapi?.pixels_per_360}
									<code>{steaminputapi?.pixels_per_360 ?? ''}</code>
								</div>
							{/if}
						{/if}
					</div>
				</section>
			{/if}
			{#if (steaminputdbInfo?.hw_features || []).length}
				{@const hwFeatures = (steaminputdbInfo?.hw_features || []).filter((f) =>
					isHWFeatureForFamily(f.controller_family, f.feature)
				)}
				{#if hwFeatures.length}
					<section id="hw-features" class="info-group">
						<h3>
							<IconMdiGamepad style="width: 1.6em; height: 1.6em;" /> Hardware Features
						</h3>
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
										<IcoCtrlGeneric style="width: 1.4em;" />
										<span>Common</span>
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
			{/if}
		{/if}
	</div>
	<!-- eslint-disable-next-line -->
	{#if
        steaminputdbInfo?.controller_support_notes
        || steaminputdbInfo?.mixed_input?.notes
        || steaminputdbInfo?.glyphs?.notes
        || steaminputdbInfo?.hw_feature_notes
        || steaminputdbInfo?.steaminputapi_support?.notes
        || (steaminputdbInfo?.mixed_input?.mixed_input_mods || []).length
        }
		{@const generalNotes = steaminputdbInfo?.controller_support_notes}
		{@const mixedInputNotes = steaminputdbInfo?.mixed_input?.notes}
		{@const glyphNotes = steaminputdbInfo?.glyphs?.notes}
		{@const hwFeatureNotes = steaminputdbInfo?.hw_feature_notes}
		{@const steamInputNotes = steaminputdbInfo?.steaminputapi_support?.notes}
		{@const mixedInputMods = steaminputdbInfo?.mixed_input?.mixed_input_mods ?? []}
		<aside id="controller-support-notes" class="card glass">
			<div class="notes-container scrollable">
				<div>
					{#if mixedInputNotes}
						<h3><IcoMixedInput style="width: 1.4em; height: 1.4em;" /> Mixed Input</h3>
						<div class="notes mdcontainer">
							<MarkdownSSR content={mixedInputNotes} />
						</div>
					{/if}
					{#if glyphNotes}
						<h3><IconMdiGamepadCircle style="width: 1.2em; height: 1.2em;" /> Glyphs</h3>
						<div class="notes mdcontainer">
							<MarkdownSSR content={glyphNotes} />
						</div>
					{/if}
					{#if steamInputNotes}
						<h3><IcoSIAPI style="width: 1.4em; height: 1.4em;" /> Steam Input</h3>
						<div class="notes mdcontainer">
							<MarkdownSSR content={steamInputNotes} />
						</div>
					{/if}
					{#if hwFeatureNotes}
						<h3><IconMdiGamepad style="width: 1.2em; height: 1.2em;" /> Hardware Features</h3>
						<div class="notes mdcontainer">
							<MarkdownSSR content={hwFeatureNotes} />
						</div>
					{/if}
					{#if generalNotes}
						<h3>Additional Information</h3>
						<div class="notes mdcontainer">
							<MarkdownSSR content={generalNotes} />
						</div>
					{/if}
				</div>
			</div>

			{#if mixedInputMods.length}
				<div class="mod-link-list">
					<h3 style="font-weight: bold; opacity: 0.9;">
						Mixed Input / Controller support Mod{#if mixedInputMods.length > 1}s{/if}
					</h3>
					<div>
						{#each mixedInputMods as mod, idx (idx)}
							<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
							<a href={mod.url} target="_blank" rel="noopener noreferrer"
								>{mod.name
									? `${mod.name} (${getModHostNameNice(mod.url)})`
									: getModUrlName(mod.url)}</a
							>
						{/each}
					</div>
				</div>
			{/if}
		</aside>
	{/if}
{/snippet}

<style lang="postcss">
#controller-support {
	position: relative;
	display: flex;
	flex-flow: row wrap;
	overflow: clip;
	overflow-clip-margin: 2em;

	--info-min-width: 42ch;
	width: 100%;
	align-items: stretch;
	height: 100%;

	gap: 1em;
	& > :first-child {
		flex: 3 1 var(--info-min-width);
		min-width: 0;
		display: flex;
		flex-flow: row wrap;
		gap: 1em;
		justify-content: space-evenly;
		align-items: center;
		& > section {
			margin: 0 2.5em;
		}
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
			& > div {
				display: grid;
				gap: 0.5em;
				margin-left: 1.5em;
			}
			& a {
				width: fit-content;
				font-weight: bold;
			}
		}

		& .scrollable {
			flex: 1 1 auto;
			min-height: 0;
			overflow: auto;
			overflow: auto;
			& > :first-child {
				max-height: 100%;
			}
		}
		max-height: 46em;
	}
}

.notes-container {
	& > :first-child {
		& > h3 {
			position: relative;
			font-weight: bold;
			opacity: 0.9;
			padding-top: 1em;
			filter: drop-shadow(2px 2px 2px rgba(0, 0, 0, 0.226))
				drop-shadow(0px 0px 0.5em rgba(0, 0, 0, 0.342));

			&::before {
				content: '';
				position: absolute;
				top: -0.42em;
				left: 0;
				width: 100%;
				height: 1px;
				background: var(--text-color);
				opacity: 0.2;
			}
		}
		& > h3:first-child {
			padding-top: 0;
			&::before {
				display: none;
			}
		}

		display: grid;
		gap: 1em;
	}
	& .mdcontainer {
		padding: 0 1em;
	}
}

section > :first-child:is(h3) {
	font-weight: bold;
	font-size: 1.3em;
	filter: drop-shadow(2px 2px 2px rgba(0, 0, 0, 0.329)) drop-shadow(0px 0px 0.5em rgba(0, 0, 0, 0.432));
}

h3 {
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
		& > div > :global(:first-child),
		& > div > :nth-child(2) {
			font-size: 1.2em;
			filter: drop-shadow(2px 2px 2px rgba(0, 0, 0, 0.329))
				drop-shadow(0px 0px 0.5em rgba(0, 0, 0, 0.322));
		}
	}
}

.ctrl-support-help {
	position: absolute;
	top: 1em;
	right: 1em;
	z-index: 10;
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
	height: fit-content;
	&:is(:not(.no-divider)) {
		&::after {
			content: '';
			width: 100%;
			height: 1px;
			background: var(--text-color);
			opacity: 0.2;
			margin: 1em 0;
		}
	}
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
