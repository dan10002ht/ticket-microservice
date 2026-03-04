"use client";

import Link from "next/link";
import {
  ArrowRight,
  CalendarDays,
  Music,
  Monitor,
  Laugh,
  Palette,
  Briefcase,
  Trophy,
  Shield,
  Zap,
  Users,
  Ticket,
  TrendingUp,
} from "lucide-react";
import { m } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { EventCard } from "@/components/organisms/events/event-card";
import { MotionSection, MotionDiv } from "@/components/molecules/motion-wrapper";
import { fadeIn, fadeInUp, staggerContainer } from "@/lib/motion";
import { useEvents } from "@/lib/api/queries";

const categories = [
  { label: "Music", icon: Music, href: "/events?category=music" },
  { label: "Tech", icon: Monitor, href: "/events?category=tech" },
  { label: "Comedy", icon: Laugh, href: "/events?category=comedy" },
  { label: "Art", icon: Palette, href: "/events?category=art" },
  { label: "Business", icon: Briefcase, href: "/events?category=business" },
  { label: "Sports", icon: Trophy, href: "/events?category=sports" },
];

const features = [
  {
    icon: CalendarDays,
    title: "Browse Events",
    description: "Discover upcoming events across Vietnam",
  },
  {
    icon: Zap,
    title: "Instant Booking",
    description: "Real-time seat selection and instant confirmation",
  },
  {
    icon: Shield,
    title: "Secure Payments",
    description: "Safe checkout powered by Stripe",
  },
];

const stats = [
  { value: "10,000+", label: "Happy Customers", icon: Users },
  { value: "500+", label: "Events Hosted", icon: CalendarDays },
  { value: "50,000+", label: "Tickets Sold", icon: Ticket },
  { value: "99.9%", label: "Uptime", icon: TrendingUp },
];

function EventCardSkeleton() {
  return (
    <div className="rounded-lg border overflow-hidden">
      <Skeleton className="aspect-[16/9] w-full" />
      <div className="p-4 space-y-2">
        <Skeleton className="h-5 w-3/4" />
        <Skeleton className="h-4 w-1/2" />
        <Skeleton className="h-4 w-2/3" />
      </div>
      <div className="border-t px-4 py-3 flex justify-between">
        <Skeleton className="h-4 w-20" />
        <Skeleton className="h-4 w-16" />
      </div>
    </div>
  );
}

export default function HomePage() {
  const { data, isLoading } = useEvents({ limit: 3 });
  const events = data?.items ?? [];

  return (
    <>
      {/* Hero */}
      <section className="relative overflow-hidden gradient-brand-subtle py-20 lg:py-32">
        {/* Decorative blurred shapes */}
        <div className="absolute -left-32 -top-32 h-64 w-64 rounded-full bg-primary/10 blur-3xl" />
        <div className="absolute -right-32 bottom-0 h-80 w-80 rounded-full bg-accent/20 blur-3xl" />

        <m.div
          initial="hidden"
          animate="visible"
          variants={staggerContainer}
          className="container relative mx-auto px-4 text-center"
        >
          <m.div variants={fadeIn}>
            <Badge variant="secondary" className="mb-6 px-4 py-1.5 text-sm">
              Vietnam&apos;s #1 Event Platform
            </Badge>
          </m.div>

          <m.h1
            variants={fadeInUp}
            className="mx-auto max-w-3xl text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl"
          >
            Discover, Book &{" "}
            <span className="text-gradient-brand">Experience</span>{" "}
            Amazing Events
          </m.h1>

          <m.p
            variants={fadeInUp}
            className="mx-auto mt-6 max-w-xl text-lg text-muted-foreground"
          >
            Browse concerts, conferences, and more. Book tickets instantly with
            real-time seat selection.
          </m.p>

          <m.div
            variants={fadeInUp}
            className="mt-8 flex flex-col items-center gap-4 sm:flex-row sm:justify-center"
          >
            <Button size="lg" asChild>
              <Link href="/events">
                Browse Events <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
            <Button variant="outline" size="lg" asChild>
              <Link href="/register">Create Account</Link>
            </Button>
          </m.div>
        </m.div>
      </section>

      {/* Features */}
      <MotionSection
        variants={staggerContainer}
        className="container mx-auto px-4 py-16"
      >
        <div className="grid gap-8 md:grid-cols-3">
          {features.map((feature) => (
            <MotionDiv key={feature.title} variants={fadeIn}>
              <div className="rounded-xl border bg-card p-6 text-center transition-shadow hover:shadow-md">
                <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-lg gradient-brand">
                  <feature.icon className="h-6 w-6 text-white" />
                </div>
                <h3 className="mt-4 font-semibold">{feature.title}</h3>
                <p className="mt-2 text-sm text-muted-foreground">
                  {feature.description}
                </p>
              </div>
            </MotionDiv>
          ))}
        </div>
      </MotionSection>

      {/* Categories */}
      <MotionSection
        variants={staggerContainer}
        className="container mx-auto px-4 py-16"
      >
        <h2 className="text-center text-2xl font-bold">Browse by Category</h2>
        <div className="mt-8 grid grid-cols-3 gap-3 sm:grid-cols-6">
          {categories.map((cat) => (
            <MotionDiv key={cat.label} variants={fadeIn}>
              <Link
                href={cat.href}
                className="flex flex-col items-center gap-2 rounded-xl border bg-card p-4 transition-all hover:border-primary hover:shadow-md"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                  <cat.icon className="h-5 w-5 text-primary" />
                </div>
                <span className="text-xs font-medium">{cat.label}</span>
              </Link>
            </MotionDiv>
          ))}
        </div>
      </MotionSection>

      {/* Upcoming Events */}
      <MotionSection
        variants={staggerContainer}
        className="container mx-auto px-4 py-16"
      >
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold">Upcoming Events</h2>
          <Button variant="ghost" asChild>
            <Link href="/events">
              View all <ArrowRight className="ml-1 h-4 w-4" />
            </Link>
          </Button>
        </div>
        <div className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {isLoading
            ? Array.from({ length: 3 }).map((_, i) => (
                <EventCardSkeleton key={i} />
              ))
            : events.map((event) => (
                <MotionDiv key={event.id} variants={fadeIn}>
                  <EventCard event={event} />
                </MotionDiv>
              ))}
        </div>
        {!isLoading && events.length === 0 && (
          <p className="mt-8 text-center text-muted-foreground">
            No upcoming events at the moment. Check back soon!
          </p>
        )}
      </MotionSection>

      {/* Social Proof */}
      <MotionSection
        variants={staggerContainer}
        className="gradient-brand-subtle py-16"
      >
        <div className="container mx-auto px-4">
          <h2 className="text-center text-2xl font-bold">
            Trusted by Event-Goers Nationwide
          </h2>
          <div className="mt-10 grid grid-cols-2 gap-6 lg:grid-cols-4">
            {stats.map((stat) => (
              <MotionDiv
                key={stat.label}
                variants={fadeIn}
                className="text-center"
              >
                <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                  <stat.icon className="h-6 w-6 text-primary" />
                </div>
                <p className="mt-3 text-3xl font-bold tabular-nums">
                  {stat.value}
                </p>
                <p className="mt-1 text-sm text-muted-foreground">
                  {stat.label}
                </p>
              </MotionDiv>
            ))}
          </div>
        </div>
      </MotionSection>
    </>
  );
}
