import type { ReactNode } from "react";
import type { UseQueryResult } from "@tanstack/react-query";
import { AlertCircleIcon, FileX, Loader2Icon } from "lucide-react";

export const DefaultLoadingFallback = () => (
  <div className="flex items-center justify-center h-full">
    <Loader2Icon className="size-8 animate-spin text-muted-foreground" />
  </div>
);

export const DefaultErrorFallback = () => (
  <div className="flex flex-col items-center justify-center h-full gap-2">
    <AlertCircleIcon className="size-8 text-destructive" />
    <p className="text-sm text-muted-foreground">
      Something went wrong hkjebjikhw
    </p>
  </div>
);

export const DefaultEmptyFallback = () => (
  <div className="flex flex-col items-center justify-center h-full gap-2">
    <FileX className="size-8 text-muted-foreground" />
    <p className="text-sm text-muted-foreground">No data found</p>
  </div>
);

interface QueryBoundaryProps<T> {
  query: UseQueryResult<T | null | undefined>;
  loadingFallback?: ReactNode;
  errorFallback?: ReactNode;
  emptyFallback?: ReactNode;
  children: (data: NonNullable<T>) => ReactNode;
}

export function QueryBoundary<T>({
  query,
  loadingFallback = <DefaultLoadingFallback />,
  errorFallback = <DefaultErrorFallback />,
  emptyFallback = <DefaultEmptyFallback />,
  children,
}: QueryBoundaryProps<T>) {
  if (query.isLoading) {
    return <>{loadingFallback}</>;
  }

  if (query.isError) {
    return <>{errorFallback}</>;
  }

  if (query.data === null || query.data === undefined) {
    return <>{emptyFallback}</>;
  }

  return <>{children(query.data)}</>;
}
