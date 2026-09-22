export const getModUrlName = (url: string) => {
    if (url.toLowerCase().includes('steamcommunity.com')) {
        return 'Steam Workshop';
    }
    if (url.toLowerCase().includes('nexusmods.com')) {
        return 'Nexus Mods';
    }
    if (url.toLowerCase().includes('github.com')) {
        return 'GitHub';
    }
    if (url.toLowerCase().includes('gitlab')) {
        return 'GitLab';
    }
    if (url.toLowerCase().includes('moddb')) {
        return 'Mod DB';
    }
    if (url.toLowerCase().includes('codeberg')) {
        return 'Codeberg';
    }
    return url.replaceAll(/http(s)?:\/\/(www\.)?/g, '');
};


export const getModHostNameNice = (url: string) => {
    let name = getModUrlName(url);
    name = name.split('/')?.[0] || name;
    try {
        const host = new URL(name).hostname;
        return host.replaceAll(/http(s)?:\/\/(www\.)?/g, '');
    } catch {
        return name;
    }
};
