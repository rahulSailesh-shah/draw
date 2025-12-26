import {
  useMutation,
  useQueryClient,
  useQuery,
  keepPreviousData,
} from "@tanstack/react-query";
import { createBoard, getBoard, getBoards, updateBoard } from "../api";
import type { UpdateBoardRequest } from "../types";

// Query hooks
export const useQueryGetBoard = (id: string) => {
  return useQuery({
    queryKey: ["board", id],
    queryFn: () => getBoard(id),
  });
};

export const useQueryGetBoards = () => {
  return useQuery({
    queryKey: ["boards"],
    queryFn: () => getBoards(),
    placeholderData: keepPreviousData,
  });
};

// Mutation hooks
export const useMutationCreateBoard = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (name: string) => createBoard(name),
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: ["board", data?.boardId.toString()],
      });
    },
    onError: () => {},
  });
};

export const useMutationUpdateBoard = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, req }: { id: string; req: UpdateBoardRequest }) =>
      updateBoard(id, req),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["board", variables.id],
      });
    },
    onError: (error) => {
      console.error(error);
    },
  });
};
