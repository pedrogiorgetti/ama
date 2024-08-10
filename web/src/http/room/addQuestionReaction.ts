interface AddQuestionReactionRequest {
  params: {
    id: string;
    roomId: string;
  };
  functions: {
    onSuccess: () => void;
  };
}

export async function addQuestionReactionRequest({
  params,
  functions,
}: AddQuestionReactionRequest): Promise<void> {
  const response = await fetch(
    `${import.meta.env.VITE_APP_API_URL}/rooms/${params.roomId}/questions/${params.id}`,
    {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
    },
  );

  await response.json();

  functions.onSuccess();
}
