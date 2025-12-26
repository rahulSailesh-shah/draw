import { isWeb } from "@livekit/components-core";
import * as React from "react";
import {
  ControlBar,
  LayoutContextProvider,
  RoomAudioRenderer,
  ConnectionStateToast,
  useCreateLayoutContext,
  useRoomContext,
} from "@livekit/components-react";

export interface DraggableControlsLayoutProps extends React.HTMLAttributes<HTMLDivElement> {}

// Sentiment analysis result
export interface SentimentData {
  text: string;
  sentiment: string;
  score: number;
  emotions: Record<string, number>;
  timestamp: string;
  source: string;
}

// Transcript message
export interface TranscriptData {
  role: string; // "user" or "ai"
  name: string; // Speaker's name
  content: string; // Transcript text
  timestamp: string; // When the segment was captured
}

// Wrapper for all stream messages
export interface StreamTextData {
  type: "sentiment" | "transcript";
  data: SentimentData | TranscriptData;
}

// Custom hook for dragging functionality
const useDraggable = (initialPosition = { x: 0, y: 0 }) => {
  // Initialize to center-top position
  const [position, setPosition] = React.useState(() => {
    if (typeof window !== "undefined") {
      return {
        x: window.innerWidth - 100,
        y: 0,
      };
    }
    return initialPosition;
  });
  const [isDragging, setIsDragging] = React.useState(false);
  const [hasBeenDragged, setHasBeenDragged] = React.useState(false);
  const dragRef = React.useRef<HTMLDivElement>(null);
  const startPos = React.useRef({ x: 0, y: 0 });

  const handleMouseDown = React.useCallback((e: React.MouseEvent) => {
    if (!dragRef.current) return;

    // Don't start dragging if clicking on a button or interactive element
    const target = e.target as HTMLElement;
    if (
      target.tagName === "BUTTON" ||
      target.closest("button") ||
      target.closest("[role='button']") ||
      target.closest("a") ||
      target.closest("input") ||
      target.closest("select")
    ) {
      return;
    }

    setIsDragging(true);
    const rect = dragRef.current.getBoundingClientRect();
    // Calculate offset from mouse position to element's top-left corner
    startPos.current = {
      x: e.clientX - rect.left,
      y: e.clientY - rect.top,
    };
    e.preventDefault();
    e.stopPropagation();
  }, []);

  const handleTouchStart = React.useCallback((e: React.TouchEvent) => {
    if (!dragRef.current) return;

    // Don't start dragging if touching a button or interactive element
    const target = e.target as HTMLElement;
    if (
      target.tagName === "BUTTON" ||
      target.closest("button") ||
      target.closest("[role='button']") ||
      target.closest("a") ||
      target.closest("input") ||
      target.closest("select")
    ) {
      return;
    }

    setIsDragging(true);
    const touch = e.touches[0];
    const rect = dragRef.current.getBoundingClientRect();
    startPos.current = {
      x: touch.clientX - rect.left,
      y: touch.clientY - rect.top,
    };
    e.preventDefault();
    e.stopPropagation();
  }, []);

  React.useEffect(() => {
    if (!isDragging) return;

    const handleMouseMove = (e: MouseEvent) => {
      if (!dragRef.current) return;

      setHasBeenDragged(true);
      const rect = dragRef.current.getBoundingClientRect();

      // Calculate new absolute position
      const newX = e.clientX - startPos.current.x;
      const newY = e.clientY - startPos.current.y;

      // Constrain to viewport bounds
      const maxX = window.innerWidth - rect.width;
      const maxY = window.innerHeight - rect.height;
      
      // Minimum top position to prevent toolbar from going above viewport
      const minY = 0;

      setPosition({
        x: Math.max(0, Math.min(newX, maxX)),
        y: Math.max(minY, Math.min(newY, maxY)),
      });
    };

    const handleTouchMove = (e: TouchEvent) => {
      if (!dragRef.current) return;
      const touch = e.touches[0];

      setHasBeenDragged(true);
      const rect = dragRef.current.getBoundingClientRect();

      // Calculate new absolute position
      const newX = touch.clientX - startPos.current.x;
      const newY = touch.clientY - startPos.current.y;

      // Constrain to viewport bounds
      const maxX = window.innerWidth - rect.width;
      const maxY = window.innerHeight - rect.height;
      
      // Minimum top position to prevent toolbar from going above viewport
      const minY = 0;

      setPosition({
        x: Math.max(0, Math.min(newX, maxX)),
        y: Math.max(minY, Math.min(newY, maxY)),
      });
    };

    const handleMouseUp = () => {
      setIsDragging(false);
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
    document.addEventListener("touchmove", handleTouchMove, { passive: false });
    document.addEventListener("touchend", handleMouseUp);

    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
      document.removeEventListener("touchmove", handleTouchMove);
      document.removeEventListener("touchend", handleMouseUp);
    };
  }, [isDragging]);

  return {
    dragRef,
    position,
    isDragging,
    hasBeenDragged,
    handleMouseDown,
    handleTouchStart,
  };
};

