import axios from "axios";
import { attachAuthToken, createRefreshInterceptor } from "./interceptors/auth";
import { handleResponseError, handleResponseSuccess } from "./interceptors/error";

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:53000/api/v1",
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request: attach JWT token
apiClient.interceptors.request.use(attachAuthToken, (error) =>
  Promise.reject(error)
);

// Response: auto-refresh on 401 (must register BEFORE error handler)
apiClient.interceptors.response.use(
  (response) => response,
  createRefreshInterceptor(apiClient)
);

// Response: normalize errors + auto-toast
apiClient.interceptors.response.use(handleResponseSuccess, handleResponseError);
