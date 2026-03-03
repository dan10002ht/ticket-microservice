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
import { EventCard } from "@/components/organisms/events/event-card";
import { MotionSection, MotionDiv } from "@/components/molecules/motion-wrapper";
import { fadeIn, fadeInUp, staggerContainer } from "@/lib/motion";
import type { Event } from "@/types";

const sampleEvents: Event[] = [
  {
    id: "1",
    title: "Summer Music Festival 2026",
    description: "The biggest music festival of the year",
    venue: "Sân vận động Mỹ Đình",
    address: "Hà Nội",
    startDate: "2026-07-15T18:00:00Z",
    endDate: "2026-07-15T23:00:00Z",
    status: "published",
    organizerId: "org-1",
    totalCapacity: 5000,
    availableCapacity: 1200,
    minPrice: 500000,
    maxPrice: 2000000,
    category: "Music",
    createdAt: "2026-01-01T00:00:00Z",
  },
  {
    id: "2",
    title: "Tech Conference Vietnam",
    description: "Annual tech conference for developers",
    venue: "GEM Center",
    address: "TP. Hồ Chí Minh",
    startDate: "2026-08-20T09:00:00Z",
    endDate: "2026-08-21T17:00:00Z",
    status: "published",
    organizerId: "org-2",
    totalCapacity: 800,
    availableCapacity: 350,
    minPrice: 300000,
    maxPrice: 1500000,
    category: "Tech",
    createdAt: "2026-02-01T00:00:00Z",
  },
  {
    id: "3",
    title: "Stand-up Comedy Night",
    description: "An evening of laughs with top comedians",
    venue: "Nhà hát Lớn",
    address: "Hà Nội",
    startDate: "2026-06-10T19:30:00Z",
    endDate: "2026-06-10T22:00:00Z",
    status: "published",
    organizerId: "org-3",
    totalCapacity: 300,
    availableCapacity: 45,
    minPrice: 200000,
    maxPrice: 800000,
    category: "Comedy",
    createdAt: "2026-03-01T00:00:00Z",
  },
];

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

export default function HomePage() {
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
          {sampleEvents.map((event) => (
            <MotionDiv key={event.id} variants={fadeIn}>
              <EventCard event={event} />
            </MotionDiv>
          ))}
        </div>
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
