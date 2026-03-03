import type { Metadata } from "next";
import { PageHeader } from "@/components/molecules/page-header";
import { AvatarWithName } from "@/components/molecules/avatar-with-name";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";

export const metadata: Metadata = {
  title: "Profile",
  description: "Manage your account settings.",
};

export default function ProfilePage() {
  return (
    <>
      <PageHeader title="Profile" description="Manage your account settings" />

      <div className="mt-6 grid gap-6 lg:grid-cols-3">
        {/* Profile card */}
        <Card>
          <CardContent className="flex flex-col items-center py-8">
            <AvatarWithName
              name="John Doe"
              email="john@example.com"
              size="lg"
            />
            <Separator className="my-6" />
            <div className="w-full space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Role</span>
                <span className="font-medium">Individual</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Member since</span>
                <span className="font-medium">Mar 2026</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Total bookings</span>
                <span className="font-medium">5</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Edit form */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Personal Information</CardTitle>
          </CardHeader>
          <CardContent>
            <form className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="firstName">First name</Label>
                  <Input id="firstName" defaultValue="John" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="lastName">Last name</Label>
                  <Input id="lastName" defaultValue="Doe" />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  defaultValue="john@example.com"
                  disabled
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="phone">Phone</Label>
                <Input
                  id="phone"
                  type="tel"
                  placeholder="+84 xxx xxx xxx"
                />
              </div>

              <div className="flex justify-end">
                <Button type="submit">Save changes</Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </>
  );
}
