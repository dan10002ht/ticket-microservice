import type { Metadata } from "next";
import { EventCard } from "@/components/organisms/events/event-card";
import { PageHeader } from "@/components/molecules/page-header";
import { SearchInput } from "@/components/molecules/search-input";
import type { Event } from "@/types";

export const metadata: Metadata = {
  title: "Events",
  description: "Browse all upcoming events and find your next experience.",
};

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
    description: "Annual tech conference",
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
    description: "An evening of laughs",
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
  {
    id: "4",
    title: "Jazz Night at the Rooftop",
    description: "Smooth jazz under the stars",
    venue: "Rex Hotel Rooftop",
    address: "TP. Hồ Chí Minh",
    startDate: "2026-09-05T20:00:00Z",
    endDate: "2026-09-05T23:00:00Z",
    status: "published",
    organizerId: "org-1",
    totalCapacity: 150,
    availableCapacity: 80,
    minPrice: 400000,
    maxPrice: 1200000,
    category: "Music",
    createdAt: "2026-03-01T00:00:00Z",
  },
  {
    id: "5",
    title: "Startup Pitch Day",
    description: "Watch startups pitch to investors",
    venue: "Dreamplex",
    address: "TP. Hồ Chí Minh",
    startDate: "2026-10-12T08:00:00Z",
    endDate: "2026-10-12T17:00:00Z",
    status: "published",
    organizerId: "org-2",
    totalCapacity: 200,
    availableCapacity: 120,
    minPrice: 100000,
    maxPrice: 500000,
    category: "Business",
    createdAt: "2026-03-01T00:00:00Z",
  },
  {
    id: "6",
    title: "Art Exhibition: Modern Vietnam",
    description: "Contemporary art from Vietnamese artists",
    venue: "Vietnam Fine Arts Museum",
    address: "Hà Nội",
    startDate: "2026-11-01T10:00:00Z",
    endDate: "2026-11-30T18:00:00Z",
    status: "published",
    organizerId: "org-3",
    totalCapacity: 500,
    availableCapacity: 400,
    minPrice: 50000,
    maxPrice: 150000,
    category: "Art",
    createdAt: "2026-03-01T00:00:00Z",
  },
];

export default function EventsPage() {
  return (
    <section className="container mx-auto px-4 py-8">
      <PageHeader
        title="All Events"
        description="Browse and discover upcoming events"
      />

      <div className="mt-6">
        <SearchInput
          placeholder="Search events by name, venue..."
          onSearch={() => {}}
          className="max-w-md"
        />
      </div>

      <div className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {sampleEvents.map((event) => (
          <EventCard key={event.id} event={event} />
        ))}
      </div>
    </section>
  );
}
