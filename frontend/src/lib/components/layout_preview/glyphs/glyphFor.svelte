<script lang="ts" module>
export { glyphFor };
export const glyphMaps: Record<string, Record<string, string>> = {
	controller_steamcontroller_gordon: {
		left_bumper: 'sc_lb.svg',
		right_bumper: 'sc_rb.svg',
		left_trigger: 'sc_lt.svg',
		// left_trigger_click: 'sc_lt_click.svg',
		left_trigger_click: 'sc_lt.svg',
		left_trigger_edge: 'sc_lt_soft.svg',
		left_trigger_threshold: 'sc_lt_soft.svg',
		right_trigger: 'sc_rt.svg',
		// right_trigger_click: 'rt_click.svg',
		right_trigger_click: 'sc_rt.svg',
		right_trigger_edge: 'sc_rt_soft.svg',
		right_trigger_threshold: 'sc_rt_soft.svg',
		button_back_left: 'sc_lg.svg',
		button_back_right: 'sc_rg.svg',
		button_menu: 'sc_button_l_arrow.svg',
		button_escape: 'sc_button_r_arrow.svg',
		left_trackpad_dpad_north: 'sc_dpad_up.svg',
		left_trackpad_dpad_south: 'sc_dpad_down.svg',
		left_trackpad_dpad_west: 'sc_dpad_left.svg',
		left_trackpad_dpad_east: 'sc_dpad_right.svg',
		left_trackpad_edge: 'empty_Outer Ring',
		left_trackpad_touch: 'sc_dpad_touch.svg',
		left_trackpad_click: 'sc_dpad_click.svg',
		left_click: 'sc_dpad_click.svg',
		left_trackpad_scroll_clockwise: 'empty_',
		left_trackpad_scroll_counterclockwise: 'empty_',
		right_trackpad_dpad_north: 'sc_touchpad_up.svg',
		right_trackpad_dpad_south: 'sc_touchpad_down.svg',
		right_trackpad_dpad_west: 'sc_touchpad_left.svg',
		right_trackpad_dpad_east: 'sc_touchpad_right.svg',
		right_trackpad_click: 'sc_touchpad_click.svg',
		right_click: 'sc_touchpad_click.svg',
		right_trackpad_edge: 'empty_Outer Ring',
		right_trackpad_touch: 'sc_touchpad_touch.svg',
		left_trackpad_doubletap: 'shared_touch_doubletap.svg',
		right_trackpad_doubletap: 'shared_touch_doubletap.svg'
	},
	controller_neptune: {
		left_bumper: 'sd_l1.svg',
		right_bumper: 'sd_r1.svg',
		left_trigger: 'sd_l2.svg',
		left_trigger_click: 'sd_l2.svg',
		left_trigger_edge: 'sd_l2_half.svg',
		left_trigger_threshold: 'sd_l2_half.svg',
		right_trigger: 'sd_r2.svg',
		right_trigger_click: 'sd_r2.svg',
		right_trigger_edge: 'sd_r2_half.svg',
		right_trigger_threshold: 'sd_r2_half.svg',
		button_back_left_upper: 'sd_l4.svg',
		button_back_left: 'sd_l5.svg',
		button_back_right_upper: 'sd_r4.svg',
		button_back_right: 'sd_r5.svg',
		button_diamond_button_a: 'shared_button_a.svg',
		button_diamond_button_b: 'shared_button_b.svg',
		button_diamond_button_x: 'shared_button_x.svg',
		button_diamond_button_y: 'shared_button_y.svg',
		button_a: 'shared_button_a.svg',
		button_b: 'shared_button_b.svg',
		button_x: 'shared_button_x.svg',
		button_y: 'shared_button_y.svg',
		button_menu: 'sd_button_view.svg',
		button_escape: 'sd_button_menu.svg',
		left_trackpad_dpad_north: 'sd_ltrackpad_up.svg',
		left_trackpad_dpad_south: 'sd_ltrackpad_down.svg',
		left_trackpad_dpad_west: 'sd_ltrackpad_left.svg',
		left_trackpad_dpad_east: 'sd_ltrackpad_right.svg',
		left_trackpad_touch: 'sd_ltrackpad_ring.svg',
		left_trackpad_click: 'sd_ltrackpad_click.svg',
		left_click: 'sd_ltrackpad_click.svg',
		right_trackpad_dpad_north: 'sd_rtrackpad_up.svg',
		right_trackpad_dpad_south: 'sd_rtrackpad_down.svg',
		right_trackpad_dpad_west: 'sd_rtrackpad_left.svg',
		right_trackpad_dpad_east: 'sd_rtrackpad_right.svg',
		right_trackpad_click: 'sd_rtrackpad_click.svg',
		right_click: 'sd_rtrackpad_click.svg',
		right_trackpad_touch: 'sd_rtrackpad_ring.svg'
	},
	controller_triton: {
		left_bumper: 'sc_l1.svg',
		right_bumper: 'sc_r1.svg',
		left_trigger_edge: 'sc_l2_half.svg',
		right_trigger_edge: 'sc_r2_half.svg',
		left_trigger_click: 'sc_l2.svg',
		right_trigger_click: 'sc_r2.svg',
		left_trigger: 'sc_l2.svg',
		right_trigger: 'sc_r2.svg',
		button_back_left_upper: 'sc_l4.svg',
		button_back_left: 'sc_l5.svg',
		button_back_right_upper: 'sc_r4.svg',
		button_back_right: 'sc_r5.svg',
		grip_sense_left: 'sd_lgrip.svg',
		grip_sense_right: 'sd_rgrip.svg'
	},
	controller_ps5: {
		left_bumper: 'ps5_l1.svg',
		right_bumper: 'ps5_r1.svg',
		left_trigger: 'ps5_l2.svg',
		left_trigger_click: 'ps5_l2.svg',
		left_trigger_edge: 'ps5_l2_soft.svg',
		left_trigger_threshold: 'ps5_l2_soft.svg',
		right_trigger: 'ps5_r2.svg',
		right_trigger_click: 'ps5_r2.svg',
		right_trigger_edge: 'ps5_r2_soft.svg',
		right_trigger_threshold: 'ps5_r2_soft.svg',
		button_diamond_button_a: 'ps_button_x.svg',
		button_diamond_button_b: 'ps_button_circle.svg',
		button_diamond_button_x: 'ps_button_square.svg',
		button_diamond_button_y: 'ps_button_triangle.svg',
		button_a: 'ps_button_x.svg',
		button_b: 'ps_button_circle.svg',
		button_x: 'ps_button_square.svg',
		button_y: 'ps_button_triangle.svg',
		button_menu: 'ps5_button_create.svg',
		button_escape: 'ps5_button_options.svg',
		left_trackpad_dpad_north: 'ps5_trackpad_l_up.svg',
		left_trackpad_dpad_south: 'ps5_trackpad_l_down.svg',
		left_trackpad_dpad_west: 'ps5_trackpad_l_left.svg',
		left_trackpad_dpad_east: 'ps5_trackpad_l_right.svg',
		left_trackpad_touch: 'ps5_trackpad_l_ring.svg',
		left_trackpad_click: 'ps5_trackpad_l_click.svg',
		left_click: 'ps5_trackpad_l_click.svg',
		right_trackpad_dpad_north: 'ps5_trackpad_r_up.svg',
		right_trackpad_dpad_south: 'ps5_trackpad_r_down.svg',
		right_trackpad_dpad_west: 'ps5_trackpad_r_left.svg',
		right_trackpad_dpad_east: 'ps5_trackpad_r_right.svg',
		right_trackpad_click: 'ps5_trackpad_r_click.svg',
		right_click: 'ps5_trackpad_r_click.svg',
		right_trackpad_touch: 'ps5_trackpad_r_ring.svg',
		dpad_dpad_north: 'ps_dpad_up.svg',
		dpad_dpad_south: 'ps_dpad_down.svg',
		dpad_dpad_west: 'ps_dpad_left.svg',
		dpad_dpad_east: 'ps_dpad_right.svg',
		dpad_dpad: 'ps_dpad.svg',
		left_joystick_click: 'shared_l3.svg',
		right_joystick_click: 'shared_r3.svg',
		button_back_left_upper: 'ps_lb.svg',
		button_back_left: 'ps_lfn.svg',
		button_back_right_upper: 'ps_rb.svg',
		button_back_right: 'ps_rfn.svg'
	},
	controller_ps4: {
		left_bumper: 'ps4_l1.svg',
		right_bumper: 'ps4_r1.svg',
		left_trigger: 'ps4_l2.svg',
		left_trigger_click: 'ps4_l2.svg',
		left_trigger_edge: 'ps4_l2_soft.svg',
		left_trigger_threshold: 'ps4_l2_soft.svg',
		right_trigger: 'ps4_r2.svg',
		right_trigger_click: 'ps4_r2.svg',
		right_trigger_edge: 'ps4_r2_soft.svg',
		right_trigger_threshold: 'ps4_r2_soft.svg',
		button_diamond_button_a: 'ps_color_button_x.svg',
		button_diamond_button_b: 'ps_color_button_circle.svg',
		button_diamond_button_x: 'ps_color_button_square.svg',
		button_diamond_button_y: 'ps_color_button_triangle.svg',
		button_a: 'ps_color_button_x.svg',
		button_b: 'ps_color_button_circle.svg',
		button_x: 'ps_color_button_square.svg',
		button_y: 'ps_color_button_triangle.svg',
		button_menu: 'ps4_button_share.svg',
		button_escape: 'ps4_button_options.svg',
		left_trackpad_dpad_north: 'ps4_trackpad_l_up.svg',
		left_trackpad_dpad_south: 'ps4_trackpad_l_down.svg',
		left_trackpad_dpad_west: 'ps4_trackpad_l_left.svg',
		left_trackpad_dpad_east: 'ps4_trackpad_l_right.svg',
		left_trackpad_touch: 'ps4_trackpad_l_ring.svg',
		left_trackpad_click: 'ps4_trackpad_l_click.svg',
		left_click: 'ps4_trackpad_l_click.svg',
		right_trackpad_dpad_north: 'ps4_trackpad_r_up.svg',
		right_trackpad_dpad_south: 'ps4_trackpad_r_down.svg',
		right_trackpad_dpad_west: 'ps4_trackpad_r_left.svg',
		right_trackpad_dpad_east: 'ps4_trackpad_r_right.svg',
		right_trackpad_click: 'ps4_trackpad_r_click.svg',
		right_click: 'ps4_trackpad_r_click.svg',
		right_trackpad_touch: 'ps4_trackpad_r_ring.svg',
		dpad_dpad_north: 'ps_dpad_up.svg',
		dpad_dpad_south: 'ps_dpad_down.svg',
		dpad_dpad_west: 'ps_dpad_left.svg',
		dpad_dpad_east: 'ps_dpad_right.svg',
		dpad_dpad: 'ps_dpad.svg',
		left_joystick_click: 'shared_l3.svg',
		right_joystick_click: 'shared_r3.svg'
	},
	controller_switch_pro: {
		left_bumper: 'switchpro_l.svg',
		right_bumper: 'switchpro_r.svg',
		left_trigger: 'switchpro_l2.svg',
		left_trigger_click: 'switchpro_l2.svg',
		left_trigger_edge: 'switchpro_l2_soft.svg',
		left_trigger_threshold: 'switchpro_l2_soft.svg',
		right_trigger: 'switchpro_r2.svg',
		right_trigger_click: 'switchpro_r2.svg',
		right_trigger_edge: 'switchpro_r2_soft.svg',
		right_trigger_threshold: 'switchpro_r2_soft.svg',
		button_diamond_button_a: 'shared_outlined_button_a.svg',
		button_diamond_button_b: 'shared_outlined_button_b.svg',
		button_diamond_button_x: 'shared_outlined_button_x.svg',
		button_diamond_button_y: 'shared_outlined_button_y.svg',
		button_a: 'shared_outlined_button_a.svg',
		button_b: 'shared_outlined_button_b.svg',
		button_x: 'shared_outlined_button_x.svg',
		button_y: 'shared_outlined_button_y.svg',
		button_menu: 'switchpro_button_minus.svg',
		button_escape: 'switchpro_button_plus.svg',
		dpad_dpad_north: 'switchpro_dpad_up.svg',
		dpad_dpad_south: 'switchpro_dpad_down.svg',
		dpad_dpad_west: 'switchpro_dpad_left.svg',
		dpad_dpad_east: 'switchpro_dpad_right.svg',
		dpad_dpad: 'switchpro_dpad.svg',
		button_back_left: 'switchpro_gl.svg',
		button_back_right: 'switchpro_gr.svg',
		joystick_dpad_north: 'switchpro_lstick_up.svg',
		joystick_dpad_south: 'switchpro_lstick_down.svg',
		joystick_dpad_west: 'switchpro_lstick_left.svg',
		joystick_dpad_east: 'switchpro_lstick_right.svg',
		right_joystick_dpad_north: 'switchpro_rstick_up.svg',
		right_joystick_dpad_south: 'switchpro_rstick_down.svg',
		right_joystick_dpad_west: 'switchpro_rstick_left.svg',
		right_joystick_dpad_east: 'switchpro_rstick_right.svg',
		joystick_click: 'switchpro_lstick_click.svg',
		left_joystick_click: 'switchpro_lstick_click.svg',
		right_joystick_click: 'switchpro_rstick_click.svg'
	},
	controller_xboxone: {
		button_menu: 'xbox_button_select.svg',
		button_escape: 'xbox_button_start.svg'
	},
	fallback: {
		left_bumper: 'xbox_lb.svg',
		right_bumper: 'xbox_rb.svg',
		left_trigger: 'xbox_lt.svg',
		left_trigger_click: 'xbox_lt.svg',
		left_trigger_edge: 'xbox_lt_soft.svg',
		left_trigger_threshold: 'xbox_lt_soft.svg',
		right_trigger: 'xbox_rt.svg',
		right_trigger_click: 'xbox_rt.svg',
		right_trigger_edge: 'xbox_rt_soft.svg',
		right_trigger_threshold: 'xbox_rt_soft.svg',
		button_menu: 'sc_button_l_arrow.svg',
		button_escape: 'sc_button_r_arrow.svg',
		joystick_dpad_north: 'shared_lstick_up.svg',
		joystick_dpad_south: 'shared_lstick_down.svg',
		joystick_dpad_west: 'shared_lstick_left.svg',
		joystick_dpad_east: 'shared_lstick_right.svg',
		right_joystick_dpad_north: 'shared_rstick_up.svg',
		right_joystick_dpad_south: 'shared_rstick_down.svg',
		right_joystick_dpad_west: 'shared_rstick_left.svg',
		right_joystick_dpad_east: 'shared_rstick_right.svg',
		joystick_click: 'shared_lstick_click.svg',
		left_joystick_click: 'shared_lstick_click.svg',
		right_joystick_click: 'shared_rstick_click.svg',
		joystick_touch: 'shared_lstick_touch.svg',
		left_joystick_touch: 'shared_lstick_touch.svg',
		right_joystick_touch: 'shared_lstick_touch.svg',
		joystick_edge: 'empty_Outer Ring',
		button_diamond_button_a: 'shared_color_button_a.svg',
		button_diamond_button_b: 'shared_color_button_b.svg',
		button_diamond_button_x: 'shared_color_button_x.svg',
		button_diamond_button_y: 'shared_color_button_y.svg',
		button_a: 'shared_color_button_a.svg',
		button_b: 'shared_color_button_b.svg',
		button_x: 'shared_color_button_x.svg',
		button_y: 'shared_color_button_y.svg',
		gyro_dpad_west: 'empty_Lean Left',
		gyro_dpad_east: 'empty_Lean Right',
		gyro_dpad_north: 'empty_',
		gyro_dpad_south: 'empty_',
		left_stick_touch: 'shared_lstick_touch.svg',
		right_stick_touch: 'shared_rstick_touch.svg',
		dpad_dpad_north: 'shared_dpad_up.svg',
		dpad_dpad_south: 'shared_dpad_down.svg',
		dpad_dpad_west: 'shared_dpad_left.svg',
		dpad_dpad_east: 'shared_dpad_right.svg',
		dpad_dpad: 'shared_dpad.svg',
		left_joystick: 'shared_lstick.svg',
		right_joystick: 'shared_rstick.svg',
		button_back_left_upper: 'sd_l4.svg',
		button_back_left: 'sd_l5.svg',
		button_back_right_upper: 'sd_r4.svg',
		button_back_right: 'sd_r5.svg',
		grip_sense_left: 'sd_lgrip.svg',
		grip_sense_right: 'sd_rgrip.svg',
		button_leftauxcapsense: 'sd_lgrip.svg',
		button_rightauxcapsense: 'sd_rgrip.svg'
	}
};

