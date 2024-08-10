import { Question } from '../../interfaces/question';

interface GetRoomQuestionsRequest {
  params: {
    id: string;
  };
}

interface GetRoomQuestionsResponse {
  list: {
    ID: string;
    Text: string;
    RoomID: string;
    ReactionCount: number;
    Answered: boolean;
  }[];
  total: number;
}

export interface GetRoomQuestionsResponseData {
  list: Question[];
  total: number;
}

export async function getRoomQuestionsRequest({
  params,
}: GetRoomQuestionsRequest): Promise<GetRoomQuestionsResponseData> {
  const response = await fetch(
    `${import.meta.env.VITE_APP_API_URL}/rooms/${params.id}/questions`,
    {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    },
  );

  const data: GetRoomQuestionsResponse = await response.json();

  const questionsFormatted: Question[] = data.list.map(item => ({
    id: item.ID,
    text: item.Text,
    roomId: item.RoomID,
    reactionCount: item.ReactionCount,
    isAnswered: item.Answered,
  }));

  return {
    list: questionsFormatted,
    total: data.total,
  };
}
