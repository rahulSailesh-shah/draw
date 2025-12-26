export interface CreateBoardResponse {
  boardId: string;
}

export interface Board {
  id: string;
  name: string;
  ownerId: string;
  elements: Record<string, any>[];
}

export interface GetBoardResponse {
  board: Board;
  token: string;
}

export interface GetBoardsByUserIDResponse {
  boards: Board[];
}

export interface UpdateBoardRequest {
  name?: string;
  elements?: Record<string, any>[];
}
