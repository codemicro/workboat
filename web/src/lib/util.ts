import { onMount } from "svelte";

export const makeAPIURL = (path: string) => {
    const baseURL = BASE_API_URL.endsWith("/") ? BASE_API_URL.slice(0, -1) : BASE_API_URL;
    path = path.startsWith("/") ? path.slice(1) : path;
    return `${baseURL}/${path}`;
};

export const isLoggedIn = (): boolean => {
    return document.cookie.indexOf("workboat_session") !== -1;
}

export const doLogin = (nextPath: string = undefined) => {
    const requestPath = "/api/auth/newLogin" + (nextPath !== undefined ? "?next=" + nextPath : "");

    fetch(makeAPIURL(requestPath)).then((res) => {
        if (!res.ok) {
            return;
        }
        return res.json()
    }).then((jsonData) => {
        if (jsonData === undefined) {
            return;
        }

        window.location.replace(jsonData);
    });
}

export const checkLogin = () => {
    onMount(() => {
        if (!isLoggedIn()) {
            doLogin(window.location.pathname);
        }
    });
}