glyphMaps.fallback = {
	...glyphMaps.controller_steamcontroller_gordon,
	...glyphMaps.fallback
};
glyphMaps.controller_triton = {
	...glyphMaps.controller_neptune,
	...glyphMaps.controller_triton
};
glyphMaps.controller_ps5_edge = {
	...glyphMaps.controller_ps5
};
glyphMaps.controller_ps3 = {
	...glyphMaps.controller_ps4
};
glyphMaps.controller_8bitdo = {
	...glyphMaps.controller_switch_pro
};
glyphMaps.controller_switch2_pro = {
	...glyphMaps.controller_switch_pro
};
glyphMaps.controller_xboxelite = {
	...glyphMaps.controller_xboxone
};
</script>

{#snippet glyphFor(controllerType: string, device: string, originalController?: string)}
	{@const deviceLower = device.toLowerCase()}
	{@const glyphEntry = glyphMaps[controllerType]?.[deviceLower]}
	{#if !glyphEntry}
		{#if controllerType === 'fallback'}
			{device}
		{:else}
			{@render glyphFor('fallback', device, originalController ?? controllerType)}
		{/if}
	{:else if glyphEntry.startsWith('empty')}
		<span style="white-space: nowrap;"><b>{glyphEntry.split('_')?.pop()}</b></span>
	{:else}
		<svelte:boundary>
			<!-- why the suffix? my ide shows previews above, and vite needs the extension for dynamic imports;
         It's a small price to pay! -->
			{@const Svg = (await import(`./${glyphEntry.slice(0, -4)}.svg?component`))?.default}
			<Svg class={`${glyphEntry.includes('color') ? 'color-glyph' : 'glyph'}`} />
			{#snippet failed()}
				{#if controllerType === 'fallback'}
					{device}
				{:else}
					{@render glyphFor('fallback', device, originalController ?? controllerType)}
				{/if}
			{/snippet}
		</svelte:boundary>
	{/if}
{/snippet}
