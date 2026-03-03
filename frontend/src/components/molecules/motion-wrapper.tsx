"use client";

import { m } from "framer-motion";
import type { Variants } from "framer-motion";
import { fadeIn } from "@/lib/motion";
import { cn } from "@/lib/utils";

interface MotionSectionProps {
  variants?: Variants;
  className?: string;
  children: React.ReactNode;
}

export function MotionSection({
  variants = fadeIn,
  className,
  children,
}: MotionSectionProps) {
  return (
    <m.section
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, margin: "-100px" }}
      variants={variants}
      className={cn(className)}
    >
      {children}
    </m.section>
  );
}

interface MotionDivProps {
  variants?: Variants;
  className?: string;
  children: React.ReactNode;
}

export function MotionDiv({
  variants = fadeIn,
  className,
  children,
}: MotionDivProps) {
  return (
    <m.div
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, margin: "-50px" }}
      variants={variants}
      className={cn(className)}
    >
      {children}
    </m.div>
  );
}
