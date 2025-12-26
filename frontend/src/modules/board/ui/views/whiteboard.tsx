import React, { useState, useRef, useCallback } from "react";
import { Excalidraw } from "@excalidraw/excalidraw";
import "@excalidraw/excalidraw/index.css";
import { convertToExcalidrawElements } from "@excalidraw/excalidraw";
import { useDebouncedCallback } from "@/lib/use-debounced-callback";
import BoardHeader from "./board-header";
import type { Board } from "../../types";
import type { ExcalidrawElementSkeleton } from "@excalidraw/excalidraw/data/transform";
import type { OrderedExcalidrawElement } from "@excalidraw/excalidraw/element/types";

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

export function convertFromExcalidrawElements(
  elements: readonly ExcalidrawElement[]
): Record<string, any>[] {
  const activeElements = elements.filter((el) => !el.isDeleted);

  return activeElements.map((element) => {
    const baseCommand: Record<string, any> = {
      type: element.type,
      id: element.id,
      x: element.x,
      y: element.y,
      width: element.width,
      height: element.height,
    };

    if (element.angle !== undefined && element.angle !== 0) {
      baseCommand.angle = element.angle;
    }
    if (element.strokeColor) {
      baseCommand.strokeColor = element.strokeColor;
    }
    if (element.backgroundColor && element.backgroundColor !== "transparent") {
      baseCommand.backgroundColor = element.backgroundColor;
    }
    if (element.fillStyle) {
      baseCommand.fillStyle = element.fillStyle;
    }
    if (element.strokeWidth !== undefined) {
      baseCommand.strokeWidth = element.strokeWidth;
    }
    if (element.strokeStyle) {
      baseCommand.strokeStyle = element.strokeStyle;
    }
    if (element.roughness !== undefined) {
      baseCommand.roughness = element.roughness;
    }
    if (element.opacity !== undefined && element.opacity !== 100) {
      baseCommand.opacity = element.opacity;
    }
    if (element.roundness !== null && element.roundness !== undefined) {
      baseCommand.roundness = element.roundness;
    }
    if (element.locked !== undefined) {
      baseCommand.locked = element.locked;
    }
    if (element.link) {
      baseCommand.link = element.link;
    }

    if (element.type === "arrow") {
      // Convert startBinding/endBinding to start/end format
      if ((element as any).startBinding?.elementId) {
        baseCommand.start = {
          id: (element as any).startBinding.elementId,
        };
      }
      if ((element as any).endBinding?.elementId) {
        baseCommand.end = {
          id: (element as any).endBinding.elementId,
        };
      }

      // Preserve label if it exists
      if ((element as any).label?.text) {
        baseCommand.label = {
          text: (element as any).label.text,
        };
        if ((element as any).label.strokeColor) {
          baseCommand.label.strokeColor = (element as any).label.strokeColor;
        }
      }

      // Preserve arrowhead styles if they exist
      if ((element as any).startArrowhead) {
        baseCommand.startArrowhead = (element as any).startArrowhead;
      }
      if ((element as any).endArrowhead) {
        baseCommand.endArrowhead = (element as any).endArrowhead;
      }
    } else {
      if (element.boundElements && element.boundElements.length > 0) {
        baseCommand.boundElements = element.boundElements.map((bound) => ({
          id: bound.id,
          type: bound.type,
        }));
      }
    }

    return baseCommand;
  });
}

export const Whiteboard = ({ board, onStateChange }: WhiteboardProps) => {
  const boardElements = (board.elements ?? []) as ExcalidrawElementSkeleton[];

  const initialElements: OrderedExcalidrawElement[] =
    convertToExcalidrawElements(boardElements, {
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
    <div className="h-full w-full relative flex flex-col">
      <BoardHeader boardId={board.id} boardName={board.name} />
      <div
        className="flex-1 relative"
        style={{ height: "calc(100% - 3.5rem)" }}
      >
        <div className="excalidraw-wrapper" style={{ height: "100%" }}>
          <Excalidraw
            excalidrawAPI={handleAPI}
            onChange={handleChange}
            initialData={{ elements, appState: { theme: "light" } }}
          />
        </div>
      </div>
    </div>
  );
};
