// API Configuration and Utilities

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: `${API_BASE_URL}/api/v1/auth/login`,
    REGISTER: `${API_BASE_URL}/api/v1/auth/register`,
  },
  DOMAINS: {
    LIST: `${API_BASE_URL}/api/v1/domains`,
    CREATE: `${API_BASE_URL}/api/v1/domains`,
    VERIFY: `${API_BASE_URL}/api/v1/domains/verify`,
  },
  SCAN: `${API_BASE_URL}/api/v1/scan`,
  STATS: `${API_BASE_URL}/api/v1/stats`,
  ASSETS: `${API_BASE_URL}/api/v1/assets`,
  SERVICES: `${API_BASE_URL}/api/v1/services`,
  FINDINGS: `${API_BASE_URL}/api/v1/findings`,
};

export interface APIError {
  code: string;
  message: string;
  details?: any;
}

export function getAuthHeaders(): HeadersInit {
  const token = localStorage.getItem('token');
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
}

export async function handleAPIResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let error: APIError;
    try {
      error = await response.json();
    } catch {
      error = {
        code: 'UNKNOWN_ERROR',
        message: `HTTP ${response.status}: ${response.statusText}`,
      };
    }
    throw new Error(error.message || error.code || 'Request failed');
  }
  return response.json();
}
