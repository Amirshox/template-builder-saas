import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface Template {
  id: string;
  org_id: string;
  name: string;
  type: 'layout' | 'docx';
  status: 'active' | 'archived';
  created_at: string;
}

export const fetchTemplates = async (orgId?: string): Promise<Template[]> => {
  const response = await api.get('/templates', {
    params: { orgId },
  });
  return response.data;
};

export const getTemplate = async (id: string): Promise<Template> => {
  const response = await api.get(`/templates/${id}`);
  return response.data;
};

export interface TemplateVersion {
  id: string;
  version: number;
  status: string;
  created_at: string;
  template_json: any;
}

export const getVersions = async (templateId: string): Promise<TemplateVersion[]> => {
  const response = await api.get(`/templates/${templateId}/versions`);
  return response.data;
};

export const createVersion = async (id: string, data: any) => {
  const payload = {
    schemaJson: {} // TODO: infer schema
  }
  const response = await api.post(`/templates/${id}/versions`, payload)
  return response.data
}

export const createTemplate = async (
  orgId: string,
  name: string,
  type: 'layout' | 'docx'
): Promise<Template> => {
  const response = await api.post('/templates', { orgId, name, type });
  return response.data;
};


export const downloadPreview = async (templateId: string, version: number = 1): Promise<void> => {
  const response = await api.post(`/templates/${templateId}/preview`, {}, {
    responseType: 'blob', // Important for PDF
  });
  // Create blob link to download
  const url = window.URL.createObjectURL(new Blob([response.data]));
  const link = document.createElement('a');
  link.href = url;
  link.setAttribute('download', 'preview.pdf');
  document.body.appendChild(link);
  link.click();
  link.remove();
};

// Auth Types
export interface AuthResponse {
  token: string;
}

// Add interceptor for Auth
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token');
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export const register = async (email: string, password: string, name: string): Promise<AuthResponse> => {
  const response = await api.post('/register', { email, password, name });
  return response.data;
};

export const login = async (email: string, password: string): Promise<AuthResponse> => {
  const response = await api.post('/login', { email, password });
  return response.data;
};

export const logout = () => {
  localStorage.removeItem('token');
  window.location.href = '/login';
};

export default api;
