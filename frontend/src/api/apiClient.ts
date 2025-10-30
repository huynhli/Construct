const BASE_URL = import.meta.env.VITE_API_BASE_URL;

async function fetchWithAuth(endpoint: string, options: RequestInit = {}) {
    const token = localStorage.getItem('authToken');

    const headers: HeadersInit = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string>),
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const config: RequestInit = {
        ...options,
        headers,
    };

    const response = await fetch(`${BASE_URL}${endpoint}`, config);

    if (!response.ok) {
        const errorText = await response.text();
        try {
            const errorData = JSON.parse(errorText);
            throw new Error(errorData.message || 'API request failed');
        } catch (e) {
            throw new Error(errorText || 'An unknown network error occurred');
        }
    }

    const responseText = await response.text();
    return responseText ? JSON.parse(responseText) : {};
}

const apiClient = {
    get: (endpoint: string) => fetchWithAuth(endpoint),
    post: (endpoint: string, body: any) => fetchWithAuth(endpoint, { method: 'POST', body: JSON.stringify(body) }),
    put: (endpoint: string, body: any) => fetchWithAuth(endpoint, { method: 'PUT', body: JSON.stringify(body) }),
    delete: (endpoint: string) => fetchWithAuth(endpoint, { method: 'DELETE' }),
};

export default apiClient;