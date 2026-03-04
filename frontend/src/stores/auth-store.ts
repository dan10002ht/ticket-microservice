import { create } from "zustand";
import type { AuthUser } from "@/lib/api/types/auth";

interface AuthState {
  user: AuthUser | null;
  isAuthenticated: boolean;
  /** true until initial auth check (useMe) completes */
  isLoading: boolean;
  /** prevents SSR hydration mismatch */
  isHydrated: boolean;

  setUser: (user: AuthUser) => void;
  clearUser: () => void;
  setLoading: (loading: boolean) => void;
  setHydrated: (hydrated: boolean) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: true,
  isHydrated: false,

  setUser: (user) =>
    set({ user, isAuthenticated: true, isLoading: false }),
  clearUser: () =>
    set({ user: null, isAuthenticated: false, isLoading: false }),
  setLoading: (isLoading) => set({ isLoading }),
  setHydrated: (isHydrated) => set({ isHydrated }),
}));
