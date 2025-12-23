import { requireAuth } from "@/lib/auth-utils";
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_authenticated")({
  beforeLoad: requireAuth,
  component: () => <Outlet />,
});
