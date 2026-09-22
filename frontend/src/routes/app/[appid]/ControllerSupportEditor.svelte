<script lang="ts">
import type { components } from '$lib/api/openapi';
import {
	AppGlyphTagType,
	HWFeature,
	HWFeatureControllerType,
	HWFeaturesForFamily,
	SteamInputAPISupportType
} from '$lib/api/steaminputdbEnums';
import { CONTROLLER_LIST } from '$lib/components/search/controllerlist.svelte';

import IcoHaptics from '$lib/assets/icohaptics.svg?component';
import IcoMixedInput from '$lib/assets/mixedinput.svg?component';
import IcoDs5 from '$lib/assets/steam_controller_type_svgs/ps5.svg?component';
import IcoSIAPI from '$lib/assets/steam_controller_type_svgs/siapi.svg?component';
import IcoSwitch from '$lib/assets/steam_controller_type_svgs/switchpro.svg?component';
import IcoXbone from '$lib/assets/steam_controller_type_svgs/xbox.svg?component';
import IcoDotsPer360 from '$lib/components/layout_preview/glyphs/cmnd_dots_per_360_calibration_spin.svg?component';
import IcoTrigger from '$lib/components/layout_preview/glyphs/shared_trigger.svg?component';
import IcoTouchpads from '~icons/fluent/cursor-hover-16-filled';
import IcoRumble from '~icons/fluent/haptic-strong-16-filled';
import IcoGyro from '~icons/game-icons/gyroscope';
import IcoActionSets from '~icons/material-symbols/layers-rounded';
import IcoDropdown from '~icons/mdi/chevron-down';
import IcoMdiCross from '~icons/mdi/close-circle-outline';
import IconMdiGamepad from '~icons/mdi/controller';
import IcoDpad from '~icons/mdi/gamepad';
import IconMdiGamepadCircle from '~icons/mdi/gamepad-circle';
import IcoLightbar from '~icons/mdi/lightbulb-on';
import IcoMdiDash from '~icons/mdi/minus-circle-outline';
import IcoPlus from '~icons/mdi/plus';
import IcoMdiSetRight from '~icons/mdi/set-right';
import IcoSteam from '~icons/mdi/steam';
import IcoTrash from '~icons/mdi/trash-can-outline';
import IcoAudioHaptics from '~icons/mdi/volume-vibrate';

import { enhance } from '$app/forms';
import { invalidateAll } from '$app/navigation';
import { toast } from '$lib/toaster/toaster.svelte';
import IcoCross from '~icons/mdi/close';
import IcoEye from '~icons/mdi/eye';
import IcoFloppy from '~icons/mdi/floppy';
import IcoPencil from '~icons/mdi/pencil';
import ControllerSupport from './ControllerSupport.svelte';

let {
	appInfo = $bindable(),
	// eslint-disable-next-line no-useless-assignment
	editControllerSupport = $bindable<boolean>(false)
}: {
	appInfo: components['schemas']['AppInfoItem'];
	editControllerSupport: boolean;
} = $props();
let steaminputdbInfo = $derived.by(() => {
	const draft = $state({
		...appInfo
	});
	if (!draft?.steaminputdb_info) {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		draft.steaminputdb_info = {} as any;
	}
	if (!draft?.steaminputdb_info?.mixed_input) {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		draft.steaminputdb_info!.mixed_input = { type: 0 } as any;
	}
	if (!draft?.steaminputdb_info?.glyphs) {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		draft.steaminputdb_info!.glyphs = {} as any;
	}

	if (!draft?.steaminputdb_info?.steaminputapi_support) {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		draft.steaminputdb_info!.steaminputapi_support = { camera_support: 0 } as any;
	}

	return draft?.steaminputdb_info as Exclude<
		Required<components['schemas']['AppInfoItem']['steaminputdb_info']>,
		undefined
	>;
});

let preview = $state(false);
let previewAppInfo = $derived({
	...appInfo,
	steaminputdb_info: {
		...steaminputdbInfo
	}
});
</script>

