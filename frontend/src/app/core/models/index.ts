export interface User {
  id: number;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
  createdAt: string;
  updatedAt: string;
}

export interface Client {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  address: string;
  notes: string;
  userId: number;
  createdAt: string;
  updatedAt: string;
  horses?: Horse[];
}

export interface Horse {
  id: number;
  name: string;
  breed: string;
  age: number;
  gender: string;
  color: string;
  weight: number;
  notes: string;
  clientId: number;
  barnId: number;
  client?: Client;
  barn?: Barn;
  createdAt: string;
  updatedAt: string;
}

export interface Barn {
  id: number;
  name: string;
  address: string;
  phone: string;
  email: string;
  notes: string;
  userId: number;
  createdAt: string;
  updatedAt: string;
  horses?: Horse[];
}

export interface Appointment {
  id: number;
  clientId: number;
  horseId: number;
  barnId: number;
  userId: number;
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
  id: number;
  appointmentId: number;
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
  id: number;
  invoiceId: number;
  description: string;
  quantity: number;
  unitPrice: number;
  amount: number;
}

export interface Invoice {
  id: number;
  clientId: number;
  userId: number;
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
  firstName: string;
  lastName: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}
