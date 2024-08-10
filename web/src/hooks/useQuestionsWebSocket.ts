import { useEffect } from 'react';
import { GetRoomQuestionsResponseData } from '../http/room/getQuestions';
import { useQueryClient } from '@tanstack/react-query';

enum EWebsocketNotificationCategory {
  Created = 'question_created',
  ReactionIncrease = 'question_reaction_increase',
  ReactionDecrease = 'question_reaction_decrease',
  Answered = 'question_answered',
}

interface WebsocketNotificationData {
  category: EWebsocketNotificationCategory;
  value: {
    id: string;
    text: string;
    count: number;
  };
}

interface UseQuestionsWebSocket {
  roomId: string;
}

export function useQuestionsWebSocket({ roomId }: UseQuestionsWebSocket) {
  const queryClient = useQueryClient();

  useEffect(() => {
    const websocket = new WebSocket(
      `${import.meta.env.VITE_APP_WS_URL}/subscribe/${roomId}`,
    );

    websocket.onopen = () => {
      console.log('WebSocket connected');
    };

    websocket.onmessage = event => {
      const parsedData: WebsocketNotificationData = JSON.parse(event.data);

      switch (parsedData.category) {
        case EWebsocketNotificationCategory.Created:
          queryClient.setQueryData<GetRoomQuestionsResponseData>(
            ['questions', roomId],
            currentData => {
              return {
                list: [
                  ...(currentData?.list ?? []),
                  {
                    id: parsedData.value.id,
                    text: parsedData.value.text,
                    reactionCount: 0,
                    isAnswered: false,
                    roomId,
                  },
                ],
                total: (currentData?.total ?? 0) + 1,
              };
            },
          );
          break;

        case EWebsocketNotificationCategory.ReactionIncrease:
        case EWebsocketNotificationCategory.ReactionDecrease:
          queryClient.setQueryData<GetRoomQuestionsResponseData>(
            ['questions', roomId],
            currentData => {
              if (!currentData) {
                return undefined;
              }

              return {
                list: currentData.list.map(question => {
                  if (question.id === parsedData.value.id) {
                    return {
                      ...question,
                      reactionCount: parsedData.value.count,
                    };
                  }

                  return question;
                }),
                total: currentData.total,
              };
            },
          );
          break;

        case EWebsocketNotificationCategory.Answered:
          queryClient.setQueryData<GetRoomQuestionsResponseData>(
            ['questions', roomId],
            currentData => {
              if (!currentData) {
                return undefined;
              }

              return {
                list: currentData.list.map(question => {
                  if (question.id === parsedData.value.id) {
                    return {
                      ...question,
                      isAnswered: true,
                    };
                  }

                  return question;
                }),
                total: currentData.total,
              };
            },
          );
          break;
      }
    };

    return () => {
      websocket.close();
    };
  }, [queryClient, roomId]);
}
