"use client";

import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { PageHeader } from "@/components/molecules/page-header";
import { AvatarWithName } from "@/components/molecules/avatar-with-name";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { useProfile, useUpdateProfile } from "@/lib/api/queries";
import { useAuthStore } from "@/stores";
import { profileSchema, type ProfileInput } from "@/lib/validators/profile";

function ProfileSkeleton() {
  return (
    <div className="mt-6 grid gap-6 lg:grid-cols-3">
      <Card>
        <CardContent className="flex flex-col items-center py-8 space-y-4">
          <Skeleton className="h-16 w-16 rounded-full" />
          <Skeleton className="h-4 w-32" />
          <Skeleton className="h-3 w-40" />
          <Separator className="my-6" />
          <div className="w-full space-y-2">
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-full" />
          </div>
        </CardContent>
      </Card>
      <Card className="lg:col-span-2">
        <CardHeader>
          <Skeleton className="h-6 w-48" />
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
          </div>
          <Skeleton className="h-10 w-full" />
          <Skeleton className="h-10 w-full" />
        </CardContent>
      </Card>
    </div>
  );
}

export default function ProfilePage() {
  const { user } = useAuthStore();
  const { data: profile, isLoading } = useProfile();
  const updateMutation = useUpdateProfile();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isDirty },
  } = useForm<ProfileInput>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      phone: "",
    },
  });

  // Populate form when profile data loads
  useEffect(() => {
    if (profile) {
      reset({
        firstName: profile.firstName || "",
        lastName: profile.lastName || "",
        phone: profile.phone || "",
      });
    }
  }, [profile, reset]);

  const onSubmit = async (data: ProfileInput) => {
    try {
      await updateMutation.mutateAsync({
        firstName: data.firstName,
        lastName: data.lastName,
        phone: data.phone || undefined,
      });
      toast.success("Profile updated successfully");
    } catch {
      toast.error("Failed to update profile");
    }
  };

  const displayName = profile
    ? `${profile.firstName} ${profile.lastName}`.trim()
    : user
      ? `${user.firstName} ${user.lastName}`.trim()
      : "User";

  const displayEmail = profile?.email || user?.email || "";

  const memberSince = profile?.createdAt
    ? new Date(profile.createdAt).toLocaleDateString("vi-VN", {
        month: "short",
        year: "numeric",
      })
    : "—";

  return (
    <>
      <PageHeader title="Profile" description="Manage your account settings" />

      {isLoading ? (
        <ProfileSkeleton />
      ) : (
        <div className="mt-6 grid gap-6 lg:grid-cols-3">
          {/* Profile card */}
          <Card>
            <CardContent className="flex flex-col items-center py-8">
              <AvatarWithName
                name={displayName}
                email={displayEmail}
                size="lg"
              />
              <Separator className="my-6" />
              <div className="w-full space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Role</span>
                  <span className="font-medium capitalize">
                    {user?.role || "Individual"}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Member since</span>
                  <span className="font-medium">{memberSince}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Status</span>
                  <span className="font-medium">
                    {user?.isVerified ? "Verified" : "Unverified"}
                  </span>
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
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="firstName">First name</Label>
                    <Input id="firstName" {...register("firstName")} />
                    {errors.firstName && (
                      <p className="text-xs text-destructive">
                        {errors.firstName.message}
                      </p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="lastName">Last name</Label>
                    <Input id="lastName" {...register("lastName")} />
                    {errors.lastName && (
                      <p className="text-xs text-destructive">
                        {errors.lastName.message}
                      </p>
                    )}
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    value={displayEmail}
                    disabled
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="phone">Phone</Label>
                  <Input
                    id="phone"
                    type="tel"
                    placeholder="+84 xxx xxx xxx"
                    {...register("phone")}
                  />
                  {errors.phone && (
                    <p className="text-xs text-destructive">
                      {errors.phone.message}
                    </p>
                  )}
                </div>

                <div className="flex justify-end">
                  <Button
                    type="submit"
                    disabled={!isDirty || updateMutation.isPending}
                  >
                    {updateMutation.isPending
                      ? "Saving..."
                      : "Save changes"}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      )}
    </>
  );
}