export const DraggableControlsLayout = ({
  ...props
}: DraggableControlsLayoutProps) => {
  const room = useRoomContext();
  const layoutContext = useCreateLayoutContext();
  const [_transcripts, setTranscripts] = React.useState<TranscriptData[]>([]);
  const [_currentSentiment, setCurrentSentiment] = React.useState<{
    sentiment: string;
    score: number;
    maxEmotion: string;
    maxEmotionScore: number;
  } | null>(null);

  const {
    dragRef,
    position,
    isDragging,
    hasBeenDragged,
    handleMouseDown,
    handleTouchStart,
  } = useDraggable();

  // Ensure toolbar stays visible and fixed when Excalidraw changes layout
  React.useEffect(() => {
    if (!dragRef.current) return;

    const element = dragRef.current;
    
    const ensureFixedPosition = () => {
      // Force fixed positioning relative to viewport, not affected by parent transforms
      const computedStyle = window.getComputedStyle(element);
      if (computedStyle.position !== 'fixed') {
        element.style.position = 'fixed';
      }
      
      // Ensure toolbar doesn't go above viewport
      const rect = element.getBoundingClientRect();
      if (rect.top < 0) {
        element.style.top = '0px';
      }
    };

    // Use MutationObserver to detect when Excalidraw changes the DOM
    const observer = new MutationObserver(() => {
      ensureFixedPosition();
    });

    // Observe changes to the document body and Excalidraw containers
    observer.observe(document.body, {
      childList: true,
      subtree: true,
      attributes: true,
      attributeFilter: ['style', 'class'],
    });

    // Also check on resize and scroll
    window.addEventListener('resize', ensureFixedPosition);
    window.addEventListener('scroll', ensureFixedPosition, true);

    // Initial check
    ensureFixedPosition();

    return () => {
      observer.disconnect();
      window.removeEventListener('resize', ensureFixedPosition);
      window.removeEventListener('scroll', ensureFixedPosition, true);
    };
  }, [dragRef]);

  React.useEffect(() => {
    try {
      room.registerTextStreamHandler(
        "room",
        async (reader, participantInfo) => {
          const message = await reader.readAll();
          const streamData: StreamTextData = JSON.parse(message);
          console.log("Received stream data:", streamData);

          if (streamData.type === "sentiment") {
            const sentimentData = streamData.data as SentimentData;
            console.log(
              `Sentiment from ${participantInfo.identity}:`,
              sentimentData
            );

            // Find the maximum emotion
            const emotions = sentimentData.emotions;
            let maxEmotion = "";
            let maxEmotionScore = 0;

            for (const [emotion, score] of Object.entries(emotions)) {
              if (score > maxEmotionScore) {
                maxEmotionScore = score;
                maxEmotion = emotion;
              }
            }

            // Update sentiment state
            setCurrentSentiment({
              sentiment: sentimentData.sentiment,
              score: sentimentData.score,
              maxEmotion,
              maxEmotionScore,
            });
          } else if (streamData.type === "transcript") {
            const transcriptData = streamData.data as TranscriptData;
            console.log(
              `Transcript from ${participantInfo.identity}:`,
              transcriptData
            );

            // YouTube-style live transcript: Update the current speaker's message in real-time
            setTranscripts((prev) => {
              // Check if the last transcript is from the same speaker
              const lastTranscript = prev[prev.length - 1];

              if (
                lastTranscript &&
                lastTranscript.role === transcriptData.role &&
                lastTranscript.name === transcriptData.name &&
                lastTranscript.timestamp === transcriptData.timestamp
              ) {
                // Same speaker, same turn - update the content in place
                const updated = [...prev];
                updated[updated.length - 1] = transcriptData;
                return updated;
              } else {
                // New speaker or new turn - add as new entry
                return [...prev, transcriptData];
              }
            });
          }
        }
      );
    } catch (error) {
      console.warn("Text stream handler already registered:", error);
    }

    return () => {
      try {
        room.unregisterTextStreamHandler("room");
      } catch (error) {
        console.warn("Error unregistering text stream handler:", error);
      }
    };
  }, []);

  return (
    <>
      <div
        ref={dragRef}
        className="draggable-controls-container fixed z-[9999] bg-white rounded-lg shadow-lg px-3 py-2 overflow-visible border border-gray-200/50 flex flex-row items-center gap-2 select-none"
        style={{
          top: `${Math.max(0, position.y)}px`,
          left: `${position.x}px`,
          transform: !hasBeenDragged ? "translateX(-50%)" : "none",
          cursor: isDragging ? "grabbing" : "move",
          position: "fixed",
          willChange: "transform",
        }}
        onMouseDown={handleMouseDown}
        onTouchStart={handleTouchStart}
        {...props}
      >
        {/* Drag Handle - Visual indicator */}
        <div className="drag-handle flex items-center justify-center px-1.5 py-1 hover:bg-gray-100 rounded transition-colors pointer-events-none">
          <div className="flex gap-0.5">
            <div className="w-1 h-1 rounded-full bg-gray-300" />
            <div className="w-1 h-1 rounded-full bg-gray-300" />
            <div className="w-1 h-1 rounded-full bg-gray-300" />
          </div>
        </div>

        {/* Divider */}
        <div className="w-px h-6 bg-gray-200" />

        {/* Controls */}
        {isWeb() && (
          <LayoutContextProvider value={layoutContext}>
            <div className="flex items-center">
              <ControlBar
                controls={{
                  chat: false,
                  camera: false,
                  screenShare: false,
                  settings: false,
                }}
                variation="minimal"
              />
            </div>
          </LayoutContextProvider>
        )}
      </div>
      <RoomAudioRenderer />
      <ConnectionStateToast />
    </>
  );
};
