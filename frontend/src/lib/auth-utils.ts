import { queryOptions } from "@tanstack/react-query";
import { authClient } from "./auth-client";
import { redirect } from "@tanstack/react-router";
import { queryClient } from "./query-client";

export const getSession = async () => {
  const sessionQueryOptions = queryOptions({
    queryKey: ["session"],
    queryFn: async () => {
      const { data: session } = await authClient.getSession();
      return session ?? null;
    },
    staleTime: 1000 * 60 * 5,
    gcTime: 1000 * 60 * 10,
  });

  return await queryClient.ensureQueryData(sessionQueryOptions);
};

export async function requireAuth() {
  const session = await getSession();
  if (!session) {
    throw redirect({ to: "/login" });
  }
  return session;
}

export async function requireNoAuth() {
  const session = await getSession();
  if (session) {
    throw redirect({ to: "/" });
  }
  return session;
}
