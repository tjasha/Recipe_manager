import axios from "axios";
const envBaseURL = import.meta.env.VITE_API_BASE_URL;
const baseURL =
    envBaseURL || (import.meta.env.DEV ? "http://localhost:8080/api" : "/api");
if (import.meta.env.PROD && !envBaseURL) {
    throw new Error("Missing VITE_API_BASE_URL in production environment");
}

const api = axios.create({
    baseURL: baseURL,
    withCredentials: true,
});

export default api;