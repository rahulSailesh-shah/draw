import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useMutationCreateBoard } from "../../hooks/use-board";
import { useNavigate } from "@tanstack/react-router";
import { generateSlug } from "random-word-slugs";
import type { Board } from "../../types";
import { PlusIcon, FileTextIcon, Loader2Icon } from "lucide-react";
import { cn } from "@/lib/utils";

export const BoardListView = ({ boards }: { boards: Board[] }) => {
  const createBoard = useMutationCreateBoard();
  const navigate = useNavigate();

  const handleCreateBoard = () => {
    const boardName = generateSlug();
    createBoard.mutate(boardName, {
      onSuccess: (data) => {
        if (!data?.boardId) {
          return;
        }
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

  const handleBoardClick = (boardId: string) => {
    navigate({
      to: "/boards/$boardId",
      params: { boardId },
    });
  };

  return (
    <div className="flex-1 flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center justify-between px-6 py-4 border-b bg-background">
        <div>
          <h1 className="text-2xl font-semibold text-foreground">Boards</h1>
          <p className="text-sm text-muted-foreground mt-1">
            {boards.length === 0
              ? "No boards yet"
              : `${boards.length} ${boards.length === 1 ? "board" : "boards"}`}
          </p>
        </div>
        <Button
          onClick={handleCreateBoard}
          disabled={createBoard.isPending}
          className="gap-2"
        >
          {createBoard.isPending ? (
            <Loader2Icon className="size-4 animate-spin" />
          ) : (
            <PlusIcon className="size-4" />
          )}
          Create board
        </Button>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto">
        {boards.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full max-w-md mx-auto text-center px-6">
            <div className="rounded-full bg-muted p-6 mb-4">
              <FileTextIcon className="size-12 text-muted-foreground" />
            </div>
            <h2 className="text-xl font-semibold text-foreground mb-2">
              No boards yet
            </h2>
            <p className="text-sm text-muted-foreground mb-6">
              Get started by creating your first board. You can collaborate,
              draw, and share ideas.
            </p>
            <Button
              onClick={handleCreateBoard}
              disabled={createBoard.isPending}
              className="gap-2"
            >
              {createBoard.isPending ? (
                <Loader2Icon className="size-4 animate-spin" />
              ) : (
                <PlusIcon className="size-4" />
              )}
              Create your first board
            </Button>
          </div>
        ) : (
          <div className="max-w-7xl mx-auto px-6 py-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {boards.map((board) => (
                <Card
                  key={board.id}
                  onClick={() => handleBoardClick(board.id)}
                  className={cn(
                    "cursor-pointer transition-all hover:shadow-md hover:border-primary/50 group",
                    "hover:-translate-y-0.5"
                  )}
                >
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2 group-hover:text-primary transition-colors">
                      <FileTextIcon className="size-5 text-muted-foreground group-hover:text-primary transition-colors" />
                      <span className="truncate">{board.name}</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground">
                      Click to open
                    </p>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
