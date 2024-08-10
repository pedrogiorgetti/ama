import { Question } from '../../interfaces/question';

interface CreateQuestionRequest {
  params: {
    id: string;
  };
  body: {
    text: string;
  };
  functions: {
    stopLoading: () => void;
  };
}

interface CreateQuestionResponse {
  id: string;
  text: string;
  room_id: string;
  reaction_count: number;
  answered: boolean;
  created_at: string;
  updated_at: string;
}

export async function createQuestionRequest({
  params,
  body,
  functions,
}: CreateQuestionRequest): Promise<Question> {
  const response = await fetch(
    `${import.meta.env.VITE_APP_API_URL}/rooms/${params.id}/questions`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    },
  );

  const data: CreateQuestionResponse = await response.json();

  functions.stopLoading();

  return {
    id: data.id,
    text: data.text,
    roomId: data.room_id,
    reactionCount: data.reaction_count,
    isAnswered: data.answered,
    createdAt: data.created_at,
    updatedAt: data.updated_at,
  };
}
