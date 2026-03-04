// Note: User service uses camelCase (unlike other services that use snake_case)

// ── Domain types ──

export interface UserProfile {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  phone?: string;
  dateOfBirth?: string;
  avatarUrl?: string;
  preferences?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface UserAddress {
  id: string;
  street: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
  isDefault?: boolean;
}

// ── Request types ──

export interface UserProfileUpdateRequest {
  firstName?: string;
  lastName?: string;
  phone?: string;
  dateOfBirth?: string;
  avatarUrl?: string;
  preferences?: Record<string, unknown>;
}

export interface UserProfileCreateRequest {
  first_name: string;
  last_name: string;
  email: string;
  phone?: string;
  avatar_url?: string;
  date_of_birth?: string;
  preferences?: Record<string, unknown>;
}

export interface UserAddressCreateRequest {
  street: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
  isDefault?: boolean;
}
