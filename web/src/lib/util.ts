export const isLoggedIn = (): boolean => {
    return (document.cookie
        .split('; ')
        .find((row) => row.startsWith('workboat_session'))) !== undefined;
};