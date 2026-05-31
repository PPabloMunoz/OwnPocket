import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router";
import { useCallback } from "react";
import { api } from "@/lib/api";
import { useAuthStore } from "@/stores/auth-store";
import type { AuthResponse, LoginRequest, RegisterRequest } from "@/types/auth";

export function useLogin() {
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: (data: LoginRequest) => api.post<AuthResponse>("/auth/login", data),
    onSuccess: (result) => {
      setAuth(result.token, result.user);
    },
  });
}

export function useRegister() {
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: (data: RegisterRequest) => api.post<AuthResponse>("/auth/register", data),
    onSuccess: (result) => {
      setAuth(result.token, result.user);
    },
  });
}

export function useLogout() {
  const logout = useAuthStore((s) => s.logout);
  const navigate = useNavigate();

  return useCallback(() => {
    logout();
    navigate("/login");
  }, [logout, navigate]);
}
