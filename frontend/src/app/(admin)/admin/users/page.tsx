import type { Metadata } from "next";
import { PageHeader } from "@/components/molecules/page-header";
import { StatusBadge } from "@/components/molecules/status-badge";
import { AvatarWithName } from "@/components/molecules/avatar-with-name";
import { SearchInput } from "@/components/molecules/search-input";
import { DataTable, type Column } from "@/components/organisms/shared/data-table";

export const metadata: Metadata = {
  title: "Manage Users",
  description: "View and manage all platform users.",
};

interface AdminUser {
  id: string;
  name: string;
  email: string;
  role: string;
  status: string;
  joinedAt: string;
}

const users: AdminUser[] = [
  { id: "1", name: "Nguyen Van A", email: "a@example.com", role: "individual", status: "active", joinedAt: "Jan 15, 2026" },
  { id: "2", name: "Tran Thi B", email: "b@example.com", role: "organization", status: "active", joinedAt: "Feb 01, 2026" },
  { id: "3", name: "Le Van C", email: "c@example.com", role: "individual", status: "active", joinedAt: "Feb 20, 2026" },
  { id: "4", name: "Pham Thi D", email: "d@example.com", role: "organization", status: "pending", joinedAt: "Mar 01, 2026" },
  { id: "5", name: "Hoang Van E", email: "e@example.com", role: "individual", status: "cancelled", joinedAt: "Mar 02, 2026" },
];

const columns: Column<AdminUser>[] = [
  {
    key: "name",
    header: "User",
    render: (user) => (
      <AvatarWithName name={user.name} email={user.email} size="sm" />
    ),
  },
  { key: "role", header: "Role" },
  {
    key: "status",
    header: "Status",
    render: (user) => <StatusBadge status={user.status} />,
  },
  { key: "joinedAt", header: "Joined" },
];

export default function AdminUsersPage() {
  return (
    <>
      <PageHeader
        title="Users"
        description="Manage all platform users"
      />

      <div className="mt-6">
        <SearchInput
          placeholder="Search by name or email..."
          onSearch={() => {}}
          className="max-w-md"
        />
      </div>

      <div className="mt-6">
        <DataTable
          columns={columns}
          data={users}
          keyExtractor={(u) => u.id}
        />
      </div>
    </>
  );
}
