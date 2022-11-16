export const makeAPIURL = (path: string) => {
    const baseURL = BASE_API_URL.endsWith("/") ? BASE_API_URL.slice(0, -1) : BASE_API_URL;
    path = path.startsWith("/") ? path.slice(1) : path;
    return `${baseURL}/${path}`;
};