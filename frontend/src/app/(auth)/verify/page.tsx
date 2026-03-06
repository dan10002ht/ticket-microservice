"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { MailCheck } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { showToast } from "@/lib/toast";
import { useAuthStore } from "@/stores";
import { useVerifyUser, useSendVerificationEmail } from "@/lib/api/queries";
import type { ApiError } from "@/lib/api/types/common";

const RESEND_COOLDOWN = 60;

export default function VerifyPage() {
  const router = useRouter();
  const { user, setUser } = useAuthStore();
  const verifyMutation = useVerifyUser();
  const resendMutation = useSendVerificationEmail();

  const [pin, setPin] = useState("");
  const [cooldown, setCooldown] = useState(RESEND_COOLDOWN);

  // Start cooldown timer on mount (email was just sent after register)
  useEffect(() => {
    if (cooldown <= 0) return;
    const timer = setInterval(() => setCooldown((c) => c - 1), 1000);
    return () => clearInterval(timer);
  }, [cooldown]);

  // Redirect if already verified or not logged in
  useEffect(() => {
    if (!user) {
      router.replace("/login");
    } else if (user.isVerified) {
      router.replace("/");
    }
  }, [user, router]);

  if (!user || user.isVerified) return null;

  const handleVerify = async (value: string) => {
    if (value.length !== 6) return;
    try {
      const result = await verifyMutation.mutateAsync({
        user_id: user.id,
        pin_code: value,
      });
      if (result.user) {
        setUser({ ...user, isVerified: true });
      }
      showToast.success("Email verified successfully!");
      router.push("/");
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      const code = apiErr?.error?.code;
      if (code === "PIN_CODE_EXPIRED") {
        showToast.error("PIN expired. Please request a new one.");
      } else if (code === "INVALID_PIN_CODE") {
        showToast.error("Invalid PIN. Please check and try again.");
      } else {
        showToast.error(apiErr?.error?.message || "Verification failed.");
      }
      setPin("");
    }
  };

  const handleResend = async () => {
    try {
      await resendMutation.mutateAsync({ email: user.email });
      showToast.success("Verification email sent!");
      setCooldown(RESEND_COOLDOWN);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      if (apiErr?.error?.code === "EMAIL_ALREADY_VERIFIED") {
        setUser({ ...user, isVerified: true });
        router.push("/");
        return;
      }
      showToast.error(apiErr?.error?.message || "Failed to resend email.");
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col items-center gap-3 text-center">
        <div className="flex h-14 w-14 items-center justify-center rounded-full bg-primary/10">
          <MailCheck className="h-7 w-7 text-primary" />
        </div>
        <h1 className="text-2xl font-bold">Verify your email</h1>
        <p className="text-sm text-muted-foreground">
          We sent a 6-digit code to{" "}
          <span className="font-medium text-foreground">{user.email}</span>
        </p>
      </div>

      <div className="flex flex-col items-center gap-4">
        <InputOTP
          maxLength={6}
          value={pin}
          onChange={setPin}
          onComplete={handleVerify}
          disabled={verifyMutation.isPending}
        >
          <InputOTPGroup>
            <InputOTPSlot index={0} />
            <InputOTPSlot index={1} />
            <InputOTPSlot index={2} />
            <InputOTPSlot index={3} />
            <InputOTPSlot index={4} />
            <InputOTPSlot index={5} />
          </InputOTPGroup>
        </InputOTP>

        <Button
          className="w-full"
          onClick={() => handleVerify(pin)}
          disabled={pin.length !== 6 || verifyMutation.isPending}
        >
          {verifyMutation.isPending ? "Verifying..." : "Verify"}
        </Button>
      </div>

      <div className="text-center">
        <p className="text-sm text-muted-foreground">
          Didn&apos;t receive the email?{" "}
          {cooldown > 0 ? (
            <span className="text-muted-foreground">
              Resend in {cooldown}s
            </span>
          ) : (
            <button
              type="button"
              onClick={handleResend}
              disabled={resendMutation.isPending}
              className="font-medium text-primary hover:underline disabled:opacity-50"
            >
              {resendMutation.isPending ? "Sending..." : "Resend code"}
            </button>
          )}
        </p>
        <p className="mt-2 text-xs text-muted-foreground">
          Code expires in 15 minutes
        </p>
      </div>
    </div>
  );
}
