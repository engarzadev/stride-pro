export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
  createdAt: string;
  updatedAt: string;
}

export interface Client {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  address: string;
  notes: string;
  userId: string;
  createdAt: string;
  updatedAt: string;
  horses?: Horse[];
}

export interface Horse {
  id: string;
  name: string;
  breed: string;
  age: number;
  gender: string;
  color: string;
  weight: number;
  notes: string;
  clientId: string;
  barnId: string | null;
  client?: Client;
  barn?: Barn;
  createdAt: string;
  updatedAt: string;
}

export interface Barn {
  id: string;
  name: string;
  contactName: string;
  address: string;
  phone: string;
  email: string;
  notes: string;
  userId: string;
  createdAt: string;
  updatedAt: string;
  horses?: Horse[];
}

export interface Appointment {
  id: string;
  clientId: string;
  horseId: string;
  barnId: string | null;
  userId: string;
  date: string;
  time: string;
  duration: number;
  type: string;
  status: string;
  notes: string;
  client?: Client;
  horse?: Horse;
  barn?: Barn;
  session?: Session;
  createdAt: string;
  updatedAt: string;
}

export interface Session {
  id: string;
  appointmentId: string;
  type: string;
  bodyZones: string[];
  notes: string;
  findings: string;
  recommendations: string;
  appointment?: Appointment;
  createdAt: string;
  updatedAt: string;
}

export interface InvoiceItem {
  id: string;
  invoiceId: string;
  description: string;
  quantity: number;
  unitPrice: number;
  amount: number;
}

export interface Invoice {
  id: string;
  clientId: string;
  userId: string;
  invoiceNumber: string;
  date: string;
  dueDate: string;
  status: string;
  subtotal: number;
  tax: number;
  total: number;
  notes: string;
  client?: Client;
  items?: InvoiceItem[];
  createdAt: string;
  updatedAt: string;
}

export interface ApiResponse<T> {
  data: T;
  error?: string;
  meta?: Record<string, unknown>;
}

export interface PaginationParams {
  page: number;
  limit: number;
  sort?: string;
  order?: 'asc' | 'desc';
  search?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    total: number;
    page: number;
    limit: number;
    totalPages: number;
  };
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
}

export interface AuthResponse {
  tokens: AuthTokens;
  user: User;
}