<form
	id="controller-support"
	method="POST"
	action="?/patchControllerSupportInfo"
	use:enhance={({ formData }) => {
		formData.set('info', JSON.stringify(previewAppInfo.steaminputdb_info));
		appInfo = previewAppInfo;
		editControllerSupport = false;
		return async ({ result }) => {
			if (result.type !== 'success') {
				toast({
					message: 'Failed to save controller support info',
					color: 'firebrick'
				});
			} else {
				editControllerSupport = false;
				window.scrollTo(0, 0);
				toast({
					message: 'Successfully saved controller support info',
					color: 'green'
				});
				await invalidateAll();
			}
		};
	}}
>
	{#if preview}
		<strong
			style="color: red; font-size: 1.4em; text-align: center; filter: drop-shadow(0 0 0.75rem black);"
			>Preview Mode</strong
		>
		<ControllerSupport appInfo={previewAppInfo} />
	{:else}
		{@render controllerSupportContent()}
	{/if}
	<div class="form-buttons">
		{#if !preview}
			<button
				class="edit-info-button"
				type="button"
				onclick={() => {
					preview = true;
					window.scrollTo(0, 0);
				}}
			>
				<IcoEye style="width: 1.2em; height: 1.2em; " />Preview
			</button>
		{:else}
			<button class="edit-info-button" type="button" onclick={() => (preview = false)}>
				<IcoPencil style="width: 1.2em; height: 1.2em; " />Edit
			</button>
		{/if}
		<button class="edit-info-button" type="button" onclick={() => (editControllerSupport = false)}>
			<IcoCross style="width: 1.2em; height: 1.2em; " />Cancel Edit
		</button>
		<button class="edit-info-button">
			<IcoFloppy style="width: 1.2em; height: 1.2em; " />Save
		</button>
	</div>
</form>

{#snippet controllerGlyphSelect(type: string | undefined, niceName: string)}
	{#each CONTROLLER_LIST.filter((c) => {
		return c.type === type;
	}) as controller (controller.type)}
		<div
			style="display: grid; gap: 0.25em; place-items: center; border: 1px solid #ccc2; padding: 0.5em; border-radius: 0.5em;"
		>
			<label for={`glyph-controller-${type}`} style="margin-bottom: 0.5em">
				<input
					type="checkbox"
					id={`glyph-controller-${type}`}
					name={`glyph-controller-${type}`}
					checked={steaminputdbInfo.glyphs.controllers?.some((t) => t.controller_type === type)}
					onchange={() => {
						steaminputdbInfo.glyphs.controllers = steaminputdbInfo.glyphs.controllers ?? [];
						const index = steaminputdbInfo.glyphs.controllers.findIndex(
							(t) => t.controller_type === type
						);
						if (index === -1) {
							steaminputdbInfo.glyphs.controllers.push({ controller_type: type });
						} else {
							steaminputdbInfo.glyphs.controllers.splice(index, 1);
						}
					}}
				/>
				{#if type?.includes('xbox')}
					<IcoXbone style="width: 1.4em; height: 1.4em;" />
				{:else}
					<controller.icon style="width: 1.4em; height: 1.4em;" />
				{/if}
				<span>
					{#if type?.includes?.('xbox')}
						Xbox
					{:else if type?.includes?.('switch')}
						Switch
					{:else}
						{niceName}
					{/if}
				</span>
			</label>
			<input
				type="text"
				placeholder="Tooltip notes..."
				disabled={steaminputdbInfo.glyphs.controllers?.find((t) => t.controller_type === type) ===
					undefined}
				value={steaminputdbInfo.glyphs.controllers?.find((t) => t.controller_type === type)?.notes ??
					''}
				onchange={(e) => {
					const controller = steaminputdbInfo.glyphs.controllers?.find(
						(t) => t.controller_type === type
					);
					if (controller) {
						// eslint-disable-next-line @typescript-eslint/no-explicit-any
						controller.notes = (e.target as any)?.value;
					}
				}}
			/>
		</div>
	{/each}
{/snippet}

{#snippet steamInputAPIGlyphTag(tag: number, disabled: boolean)}
	<label
		for={`siapi-glyph-${tag}`}
		class={disabled ? 'feature disabled' : 'feature '}
		style="background-color: {tag === AppGlyphTagType.SteamInputAPIButtons
			? '#00dfff'
			: tag === AppGlyphTagType.InGameButtons
				? '#12ff12'
				: 'red'}; gap: 0.5em;"
	>
		<input
			type="checkbox"
			id={`siapi-glyph-${tag}`}
			name={`siapi-glyph-${tag}`}
			disabled={disabled}
			checked={steaminputdbInfo.steaminputapi_support?.glyphs?.includes(tag)}
			onchange={() => {
				steaminputdbInfo.steaminputapi_support = steaminputdbInfo.steaminputapi_support ?? {};
				steaminputdbInfo.steaminputapi_support.glyphs =
					steaminputdbInfo.steaminputapi_support.glyphs ?? [];
				const index = steaminputdbInfo.steaminputapi_support.glyphs.findIndex((t) => t === tag);
				if (index === -1) {
					steaminputdbInfo.steaminputapi_support.glyphs.push(tag);
				} else {
					steaminputdbInfo.steaminputapi_support.glyphs.splice(index, 1);
				}
			}}
		/>
		{#if tag === AppGlyphTagType.SteamInputAPIButtons}
			<IcoSIAPI style="width: 1.4em; height: 1.4em;" /> Steam Input API Glyphs
		{:else if tag === AppGlyphTagType.InGameButtons}
			<IconMdiGamepad style="width: 1.2em; height: 1.2em;" /> In-Game Glyphs
		{:else if tag === AppGlyphTagType.SteamInputAPITexts}
			<IcoMdiDash style="width: 1.2em; height: 1.2em;" /> Steam Input API Button Text
		{/if}
	</label>
{/snippet}

{#snippet hwFeature(feature: number, family: number)}
	{@const entry = steaminputdbInfo.hw_features?.find(
		(f) => f.feature === feature && f.controller_family === family
	)}
	<div style="display: grid; gap: 0.25em; justify-items: center; align-self: start;">
		<label for={`hw-feature-${family}-${feature}`} class="feature">
			<input
				type="checkbox"
				id={`hw-feature-${family}-${feature}`}
				name={`hw-feature-${family}-${feature}`}
				checked={!!entry}
				onchange={() => {
					steaminputdbInfo.hw_features = steaminputdbInfo.hw_features ?? [];
					const index = steaminputdbInfo.hw_features.findIndex(
						(f) => f.feature === feature && f.controller_family === family
					);
					if (index === -1) {
						steaminputdbInfo.hw_features.push({ feature, controller_family: family });
					} else {
						steaminputdbInfo.hw_features.splice(index, 1);
					}
				}}
			/>
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
			{/if}
		</label>
		{#if entry}
			<input
				type="text"
				placeholder="Tooltip Notes..."
				style="width: 0; min-width: 100%; box-sizing: border-box;"
				value={entry.notes ?? ''}
				onchange={(e) => {
					entry.notes = e.currentTarget.value;
				}}
			/>
		{/if}
	</div>
{/snippet}

{#snippet steamInputAPISupportTags(tag: number)}
	{#if tag === SteamInputAPISupportType.InGameActions}
		<label for={`siapi-support-${tag}`} class="feature" style="background-color: turquoise; gap: 0.5em;">
			<input
				type="checkbox"
				id={`siapi-support-${tag}`}
				name={`siapi-support-${tag}`}
				checked={steaminputdbInfo.steaminputapi_support?.support_tags?.includes(tag)}
				onchange={() => {
					steaminputdbInfo.steaminputapi_support = steaminputdbInfo.steaminputapi_support ?? {};
					steaminputdbInfo.steaminputapi_support.support_tags =
						steaminputdbInfo.steaminputapi_support.support_tags ?? [];
					const index = steaminputdbInfo.steaminputapi_support.support_tags.findIndex(
						(t) => t === tag
					);
					if (index === -1) {
						steaminputdbInfo.steaminputapi_support.support_tags.push(tag);
					} else {
						steaminputdbInfo.steaminputapi_support.support_tags.splice(index, 1);
					}
				}}
			/>
			<IcoMdiSetRight style="width: 1.2em; height: 1.2em;" />
			<span> In-Game Actions </span>
		</label>
	{:else if tag === SteamInputAPISupportType.ActionSets}
		<label for={`siapi-support-${tag}`} class="feature" style="background-color: yellowgreen">
			<input
				type="checkbox"
				id={`siapi-support-${tag}`}
				name={`siapi-support-${tag}`}
				checked={steaminputdbInfo.steaminputapi_support?.support_tags?.includes(tag)}
				onchange={() => {
					steaminputdbInfo.steaminputapi_support = steaminputdbInfo.steaminputapi_support ?? {};
					steaminputdbInfo.steaminputapi_support.support_tags =
						steaminputdbInfo.steaminputapi_support.support_tags ?? [];
					const index = steaminputdbInfo.steaminputapi_support.support_tags.findIndex(
						(t) => t === tag
					);
					if (index === -1) {
						steaminputdbInfo.steaminputapi_support.support_tags.push(tag);
					} else {
						steaminputdbInfo.steaminputapi_support.support_tags.splice(index, 1);
					}
				}}
			/>
			<IcoActionSets style="width: 1.2em; height: 1.2em;" /> Action Sets
		</label>
	{:else if tag === SteamInputAPISupportType.XInputActions}
		<label for={`siapi-support-${tag}`} class="feature" style="background-color: lightgreen; gap: 0.5em;">
			<input
				type="checkbox"
				id={`siapi-support-${tag}`}
				name={`siapi-support-${tag}`}
				checked={steaminputdbInfo.steaminputapi_support?.support_tags?.includes(tag)}
				onchange={() => {
					steaminputdbInfo.steaminputapi_support = steaminputdbInfo.steaminputapi_support ?? {};
					steaminputdbInfo.steaminputapi_support.support_tags =
						steaminputdbInfo.steaminputapi_support.support_tags ?? [];
					const index = steaminputdbInfo.steaminputapi_support.support_tags.findIndex(
						(t) => t === tag
					);
					if (index === -1) {
						steaminputdbInfo.steaminputapi_support.support_tags.push(tag);
					} else {
						steaminputdbInfo.steaminputapi_support.support_tags.splice(index, 1);
					}
				}}
			/>
			<IcoXbone style="width: 1.2em; height: 1.2em;" /> XInput Style Actions
		</label>
	{:else if tag === SteamInputAPISupportType.Discontinued}
		<label for={`siapi-support-${tag}`} class="feature" style="background-color: red; gap: 0.5em;">
			<input
				type="checkbox"
				id={`siapi-support-${tag}`}
				name={`siapi-support-${tag}`}
				checked={steaminputdbInfo.steaminputapi_support?.support_tags?.includes(tag)}
				onchange={() => {
					steaminputdbInfo.steaminputapi_support = steaminputdbInfo.steaminputapi_support ?? {};
					steaminputdbInfo.steaminputapi_support.support_tags =
						steaminputdbInfo.steaminputapi_support.support_tags ?? [];
					const index = steaminputdbInfo.steaminputapi_support.support_tags.findIndex(
						(t) => t === tag
					);
					if (index === -1) {
						steaminputdbInfo.steaminputapi_support.support_tags.push(tag);
					} else {
						steaminputdbInfo.steaminputapi_support.support_tags.splice(index, 1);
					}
				}}
			/>
			<IcoMdiCross style="width: 1.2em; height: 1.2em;" /> Discontinued
		</label>
	{/if}
{/snippet}

{#snippet controllerSupportContent()}
	<div class="card glass">
		<section class="group">
			<section id="mixed-input" class="info-group">
				<h3><IcoMixedInput style="width: 2em; height: 2em;" />Mixed Input</h3>
				<div>
					<div>
						<label for="mixed-input-type" class="dropdown">
							<span>Mixed Input: </span>
							<select
								id="mixed-input-type"
								name="mixed-input-type"
								bind:value={steaminputdbInfo.mixed_input.type}
							>
								<option value={0}>Unknown</option>
								<option value={1}>Unsupported</option>
								<option value={2}>Supported</option>
								<option value={3}>Requires Mod</option>
								<option value={4}>Partial</option>
							</select>

							<IcoDropdown />
						</label>
						<label for="glypth-flicker">
							<label for="mixed-input-type" class="dropdown">
								<span>Glyph Flicker: </span>
								<select
									id="glyph-flicker"
									name="glyph-flicker"
									bind:value={steaminputdbInfo.mixed_input.glyph_flicker}
								>
									<option value={null}>Unknown</option>
									<option value={true}>Yes</option>
									<option value={false}>No</option>
								</select>
								<IcoDropdown />
							</label>
						</label>
					</div>
					<div class="notes">
						<span>Mixed Input specific notes:</span>
						<textarea
							id="mixed-input-notes"
							name="mixed-input-notes"
							bind:value={steaminputdbInfo.mixed_input.notes}></textarea>
					</div>
				</div>
			</section>
			<section id="glyphs" class="info-group">
				<h3><IconMdiGamepadCircle style="width: 1.6em; height: 1.6em;" /> Glyphs</h3>
				<div>
					<div class="ctrl-glyphs">
						<!-- eslint-disable prettier/prettier -->
						{#each CONTROLLER_LIST.filter((c) => {
                            return c.type === 'controller_xboxone'
                                || c.type === 'controller_triton'
                                || c.type === 'controller_steamcontroller_gordon'
                                || c.type === 'controller_neptune'
                                || c.type === 'controller_steamframe_pair'
                                || c.type === 'controller_ps5'
                                || c.type === 'controller_ps4'
                                || c.type === 'controller_switch_pro'
                                || c.type === 'controller_generic'
                                ;				
						}) as controller (controller.type)}
							{@render controllerGlyphSelect(controller.type, controller.niceName)}
						{/each}
                    <!-- eslint-enable prettier/prettier -->
					</div>
					<div>
						<label for="glyph-detect">
							<label for="glyph-detect" class="dropdown">
								<span>Glyph Autodetect: </span>
								<select
									id="glyph-detect"
									name="glyph-detect"
									bind:value={steaminputdbInfo.glyphs.autodetect}
								>
									<option value={null}>Unknown</option>
									<option value={true}>Yes</option>
									<option value={false}>No</option>
								</select>
								<IcoDropdown />
							</label>
						</label>
						<label for="manual-glyph-select">
							<label for="manual-glyph-select" class="dropdown">
								<span>Manual Override / Lock: </span>
								<select
									id="manual-glyph-select"
									name="manual-glyph-select"
									bind:value={steaminputdbInfo.glyphs.manual_select}
								>
									<option value={null}>Unknown</option>
									<option value={true}>Yes</option>
									<option value={false}>No</option>
								</select>
								<IcoDropdown />
							</label>
						</label>
					</div>
					{#if (steaminputdbInfo.steaminputapi_support?.steam_input_type ?? 0) <= 1}
						<em style="margin-top: 1em;"
							>Requires Steam Input Type to be other than <code>Unknown</code> or
							<code>Steam Input unaware</code>
						</em>
					{/if}
					<div>
						{#each Object.values(AppGlyphTagType).filter((t) => {
							return t !== AppGlyphTagType.Unknown;
						}) as tag (tag)}
							{@render steamInputAPIGlyphTag(
								tag,
								(steaminputdbInfo.steaminputapi_support?.steam_input_type ?? 0) <= 1
							)}
						{/each}
					</div>
					<div class="notes">
						<span>Glyph specific notes:</span>
						<textarea
							id="glyph-notes"
							name="glyph-notes"
							bind:value={steaminputdbInfo.glyphs.notes}></textarea>
					</div>
				</div>
			</section>
		</section>

		<div class="group">
			<section id="hw-features" class="info-group" style="margin-top: 1em;">
				<h3><IconMdiGamepad style="width: 1.6em; height: 1.6em;" /> Hardware Features</h3>
				<div>
					{#each Object.values(HWFeatureControllerType) as family_tag (family_tag)}
						{#if family_tag == HWFeatureControllerType.Common}
							<div>
								<IcoDpad style="width: 1.4em;" />
								<span>Common</span>
							</div>
							<div>
								{#each HWFeaturesForFamily[family_tag] ?? [] as feature (feature)}
									{@render hwFeature(feature, family_tag)}
								{/each}
							</div>
						{:else if family_tag == HWFeatureControllerType.Steam}
							<div class="feature-group-header">
								<IcoSteam style="width: 1.4em;" />
								<span>Steam</span>
							</div>
							<div>
								{#each HWFeaturesForFamily[family_tag] ?? [] as feature (feature)}
									{@render hwFeature(feature, family_tag)}
								{/each}
							</div>
						{:else if family_tag == HWFeatureControllerType.PlayStation}
							<div class="feature-group-header">
								<IcoDs5 style="width: 1.2em" />
								<span>PlayStation</span>
							</div>
							<div>
								{#each HWFeaturesForFamily[family_tag] ?? [] as feature (feature)}
									{@render hwFeature(feature, family_tag)}
								{/each}
							</div>
						{:else if family_tag == HWFeatureControllerType.Xbox}
							<div class="feature-group-header">
								<IcoXbone style="width: 1.2em;" />
								<span>Xbox</span>
							</div>
							<div>
								{#each HWFeaturesForFamily[family_tag] ?? [] as feature (feature)}
									{@render hwFeature(feature, family_tag)}
								{/each}
							</div>
						{:else if family_tag == HWFeatureControllerType.Nintendo}
							<div class="feature-group-header">
								<IcoSwitch style="width:1.2em;" />
								<span>Nintendo</span>
							</div>
							<div>
								{#each HWFeaturesForFamily[family_tag] ?? [] as feature (feature)}
									{@render hwFeature(feature, family_tag)}
								{/each}
							</div>
						{/if}
					{/each}
				</div>
			</section>
			<section id="steaminputapi" class="info-group" style="margin-bottom: 1em">
				<h3><IcoSIAPI style="width: 2em; height: 2em;" /> Steam Input Support</h3>
				<div>
					<label for="steam-input-type">
						<label for="steam-input-type" class="dropdown">
							<span>Steam Input Type: </span>
							<select
								id="steam-input-type"
								name="steam-input-type"
								bind:value={steaminputdbInfo.steaminputapi_support.steam_input_type}
							>
								<option value={0}>Unknown</option>
								<option value={1}>Steam Input unaware</option>
								<option value={2}>Steam Virtual Gamepad</option>
								<option value={3}>Native Steam Input API</option>
								<option value={4}>Requires Steam Input</option>
							</select>
							<IcoDropdown />
						</label>
					</label>
				</div>
				<div>
					<div>
						{#each Object.values(SteamInputAPISupportType) as tag (tag)}
							{@render steamInputAPISupportTags(tag)}
						{/each}
					</div>
				</div>
				<div>
					<label for="siapi-camera-support">
						<label for="siapi-camera-support" class="dropdown">
							<span>SIAPI Camera Support: </span>
							<select
								id="siapi-camera-support"
								name="siapi-camera-support"
								bind:value={steaminputdbInfo.steaminputapi_support.camera_support}
							>
								<option value={0}>Unknown</option>
								<option value={1}>None</option>
								<option value={2}>Partial</option>
								<option value={3}>Full</option>
							</select>
							<IcoDropdown />
						</label>
					</label>
					<div style="gap: 1em;">
						<span style="display: flex; align-items: center; gap: 0.5ch;"
							><IcoDotsPer360 style="width: 1.6em; height: 1.6em;" />Pixels Per 360°
						</span>
						<input
							type="text"
							placeholder="Native"
							bind:value={steaminputdbInfo.steaminputapi_support.pixels_per_360}
						/>
					</div>
				</div>

				<div class="notes" style="margin-top: 1em;">
					<span>Steam Input API specific notes:</span>
					<textarea
						id="siapi-notes"
						name="siapi-notes"
						bind:value={steaminputdbInfo.steaminputapi_support.notes}></textarea>
				</div>
			</section>
		</div>
	</div>
	<aside id="controller-support-notes" class="card glass">
		<h3>Additional Info</h3>
		<span>Markdown is supported!</span>
		<a href="https://www.markdownguide.org/cheat-sheet/" target="_blank" rel="noopener noreferrer"
			>Markdown Cheat Sheet</a
		>
		<div class="notes" style="height: 100%;">
			<textarea
				placeholder="Any additional info here..."
				style="width: 0; min-width: 100%; height: max-content; resize: none;"
				bind:value={steaminputdbInfo.controller_support_notes}></textarea>
		</div>

		<div class="mod-link-list">
			<strong style="padding-top: 1.5em;">Mixed Input / Controller support Mods </strong>
			{#each steaminputdbInfo.mixed_input.mixed_input_mod_urls ?? [] as url, idx (idx)}
				<div style="display: flex; gap: 0.5em; align-items: center;">
					<!-- eslint-disable-next-line -->
					<a href={url}
						target="_blank"
						rel="noopener noreferrer"
						style="flex: 1 1 auto; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;"
					>
						{url.replace(/^https?:\/\/(www\.)?/i, '')}
					</a>
					<button
						type="button"
						style="padding: 0.75em;"
						onclick={() => {
							steaminputdbInfo.mixed_input.mixed_input_mod_urls?.splice(idx, 1);
						}}
					>
						<IcoTrash style="width: 1.2em; height: 1.2em;" />
					</button>
				</div>
			{/each}
			<div style="display: flex; gap: 0.5em; align-items: center;">
				<input type="url" placeholder="https://..." style="flex: 1 1 auto; min-width: 0;" />
				<button
					type="button"
					style="padding: 0.75em;"
					onclick={(e) => {
						const input = e.currentTarget.previousElementSibling as HTMLInputElement;
						let url = input.value.trim();
						if (!url) {
							return;
						}
						if (!/^https?:\/\//i.test(url)) {
							url = `https://${url}`;
						}
						steaminputdbInfo.mixed_input.mixed_input_mod_urls =
							steaminputdbInfo.mixed_input.mixed_input_mod_urls ?? [];
						steaminputdbInfo.mixed_input.mixed_input_mod_urls.push(url);
						input.value = '';
					}}
				>
					<IcoPlus style="width: 1.2em; height: 1.2em;" />
				</button>
			</div>
		</div>
	</aside>
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
	}
}

h3 {
	font-size: 1.1em;
	align-items: center;
	display: flex;
	gap: 0.5ch;
}

.feature-group-header {
	margin-top: 1em;
}
.rule-divider {
	opacity: 0.1;
	height: 1px;
	background: var(--text-color);
}

.group {
	display: flex;
	flex-flow: column wrap;
	gap: 1em;
	justify-content: space-evenly;
	align-items: center;
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
			gap: 0.5em;
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
			gap: 0.5em;
			align-items: center;
		}
	}
}
.notes {
	display: flex !important;
	flex-flow: column !important;
	gap: 0.5ch;
	width: 100%;
	min-width: 100%;
	justify-content: center;
	margin-top: 1em;
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

	&:is(.disabled) {
		opacity: 0.5;
		cursor: not-allowed;
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

label {
	display: flex;
	align-items: center;
	gap: 1ch;
}

label.dropdown {
	margin-left: auto;
	display: grid;
	grid-template-columns: auto auto;
	gap: 0.5em;
	align-items: center;
	position: relative;
	isolation: isolate;
	border: 1px solid color-mix(in srgb, var(--text-color), transparent 90%);
	padding: 0.5em 1em;
	box-shadow: 0 1px 4px 0 rgb(0 0 0 / 0.25);
	border-radius: 100vw;
	transition: all var(--transition-duration) var(--default-ease);

	& :global([disabled]) {
		opacity: 0.5;
	}

	&:hover,
	&:focus-within {
		outline: 0.1em solid var(--color-primary);
		box-shadow: 0 0 1.3em -0.4em var(--color-primary);
	}

	& > :first-child {
		white-space: nowrap;
	}
	:global(> :last-child) {
		content: '';
		color: var(--text-color);
		position: absolute;
		z-index: 1;
		height: 100%;
		width: 1.4em;
		top: 50%;
		translate: 0 -50%;
		right: 0.5em;
		background-size: contain;
		pointer-events: none;
	}
}

select {
	font-style: inherit;
	background: transparent;
	border: 1px solid transparent;
	outline: none;
	color: var(--text-color);
	cursor: pointer;
	appearance: none;
	padding-right: 1em;
	position: relative;
	width: 100%;

	& option {
		color: var(--text-color);
		background: var(--card-color);
	}
}

textarea {
	width: 100%;
	min-height: 100%;
	resize: vertical;
}

button {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 0.5ch;
	font-weight: bold;
	padding: 0.5em 1em;
}

.form-buttons {
	flex: 1 0 auto;
	display: grid;
	grid-auto-flow: column;
	gap: 1em;
	width: 100%;
	justify-content: end;
}
</style>
