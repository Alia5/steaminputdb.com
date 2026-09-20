
__INJECT_RETURN = (async () => {
    await SteamClient.Input.GetConfigForAppAndController(
        goTmpl('.AppID'),
        goTmpl('.ControllerIndex')
    );
    await SteamClient.Apps.DownloadWorkshopItem(241100, goTmpl('.WorkshopItemID'), true);
    try {
        await SteamClient.Input.SetSelectedConfigForApp(
            goTmpl('.AppID'),
            goTmpl('.ControllerIndex'),
            'workshop://' + goTmpl('.WorkshopItemID'),
            false,
            0
        );
    } catch (error) {
        console.error(error);
        await SteamClient.Input.SetSelectedConfigForApp(
            goTmpl('.AppID'),
            goTmpl('.ControllerIndex'),
            'workshop://' + goTmpl('.WorkshopItemID'),
            false,
            true
        ).catch((e) => e);
    }
})();
