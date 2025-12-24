import type { TrackReferenceOrPlaceholder } from "@livekit/components-core";
import {
  isEqualTrackRef,
  isTrackReference,
  isWeb,
  log,
} from "@livekit/components-core";
import { RoomEvent, Track } from "livekit-client";
import * as React from "react";
import {
  CarouselLayout,
  ConnectionStateToast,
  FocusLayout,
  FocusLayoutContainer,
  GridLayout,
  LayoutContextProvider,
  ParticipantTile,
  RoomAudioRenderer,
  useCreateLayoutContext,
  usePinnedTracks,
  useTracks,
  ControlBar,
  useRoomContext,
} from "@livekit/components-react";

export interface VideoConferenceProps extends React.HTMLAttributes<HTMLDivElement> {
  meetingId: string;
}

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

export const VideoConference = ({
  meetingId,
  ...props
}: VideoConferenceProps) => {
  const lastAutoFocusedScreenShareTrack =
    React.useRef<TrackReferenceOrPlaceholder | null>(null);

  const room = useRoomContext();
  const [transcripts, setTranscripts] = React.useState<TranscriptData[]>([]);
  const [currentSentiment, setCurrentSentiment] = React.useState<{
    sentiment: string;
    score: number;
    maxEmotion: string;
    maxEmotionScore: number;
  } | null>(null);

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

  const tracks = useTracks(
    [
      { source: Track.Source.Camera, withPlaceholder: true },
      { source: Track.Source.ScreenShare, withPlaceholder: false },
    ],
    { updateOnlyOn: [RoomEvent.ActiveSpeakersChanged], onlySubscribed: false }
  );

  const layoutContext = useCreateLayoutContext();

  const screenShareTracks = tracks
    .filter(isTrackReference)
    .filter((track) => track.publication.source === Track.Source.ScreenShare);

  const focusTrack = usePinnedTracks(layoutContext)?.[0];
  const carouselTracks = tracks.filter(
    (track) => !isEqualTrackRef(track, focusTrack)
  );

  React.useEffect(() => {
    // If screen share tracks are published, and no pin is set explicitly, auto set the screen share.
    if (
      screenShareTracks.some((track) => track.publication.isSubscribed) &&
      lastAutoFocusedScreenShareTrack.current === null
    ) {
      log.debug("Auto set screen share focus:", {
        newScreenShareTrack: screenShareTracks[0],
      });
      layoutContext.pin.dispatch?.({
        msg: "set_pin",
        trackReference: screenShareTracks[0],
      });
      lastAutoFocusedScreenShareTrack.current = screenShareTracks[0];
    } else if (
      lastAutoFocusedScreenShareTrack.current &&
      !screenShareTracks.some(
        (track) =>
          track.publication.trackSid ===
          lastAutoFocusedScreenShareTrack.current?.publication?.trackSid
      )
    ) {
      log.debug("Auto clearing screen share focus.");
      layoutContext.pin.dispatch?.({ msg: "clear_pin" });
      lastAutoFocusedScreenShareTrack.current = null;
    }
    if (focusTrack && !isTrackReference(focusTrack)) {
      const updatedFocusTrack = tracks.find(
        (tr) =>
          tr.participant.identity === focusTrack.participant.identity &&
          tr.source === focusTrack.source
      );
      if (
        updatedFocusTrack !== focusTrack &&
        isTrackReference(updatedFocusTrack)
      ) {
        layoutContext.pin.dispatch?.({
          msg: "set_pin",
          trackReference: updatedFocusTrack,
        });
      }
    }
  }, [
    screenShareTracks
      .map(
        (ref) => `${ref.publication.trackSid}_${ref.publication.isSubscribed}`
      )
      .join(),
    focusTrack?.publication?.trackSid,
    tracks,
  ]);

  return (
    <div className="h-screen w-screen flex" {...props}>
      {/* Video Conference Section - Left Side */}
      <div className="flex-1 flex flex-col bg-zinc-900">
        {isWeb() && (
          <LayoutContextProvider value={layoutContext}>
            <div className="lk-video-conference-inner h-full w-full flex flex-col">
              {!focusTrack ? (
                <div className="lk-grid-layout-wrapper flex-1 min-h-0">
                  <GridLayout tracks={tracks}>
                    <ParticipantTile />
                  </GridLayout>
                </div>
              ) : (
                <div className="lk-focus-layout-wrapper flex-1 min-h-0">
                  <FocusLayoutContainer>
                    <CarouselLayout tracks={carouselTracks}>
                      <ParticipantTile />
                    </CarouselLayout>
                    {focusTrack && <FocusLayout trackRef={focusTrack} />}
                  </FocusLayoutContainer>
                </div>
              )}
              <ControlBar
                controls={{ chat: false, settings: false }}
                saveUserChoices={true}
              />
            </div>
          </LayoutContextProvider>
        )}
        <RoomAudioRenderer />
        <ConnectionStateToast />
      </div>
    </div>
  );
};
