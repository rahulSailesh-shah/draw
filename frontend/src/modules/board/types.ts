export interface CreateBoardResponse {
  boardId: string;
  token: string;
}

export interface Board {
  id: string;
  name: string;
  ownerId: string;
  elements: Record<string, any>[];
}

export interface GetBoardResponse {
  board: Board;
}

export interface UpdateBoardRequest {
  name?: string;
  elements?: Record<string, any>[];
}
