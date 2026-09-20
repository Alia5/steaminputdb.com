(() => {
    const appId = goTmpl('.AppID');
    const roots = [window.opener, window];

    try {
        const bpmActive = roots.some(
            (w) => typeof w?.SteamUIStore?.IsGamepadUIWindowActive === 'function'
                && w.SteamUIStore.IsGamepadUIWindowActive()
        );

        if (bpmActive) {
            const uiStore = roots.find(
                (w) => typeof w?.SteamUIStore?.Navigate === 'function')?.SteamUIStore;
            if (uiStore?.Navigate) {
                uiStore.Navigate(`/app/${appId}/controllerconfigurator/main`);
                return;
            }
        } else {
            const apps = roots.find(
                (w) => typeof w?.SteamClient?.Apps?.ShowControllerConfigurator === 'function'
            )
                ?.SteamClient?.Apps;
            if (apps?.ShowControllerConfigurator) {
                apps.ShowControllerConfigurator(appId);
                return;
            }
        }
    } catch (e) {
        console.error('SIDB: opening configurator failed, falling back', e);
    }
    // fallback
    SteamClient.Apps.ShowControllerConfigurator(appId);
})();
