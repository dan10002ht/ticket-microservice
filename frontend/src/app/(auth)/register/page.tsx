"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { User, Building2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PasswordInput } from "@/components/ui/password-input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";
import { showToast } from "@/lib/toast";
import { registerSchema, type RegisterInput } from "@/lib/validators/auth";
import { useRegister } from "@/lib/api/queries";
import type { ApiError } from "@/lib/api/types/common";

// Map API field paths to react-hook-form field names
const API_FIELD_MAP: Record<string, keyof RegisterInput> = {
  email: "email",
  password: "password",
  first_name: "first_name",
  last_name: "last_name",
};

export default function RegisterPage() {
  const router = useRouter();
  const registerMutation = useRegister();

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    setError,
    formState: { errors },
  } = useForm<RegisterInput>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      first_name: "",
      last_name: "",
      email: "",
      password: "",
      confirmPassword: "",
      role: "individual",
      organization_name: "",
    },
  });

  const password = watch("password");
  const selectedRole = watch("role");

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

  const strength = getPasswordStrength(password);

  const onSubmit = async (data: RegisterInput) => {
    try {
      const { confirmPassword: _, organization_name, ...rest } = data;
      const registerData = {
        ...rest,
        ...(data.role === "organization" && organization_name
          ? { organization: { name: organization_name.trim() } }
          : {}),
      };
      await registerMutation.mutateAsync(registerData);
      showToast.success("Tạo tài khoản thành công!");
      router.push(data.role === "organization" ? "/org/dashboard" : "/");
    } catch (err: unknown) {
      // Map API validation details to form field errors
      const apiErr = err as ApiError;
      if (apiErr.error?.details && Array.isArray(apiErr.error.details)) {
        for (const detail of apiErr.error.details) {
          const d = detail as { path: string; msg: string };
          const formField = API_FIELD_MAP[d.path];
          if (formField) {
            setError(formField, { message: d.msg });
          }
        }
      }
      showToast.error(apiErr?.error?.message || "Registration failed. Please try again.");
    }
  };

  return (
    <div className="space-y-6">
      <div className="text-center">
        <h1 className="text-2xl font-bold">Create account</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Get started with TicketBox
        </p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {/* Role selection */}
        <div className="space-y-2">
          <Label>Account type</Label>
          <div className="grid grid-cols-2 gap-3">
            <button
              type="button"
              onClick={() => setValue("role", "individual", { shouldValidate: true })}
              className={cn(
                "flex flex-col items-center gap-2 rounded-lg border p-4 transition-all",
                selectedRole === "individual"
                  ? "border-primary bg-primary/5 ring-1 ring-primary"
                  : "hover:border-muted-foreground/50"
              )}
            >
              <User className="h-5 w-5" />
              <span className="text-sm font-medium">Individual</span>
              <span className="text-xs text-muted-foreground">
                Browse & book events
              </span>
            </button>
            <button
              type="button"
              onClick={() => setValue("role", "organization", { shouldValidate: true })}
              className={cn(
                "flex flex-col items-center gap-2 rounded-lg border p-4 transition-all",
                selectedRole === "organization"
                  ? "border-primary bg-primary/5 ring-1 ring-primary"
                  : "hover:border-muted-foreground/50"
              )}
            >
              <Building2 className="h-5 w-5" />
              <span className="text-sm font-medium">Organization</span>
              <span className="text-xs text-muted-foreground">
                Create & manage events
              </span>
            </button>
          </div>
        </div>

        {/* Organization name (visible only for org role) */}
        {selectedRole === "organization" && (
          <div className="space-y-2">
            <Label htmlFor="organization_name">Organization name</Label>
            <Input
              id="organization_name"
              placeholder="Your company or team name"
              {...register("organization_name")}
            />
            {errors.organization_name && (
              <p className="text-xs text-destructive">
                {errors.organization_name.message}
              </p>
            )}
          </div>
        )}

        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="first_name">First name</Label>
            <Input
              id="first_name"
              placeholder="John"
              autoComplete="given-name"
              {...register("first_name")}
            />
            {errors.first_name && (
              <p className="text-xs text-destructive">
                {errors.first_name.message}
              </p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="last_name">Last name</Label>
            <Input
              id="last_name"
              placeholder="Doe"
              autoComplete="family-name"
              {...register("last_name")}
            />
            {errors.last_name && (
              <p className="text-xs text-destructive">
                {errors.last_name.message}
              </p>
            )}
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="you@example.com"
            autoComplete="email"
            {...register("email")}
          />
          {errors.email && (
            <p className="text-xs text-destructive">{errors.email.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="password">Password</Label>
          <PasswordInput
            id="password"
            placeholder="Min 8 characters"
            autoComplete="new-password"
            {...register("password")}
          />
          {errors.password && (
            <p className="text-xs text-destructive">
              {errors.password.message}
            </p>
          )}
          {password && !errors.password && (
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
            placeholder="Confirm your password"
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
          disabled={registerMutation.isPending}
        >
          {registerMutation.isPending ? "Creating account..." : "Create account"}
        </Button>
      </form>

      <Separator />

      <p className="text-center text-sm text-muted-foreground">
        Already have an account?{" "}
        <Link
          href="/login"
          className="font-medium text-primary hover:underline"
        >
          Sign in
        </Link>
      </p>
    </div>
  );
}
