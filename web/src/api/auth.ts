import axios from 'axios';

const api = axios.create({
    baseURL: '/api',
    withCredentials: true, // Important! Sends cookies with requests
    headers: {
        'Content-Type': 'application/json',
    },
});

// Get CSRF token from cookie
function getCSRFToken(): string | null {
    const match = document.cookie.match(/csrf_token=([^;]+)/);
    return match ? match[1] : null;
}

// Add CSRF token to requests
api.interceptors.request.use((config) => {
    const csrfToken = getCSRFToken();
    if (csrfToken && config.method !== 'get') {
        config.headers['X-CSRF-Token'] = csrfToken;
        config.headers['X-Requested-With'] = 'XMLHttpRequest';
    }
    return config;
});

// Handle 401 errors (token expired)
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        // If 401 and haven't retried yet, try refreshing token
        // BUT skip if the failed request was already a refresh attempt (avoids infinite loop)
        if (
            error.response?.status === 401 && 
            !originalRequest._retry && 
            !originalRequest.url?.includes('/auth/refresh')
        ) {
            originalRequest._retry = true;

            try {
                await authAPI.refresh();
                return api(originalRequest);
            } catch (refreshError) {
                // Refresh failed, redirect to login
                window.location.href = '/login';
                return Promise.reject(refreshError);
            }
        }

        return Promise.reject(error);
    }
);

export interface RegisterData {
    email: string;
    password: string;
}

export interface LoginData {
    email: string;
    password: string;
}

export interface User {
    id: number;
    email: string;
    verified: boolean;
    created_at: string;
}

export const authAPI = {
    register: (data: RegisterData) =>
        api.post('/auth/register', data),

    login: (data: LoginData) =>
        api.post('/auth/login', data),

    logout: () =>
        api.post('/auth/logout'),

    refresh: () =>
        api.post('/auth/refresh'),

    getMe: () =>
        api.get<{ user: User }>('/user/me'),
};

export default api;
