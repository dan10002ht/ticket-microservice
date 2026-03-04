"use client";

import { Check, Ticket, ClipboardList, CreditCard, PartyPopper } from "lucide-react";
import type { BookingStep } from "@/stores/booking-store";

interface BookingStepIndicatorProps {
  currentStep: BookingStep;
}

const steps: { key: BookingStep; label: string; icon: typeof Ticket }[] = [
  { key: "select", label: "Select", icon: Ticket },
  { key: "review", label: "Review", icon: ClipboardList },
  { key: "payment", label: "Payment", icon: CreditCard },
  { key: "confirmation", label: "Done", icon: PartyPopper },
];

const stepOrder: BookingStep[] = ["select", "review", "payment", "confirmation"];

function getStepIndex(step: BookingStep) {
  return stepOrder.indexOf(step);
}

export function BookingStepIndicator({ currentStep }: BookingStepIndicatorProps) {
  const currentIndex = getStepIndex(currentStep);

  return (
    <div className="flex items-center justify-center gap-1 sm:gap-2">
      {steps.map((step, index) => {
        const isCompleted = index < currentIndex;
        const isActive = index === currentIndex;
        const Icon = step.icon;

        return (
          <div key={step.key} className="flex items-center gap-1 sm:gap-2">
            {index > 0 && (
              <div
                className={`h-px w-6 sm:w-10 ${
                  isCompleted ? "bg-primary" : "bg-border"
                }`}
              />
            )}
            <div className="flex flex-col items-center gap-1">
              <div
                className={`flex h-8 w-8 items-center justify-center rounded-full text-xs font-medium transition-colors ${
                  isCompleted
                    ? "bg-primary text-primary-foreground"
                    : isActive
                      ? "bg-primary text-primary-foreground"
                      : "bg-muted text-muted-foreground"
                }`}
              >
                {isCompleted ? (
                  <Check className="h-4 w-4" />
                ) : (
                  <Icon className="h-4 w-4" />
                )}
              </div>
              <span
                className={`text-xs ${
                  isActive ? "font-medium text-foreground" : "text-muted-foreground"
                }`}
              >
                {step.label}
              </span>
            </div>
          </div>
        );
      })}
    </div>
  );
}
