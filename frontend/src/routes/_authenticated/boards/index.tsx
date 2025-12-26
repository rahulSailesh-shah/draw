import { QueryBoundary } from "@/components/query-boundary";
import { useQueryGetBoards } from "@/modules/board/hooks/use-board";
import type { GetBoardsByUserIDResponse } from "@/modules/board/types";
import { BoardListView } from "@/modules/board/ui/views/board-list-view";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_authenticated/boards/")({
  component: RouteComponent,
});

function RouteComponent() {
  const getBoards = useQueryGetBoards();

  return (
    <QueryBoundary query={getBoards}>
      {(data: GetBoardsByUserIDResponse) => (
        <>
          <BoardListView boards={data?.boards ?? []} />
        </>
      )}
    </QueryBoundary>
  );
}
