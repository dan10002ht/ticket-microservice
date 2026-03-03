import axios from "axios";

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3000/api",
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  // TODO: Attach auth token from store
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // TODO: Handle 401 → refresh token or redirect to login
    return Promise.reject(error);
  }
);
