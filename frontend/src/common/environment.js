const development = {
    frontend: "http://localhost:3011",
    backend: "http://localhost:5011/api"
};

const production = {
    frontend: "https://ldap-auth.chalmers.it",
    backend: "https://ldap-auth.chalmers.it/api"
};

function getFrontendUrl() {
    return isDevelopment() ? development.frontend : production.frontend;
}

function getBackendUrl() {
    return isDevelopment() ? development.backend : production.backend;
}

function isDevelopment() {
    return process.env.NODE_ENV === "development";
}

export { getBackendUrl, getFrontendUrl };
