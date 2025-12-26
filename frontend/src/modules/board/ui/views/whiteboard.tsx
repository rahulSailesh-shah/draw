import React, { useState, useRef, useCallback } from "react";
import { Excalidraw } from "@excalidraw/excalidraw";
import "@excalidraw/excalidraw/index.css";
import { convertToExcalidrawElements } from "@excalidraw/excalidraw";
import { useDebouncedCallback } from "@/lib/use-debounced-callback";
import BoardHeader from "./board-header";
import type { Board } from "../../types";

export interface WhiteboardStateChange {
  elements: readonly ExcalidrawElement[];
  appState?: unknown;
  files?: unknown;
}

export interface WhiteboardProps {
  board: Board;
  onStateChange?: (state: WhiteboardStateChange) => void;
}

type ExcalidrawOnChange = Parameters<
  NonNullable<React.ComponentProps<typeof Excalidraw>["onChange"]>
>[0];
type ExcalidrawElement = ExcalidrawOnChange[number];
type ExcalidrawAPI = Parameters<
  NonNullable<React.ComponentProps<typeof Excalidraw>["excalidrawAPI"]>
>[0];

export const Whiteboard = ({ board, onStateChange }: WhiteboardProps) => {
  const initialElements = convertToExcalidrawElements([], {
    regenerateIds: false,
  });

  const excalidrawAPI = useRef<ExcalidrawAPI | null>(null);
  const previousElementsRef =
    useRef<readonly ExcalidrawElement[]>(initialElements);

  const [elements, setElements] =
    useState<readonly ExcalidrawElement[]>(initialElements);

  const elementsHaveChanged = useCallback(
    (
      prevElements: readonly ExcalidrawElement[],
      newElements: readonly ExcalidrawElement[]
    ): boolean => {
      if (prevElements.length !== newElements.length) {
        return true;
      }

      const prevMap = new Map(prevElements.map((el) => [el.id, el]));
      const newMap = new Map(newElements.map((el) => [el.id, el]));

      if (prevMap.size !== newMap.size) {
        return true;
      }

      let hasChanges = false;
      for (const [id, newEl] of newMap) {
        const prevEl = prevMap.get(id);
        if (!prevEl) {
          return true;
        }

        if (prevEl.version !== newEl.version) {
          hasChanges = true;
          break;
        }

        if (prevEl.isDeleted !== newEl.isDeleted) {
          hasChanges = true;
          break;
        }

        const positionChanged =
          Math.abs(prevEl.x - newEl.x) > 0.01 ||
          Math.abs(prevEl.y - newEl.y) > 0.01;
        const sizeChanged =
          Math.abs(prevEl.width - newEl.width) > 0.01 ||
          Math.abs(prevEl.height - newEl.height) > 0.01;

        if (positionChanged || sizeChanged) {
          hasChanges = true;
          break;
        }
      }

      return hasChanges;
    },
    []
  );

  const notifyStateChange = useDebouncedCallback(
    (updatedElements: readonly ExcalidrawElement[], appState?: unknown) => {
      if (onStateChange) {
        onStateChange({
          elements: updatedElements,
          appState,
        });
      }
    }
  );

  const handleChange = useCallback(
    (updatedElements: readonly ExcalidrawElement[], appState: unknown) => {
      const elementsChanged = elementsHaveChanged(
        previousElementsRef.current,
        updatedElements
      );
      setElements(updatedElements);
      previousElementsRef.current = updatedElements.map((el) => ({ ...el }));

      if (elementsChanged) {
        notifyStateChange(updatedElements, appState);
      }
    },
    [notifyStateChange, elementsHaveChanged]
  );

  const handleAPI = useCallback((api: ExcalidrawAPI) => {
    excalidrawAPI.current = api;
  }, []);

  return (
    <div className="h-full w-full relative">
      <BoardHeader boardId={board.id} boardName={board.name} />
      <Excalidraw
        excalidrawAPI={handleAPI}
        onChange={handleChange}
        initialData={{
          elements,
          appState: { zenModeEnabled: true, viewBackgroundColor: "#a5d8ff" },
          scrollToContent: true,
        }}
      />
    </div>
  );
};
