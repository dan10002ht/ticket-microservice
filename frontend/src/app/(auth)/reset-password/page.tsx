"use client";

import { useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, CheckCircle2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { PasswordInput } from "@/components/ui/password-input";
import { Label } from "@/components/ui/label";
import {
  resetPasswordSchema,
  type ResetPasswordInput,
} from "@/lib/validators/auth";
import { useResetPassword } from "@/lib/api/queries";
import { showToast } from "@/lib/toast";
import type { ApiError } from "@/lib/api/types/common";

export default function ResetPasswordPage() {
  const searchParams = useSearchParams();
  const token = searchParams.get("token") || "";
  const [success, setSuccess] = useState(false);
  const resetMutation = useResetPassword();

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<ResetPasswordInput>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      token,
      password: "",
      confirmPassword: "",
    },
  });

  const newPassword = watch("password");

  const getPasswordStrength = (
    pwd: string
  ): { label: string; width: string; color: string } => {
    if (!pwd) return { label: "", width: "0%", color: "" };
    let score = 0;
    if (pwd.length >= 8) score++;
    if (/[A-Z]/.test(pwd)) score++;
    if (/[0-9]/.test(pwd)) score++;
    if (/[^A-Za-z0-9]/.test(pwd)) score++;
    const map = [
      { label: "Weak", width: "25%", color: "bg-destructive" },
      { label: "Fair", width: "50%", color: "bg-yellow-500" },
      { label: "Good", width: "75%", color: "bg-blue-500" },
      { label: "Strong", width: "100%", color: "bg-green-500" },
    ];
    return map[Math.min(score, 4) - 1] ?? map[0];
  };

  const strength = getPasswordStrength(newPassword);

  const onSubmit = async (data: ResetPasswordInput) => {
    try {
      const { confirmPassword: _, ...payload } = data;
      await resetMutation.mutateAsync(payload);
      setSuccess(true);
    } catch (err: unknown) {
      const apiError = err as ApiError;
      showToast.error(
        apiError?.error?.message ||
          "Failed to reset password. The link may have expired."
      );
    }
  };

  if (!token) {
    return (
      <div className="space-y-6">
        <div className="flex flex-col items-center text-center">
          <h1 className="text-2xl font-bold">Invalid reset link</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            This password reset link is invalid or has expired. Please request a
            new one.
          </p>
        </div>

        <Button asChild className="w-full">
          <Link href="/forgot-password">Request new link</Link>
        </Button>

        <Link
          href="/login"
          className="flex items-center justify-center gap-2 text-sm text-muted-foreground hover:text-primary"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to sign in
        </Link>
      </div>
    );
  }

  if (success) {
    return (
      <div className="space-y-6">
        <div className="flex flex-col items-center text-center">
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-green-500/10">
            <CheckCircle2 className="h-6 w-6 text-green-500" />
          </div>
          <h1 className="mt-4 text-2xl font-bold">Password reset</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Your password has been successfully reset. You can now sign in with
            your new password.
          </p>
        </div>

        <Button asChild className="w-full">
          <Link href="/login">Sign in</Link>
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="text-center">
        <h1 className="text-2xl font-bold">Set new password</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Your new password must be different from previously used passwords
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <input type="hidden" {...register("token")} />

        <div className="space-y-2">
          <Label htmlFor="password">New password</Label>
          <PasswordInput
            id="password"
            placeholder="Min 8 characters"
            autoComplete="new-password"
            autoFocus
            {...register("password")}
          />
          {errors.password && (
            <p className="text-xs text-destructive">
              {errors.password.message}
            </p>
          )}
          {newPassword && !errors.password && (
            <div className="space-y-1">
              <div className="h-1.5 w-full overflow-hidden rounded-full bg-muted">
                <div
                  className={`h-full rounded-full transition-all ${strength.color}`}
                  style={{ width: strength.width }}
                />
              </div>
              <p className="text-xs text-muted-foreground">{strength.label}</p>
            </div>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="confirmPassword">Confirm password</Label>
          <PasswordInput
            id="confirmPassword"
            placeholder="Confirm your new password"
            autoComplete="new-password"
            {...register("confirmPassword")}
          />
          {errors.confirmPassword && (
            <p className="text-xs text-destructive">
              {errors.confirmPassword.message}
            </p>
          )}
        </div>

        <Button
          type="submit"
          className="w-full"
          disabled={resetMutation.isPending}
        >
          {resetMutation.isPending ? "Resetting..." : "Reset password"}
        </Button>
      </form>

      <Link
        href="/login"
        className="flex items-center justify-center gap-2 text-sm text-muted-foreground hover:text-primary"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to sign in
      </Link>
    </div>
  );
}
