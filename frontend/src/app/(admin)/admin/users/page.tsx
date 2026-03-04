"use client";

import { useState, useMemo } from "react";
import { PageHeader } from "@/components/molecules/page-header";
import { StatusBadge } from "@/components/molecules/status-badge";
import { AvatarWithName } from "@/components/molecules/avatar-with-name";
import { SearchInput } from "@/components/molecules/search-input";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import { useAdminUsers } from "@/lib/api/queries";
import type { AuthUser } from "@/lib/api/types/auth";

const columns: Column<AuthUser>[] = [
  {
    key: "name",
    header: "User",
    render: (user) => (
      <AvatarWithName
        name={`${user.first_name} ${user.last_name}`.trim()}
        email={user.email}
        size="sm"
      />
    ),
  },
  {
    key: "role",
    header: "Role",
    render: (user) => (
      <span className="capitalize">{user.role ?? "user"}</span>
    ),
  },
  {
    key: "is_active",
    header: "Status",
    render: (user) => (
      <StatusBadge
        status={
          user.is_active === false
            ? "inactive"
            : user.is_verified
              ? "active"
              : "pending"
        }
      />
    ),
  },
];

export default function AdminUsersPage() {
  const { data, isLoading } = useAdminUsers();
  const users = data?.items ?? [];
  const [search, setSearch] = useState("");

  const filtered = useMemo(() => {
    if (!search) return users;
    const q = search.toLowerCase();
    return users.filter(
      (u) =>
        u.first_name.toLowerCase().includes(q) ||
        u.last_name.toLowerCase().includes(q) ||
        u.email.toLowerCase().includes(q)
    );
  }, [users, search]);

  return (
    <>
      <PageHeader
        title="Users"
        description="Manage all platform users"
      />

      <div className="mt-6">
        <SearchInput
          placeholder="Search by name or email..."
          onSearch={setSearch}
          className="max-w-md"
        />
      </div>

      <div className="mt-6">
        <DataTable
          columns={columns}
          data={filtered}
          isLoading={isLoading}
          emptyMessage="No users found"
          keyExtractor={(u) => u.id}
        />
      </div>
    </>
  );
}
