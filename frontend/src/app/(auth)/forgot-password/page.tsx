"use client";

import { useState, useEffect, useCallback } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, Mail } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  forgotPasswordSchema,
  type ForgotPasswordInput,
} from "@/lib/validators/auth";
import { useForgotPassword } from "@/lib/api/queries";
import { showToast } from "@/lib/toast";
import type { ApiError } from "@/lib/api/types/common";

const RESEND_COOLDOWN = 60;

export default function ForgotPasswordPage() {
  const [submitted, setSubmitted] = useState(false);
  const [cooldown, setCooldown] = useState(0);
  const forgotMutation = useForgotPassword();

  const {
    register,
    handleSubmit,
    getValues,
    formState: { errors },
  } = useForm<ForgotPasswordInput>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: { email: "" },
  });

  useEffect(() => {
    if (cooldown <= 0) return;
    const timer = setInterval(() => setCooldown((c) => c - 1), 1000);
    return () => clearInterval(timer);
  }, [cooldown]);

  const onSubmit = async (data: ForgotPasswordInput) => {
    try {
      await forgotMutation.mutateAsync(data);
      setSubmitted(true);
      setCooldown(RESEND_COOLDOWN);
    } catch (err: unknown) {
      const apiError = err as ApiError;
      showToast.error(
        apiError?.error?.message ||
          "Failed to send reset email. Please try again."
      );
    }
  };

  const handleResend = useCallback(async () => {
    const email = getValues("email");
    if (!email || cooldown > 0) return;
    try {
      await forgotMutation.mutateAsync({ email });
      setCooldown(RESEND_COOLDOWN);
      showToast.success("Reset link sent again!");
    } catch (err: unknown) {
      const apiError = err as ApiError;
      showToast.error(
        apiError?.error?.message || "Failed to resend. Please try again."
      );
    }
  }, [cooldown, forgotMutation, getValues]);

  if (submitted) {
    return (
      <div className="space-y-6">
        <div className="flex flex-col items-center text-center">
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
            <Mail className="h-6 w-6 text-primary" />
          </div>
          <h1 className="mt-4 text-2xl font-bold">Check your email</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            We sent a password reset link to{" "}
            <span className="font-medium text-foreground">
              {getValues("email")}
            </span>
          </p>
        </div>

        <div className="space-y-3">
          <p className="text-center text-xs text-muted-foreground">
            Didn&apos;t receive the email? Check your spam folder or
          </p>
          <Button
            variant="outline"
            className="w-full"
            disabled={cooldown > 0 || forgotMutation.isPending}
            onClick={handleResend}
          >
            {forgotMutation.isPending
              ? "Sending..."
              : cooldown > 0
                ? `Resend in ${cooldown}s`
                : "Resend reset link"}
          </Button>
          <Button
            variant="ghost"
            className="w-full"
            onClick={() => {
              setSubmitted(false);
              forgotMutation.reset();
            }}
          >
            Try another email
          </Button>
        </div>

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

  return (
    <div className="space-y-6">
      <div className="text-center">
        <h1 className="text-2xl font-bold">Forgot password?</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          No worries, we&apos;ll send you reset instructions
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="you@example.com"
            autoComplete="email"
            autoFocus
            {...register("email")}
          />
          {errors.email && (
            <p className="text-xs text-destructive">{errors.email.message}</p>
          )}
        </div>

        <Button
          type="submit"
          className="w-full"
          disabled={forgotMutation.isPending}
        >
          {forgotMutation.isPending ? "Sending..." : "Send reset link"}
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
