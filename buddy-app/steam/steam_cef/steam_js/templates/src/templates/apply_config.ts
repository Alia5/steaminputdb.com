__INJECT_RETURN = (async () => {
    await SteamClient.Input.GetConfigForAppAndController(
        goTmpl('.AppID'),
        goTmpl('.ControllerIndex')
    );
    await SteamClient.Apps.DownloadWorkshopItem(241100, goTmpl('.WorkshopItemID'), true);
    await SteamClient.Input.SetSelectedConfigForApp(
        goTmpl('.AppID'),
        goTmpl('.ControllerIndex'),
        'workshop://'+goTmpl('.WorkshopItemID'),
        false,
        true
    );
})();
