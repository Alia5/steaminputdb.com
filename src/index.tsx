import {
    PanelSection,
    PanelSectionRow,
    ToggleField,
    staticClasses
} from "@decky/ui";
import {
    definePlugin,
    toaster,
} from "@decky/api"
import { useState, useEffect } from "react";
import { Logo } from "./Logo";

const BUDDY_URL = "http://localhost:5119";

interface BuddySettings {
    addDesktopUIEntries?: boolean;
    addBigPictureUIEntries?: boolean;
}

const errToast = (msg: string) => {
    toaster.toast({
        title: "SteamInputDB-Buddy",
        body: msg,
        critical: true,
        eType: 2,
    });
}

const getSettings = async (): Promise<BuddySettings> => {
    const res = await fetch(`${BUDDY_URL}/v1/settings`);
    if (!res.ok) {
        throw new Error(`GET /v1/settings failed: ${res.status}`);
    }
    return res.json();
};

const putSettings = async (settings: BuddySettings): Promise<BuddySettings> => {
    const res = await fetch(`${BUDDY_URL}/v1/settings`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(settings),
    });
    if (!res.ok) {
        throw new Error(`PUT /v1/settings failed: ${res.status}`);
    }
    return res.json();
};

function Content() {
    const [enabled, setEnabled] = useState(false);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        getSettings()
            .then((s) => {
                setEnabled(s.addBigPictureUIEntries === true && s.addDesktopUIEntries === true);
            })
            .catch((e) => {
                console.error("SteamInputDB-Buddy: failed to get settings", e);
                errToast("Failed to get settings from SteamInputDB-Buddy");
            })
            .finally(() => {
                setLoading(false);
            });
    }, []);

    const onToggle = (checked: boolean) => {
        setEnabled(checked);
        putSettings({
            addDesktopUIEntries: checked,
            addBigPictureUIEntries: checked,
        }).catch((e) => {
            console.error("SteamInputDB-Buddy: failed to update settings", e);
            errToast("Failed to update SteamInputDB-Buddy settings");
            setEnabled(!checked);
        });
    };

    return (
        <PanelSection>
            <PanelSectionRow>
                <ToggleField
                    label="SteamInputDB Browse buttons"
                    description="Show SteamInputDB buttons in Steams UI"
                    checked={enabled}
                    disabled={loading}
                    onChange={onToggle}
                />
            </PanelSectionRow>
        </PanelSection>
    );
};


export default definePlugin(() => {
    console.log("SteamInputDB-Buddy Decky plugin initializing...")

    return {
        name: "SteamInputDB-Buddy",
        titleView: <div className={staticClasses.Title}>SteamInputDB-Buddy</div>,
        content: <Content />,
        icon: <Logo />,
        onDismount() {
            console.log("Unloading")
        },
    };
});
