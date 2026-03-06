"use client";

import { useState, useMemo } from "react";
import { PageHeader } from "@/components/molecules/page-header";
import { StatusBadge } from "@/components/molecules/status-badge";
import { AvatarWithName } from "@/components/molecules/avatar-with-name";
import { SearchInput } from "@/components/molecules/search-input";
import { Pagination } from "@/components/molecules/pagination";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";
import { useAdminUsers } from "@/lib/api/queries";
import { getTotalPages } from "@/lib/utils";
import type { AuthUser } from "@/lib/api/types/auth";

const LIMIT = 20;

const columns: Column<AuthUser>[] = [
  {
    key: "name",
    header: "User",
    render: (user) => (
      <AvatarWithName
        name={`${user.firstName} ${user.lastName}`.trim()}
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
    key: "isActive",
    header: "Status",
    render: (user) => (
      <StatusBadge
        status={
          user.isActive === false
            ? "inactive"
            : user.isVerified
              ? "active"
              : "pending"
        }
      />
    ),
  },
];

export default function AdminUsersPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");

  const { data, isLoading } = useAdminUsers({ page, limit: LIMIT });
  const users = data?.items ?? [];
  const totalPages = getTotalPages(data?.total ?? 0, LIMIT);

  const filtered = useMemo(() => {
    if (!search) return users;
    const q = search.toLowerCase();
    return users.filter(
      (u) =>
        u.firstName.toLowerCase().includes(q) ||
        u.lastName.toLowerCase().includes(q) ||
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

      <div className="mt-6">
        <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
      </div>
    </>
  );
}
