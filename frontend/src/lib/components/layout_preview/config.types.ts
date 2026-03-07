export interface LayoutFile {
    controller_mappings: ControllerMappings;
}

export interface ControllerMappings {
    version: string;
    revision: string;
    title: string;
    description?: string;
    creator?: string;
    controller_type: string;
    major_revision?: string;
    minor_revision?: string;
    Timestamp?: string;
    actions?: Record<string, Action>;
    action_layers?: Record<string, ActionLayer>;
    localization?: Record<string, LocalizationEntry>;
    group: ControllerGroup | ControllerGroup[];
    preset: ControllerPreset | ControllerPreset[];
    settings?: ControllerMappingsSettings;
}

export interface Action {
    title?: string;
    legacy_set?: string;
    [key: string]: string | undefined;
}

export interface ActionLayer extends Action {
    set_layer?: string;
    parent_set_name?: string;
}

export interface LocalizationEntry {
    title?: string;
    description?: string;
    [key: string]: string | undefined;
}

export type GroupMode =
    | 'four_buttons'
    | 'dpad'
    | 'absolute_mouse'
    | 'trigger'
    | 'mouse_joystick'
    | 'joystick_move'
    | 'touch_menu'
    | 'mouse_region'
    | 'scrollwheel'
    | 'radial_menu'
    | 'single_button'
    | 'switches'
    | 'joystick_mouse'
    | (string & {});

export interface ControllerGroup {
    id: string;
    mode: GroupMode;
    inputs?: Record<string, GroupInput>;
    settings?: GroupSettings;
}

export interface GroupInput {
    activators?: ActivatorMap;
    [key: string]: unknown;
}

export type ActivatorType =
    | 'Full_Press'
    | 'Long_Press'
    | 'Start_Press'
    | (string & {});

export type ActivatorMap = {
    [key in ActivatorType]?: Activator | Activator[];
};

export interface Activator {
    bindings?: BindingBlock;
    settings?: ActivatorSettings;
    [key: string]: unknown;
}

export interface BindingBlock {
    binding?: string | string[];
    [key: string]: string | string[] | undefined;
}

export interface ActivatorSettings {
    [key: string]: string | undefined;
}

export interface GroupSettings {
    [key: string]: string | undefined;
}

export interface ControllerPreset {
    id: string;
    name: string;
    group_source_bindings: Record<string, string>;
}

export interface ControllerMappingsSettings {
    left_trackpad_mode?: string;
    right_trackpad_mode?: string;
    action_set_trigger_cursor_show?: string;
    action_set_trigger_cursor_hide?: string;
    [key: string]: string | undefined;
}

