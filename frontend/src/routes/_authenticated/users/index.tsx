import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { authClient } from "@/lib/auth-client";
import { queryClient } from "@/lib/query-client";
import { QueryBoundary } from "@/components/query-boundary";

export const Route = createFileRoute("/_authenticated/users/")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const userData = useQuery({
    queryKey: ["repoData"],
    queryFn: async () => {
      const response = await fetch(
        "https://api.github.com/repos/TanStack/query"
      );
      return await response.json();
    },
  });

  return (
    <>
      <QueryBoundary query={userData}>
        {(data: any) => (
          <>
            <div>
              <h1>{data.full_name}</h1>
              <p>{data.description}</p>
              <strong>👀 {data.subscribers_count}</strong>{" "}
              <strong>✨ {data.stargazers_count}</strong>{" "}
              <strong>🍴 {data.forks_count}</strong>
            </div>

            <Button
              onClick={async () => {
                try {
                  await authClient.signOut({
                    fetchOptions: {
                      onSuccess: () => {
                        queryClient.removeQueries({ queryKey: ["session"] });
                        router.navigate({ to: "/login", replace: true });
                      },
                    },
                  });
                } catch (err) {
                  console.error("Error signing out:", err);
                }
              }}
            >
              Sign out
            </Button>
          </>
        )}
      </QueryBoundary>
    </>
  );
}
