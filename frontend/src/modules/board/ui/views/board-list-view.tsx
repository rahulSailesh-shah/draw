import { Button } from "@/components/ui/button";
import { useMutationCreateBoard } from "../../hooks/use-board";
import { useNavigate } from "@tanstack/react-router";
import { generateSlug } from "random-word-slugs";

export const BoardListView = () => {
  const createBoard = useMutationCreateBoard();
  const navigate = useNavigate();

  const handleCreateBoard = () => {
    const boardName = generateSlug();
    createBoard.mutate(boardName, {
      onSuccess: (data) => {
        if (!data?.token || !data?.boardId) {
          return;
        }
        sessionStorage.setItem(`room-token`, data.token);
        navigate({
          to: "/boards/$boardId",
          params: { boardId: data.boardId.toString() },
        });
      },
      onError: (error) => {
        console.error("Failed to create board:", error);
      },
    });
  };

  return (
    <div className="flex-1 flex justify-center">
      <Button onClick={handleCreateBoard}>Create board</Button>
    </div>
  );
};
