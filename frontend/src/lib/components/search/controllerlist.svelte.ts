import EightBitDo from '$lib/assets/steam_controller_type_svgs/8bitdo_ultimate.svg?component';
import Hori from '$lib/assets/steam_controller_type_svgs/hori.svg?component';
import PS4 from '$lib/assets/steam_controller_type_svgs/ps4.svg?component';
import PS5 from '$lib/assets/steam_controller_type_svgs/ps5.svg?component';
import Gordon from '$lib/assets/steam_controller_type_svgs/steam.svg?component';
import SwitchPro from '$lib/assets/steam_controller_type_svgs/switchpro.svg?component';
import Triton from '$lib/assets/steam_controller_type_svgs/triton.svg?component';
import XBox from '$lib/assets/steam_controller_type_svgs/xbox.svg?component';
import IconXboxOne from '~icons/fluent/xbox-controller-24-filled';
import IconGameIconsSpartanHelmet from '~icons/game-icons/spartan-helmet';
import IconMdiCellPhone from '~icons/mdi/cellphone';
import IconMdiGamepad from '~icons/mdi/gamepad';
import IconSimpleIconsRepublicOfGamers from '~icons/simple-icons/republicofgamers';
import IconSD from '~icons/simple-icons/steamdeck';

export const CONTROLLER_LIST = [
    {
        type: 'controller_triton',
        icon: Triton,
        niceName: 'Steam Controller'
    },
    {
        type: 'controller_steamcontroller_gordon',
        icon: Gordon,
        niceName: 'Steam Controller (2015)'
    },
    {
        type: 'controller_neptune',
        icon: IconSD,
        niceName: 'Steam Deck'
    },
    {
        type: 'controller_ps5',
        icon: PS5,
        niceName: 'DualSense'
    },
    {
        type: 'controller_ps4',
        icon: PS4,
        niceName: 'DualShock 4'
    },
    {
        type: 'controller_xbox360',
        icon: XBox,
        niceName: 'Xbox 360'
    },
    {
        type: 'controller_xboxone',
        icon: IconXboxOne,
        niceName: 'Xbox One'
    },
    {
        type: 'controller_xboxelite',
        icon: IconXboxOne,
        niceName: 'Xbox Elite'
    },
    {
        type: 'controller_switch_pro',
        icon: SwitchPro,
        niceName: 'Switch Pro'
    },
    {
        type: 'controller_switch2_pro',
        icon: SwitchPro,
        niceName: 'Switch 2 Pro'
    },
    {
        type: 'controller_8bitdo',
        icon: EightBitDo,
        niceName: '8BitDo'
    },
    {
        type: 'controller_generic',
        icon: IconMdiGamepad,
        niceName: 'Generic'
    },
    {
        type: 'controller_steamcontroller_headcrab',
        icon: Gordon,
        niceName: 'Steam Controller (Headcrab)'
    },
    {
        type: 'controller_ps5_edge',
        icon: PS5,
        niceName: 'DualSense Edge'
    },
    {
        type: 'controller_ps3',
        icon: IconMdiGamepad,
        niceName: 'DualShock 3'
    },
    {
        type: 'controller_hori_steam',
        icon: Hori,
        niceName: 'HoriPad Steam'
    },
    {
        type: 'controller_mobile_touch',
        icon: IconMdiCellPhone,
        niceName: 'Mobile Touch'
    },
    {
        type: 'controller_rog_ally',
        icon: IconSimpleIconsRepublicOfGamers,
        niceName: 'ASUS ROG Ally'
    },
    {
        type: 'controller_legion_go_s',
        icon: IconGameIconsSpartanHelmet,
        niceName: 'Lenovo Legion Go S'
    }
] as const;
