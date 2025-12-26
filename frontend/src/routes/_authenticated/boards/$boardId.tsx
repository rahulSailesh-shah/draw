import { QueryBoundary } from "@/components/query-boundary";
import { useQueryGetBoard } from "@/modules/board/hooks/use-board";
import type { GetBoardResponse } from "@/modules/board/types";
import { BoardRoomView } from "@/modules/board/ui/views/board-room-view";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_authenticated/boards/$boardId")({
  component: RouteComponent,
});

function RouteComponent() {
  const { boardId } = Route.useParams();
  const boardQuery = useQueryGetBoard(boardId);

  return (
    <>
      <QueryBoundary query={boardQuery}>
        {(data: GetBoardResponse) => (
          <>
            <BoardRoomView board={data.board} />
          </>
        )}
      </QueryBoundary>
    </>
  );
}
