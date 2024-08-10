interface RemoveQuestionReactionRequest {
  params: {
    id: string;
    roomId: string;
  };
  functions: {
    onSuccess: () => void;
  };
}

export async function removeQuestionReactionRequest({
  params,
  functions,
}: RemoveQuestionReactionRequest): Promise<void> {
  const response = await fetch(
    `${import.meta.env.VITE_APP_API_URL}/rooms/${params.roomId}/questions/${params.id}`,
    {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    },
  );

  await response.json();

  functions.onSuccess();
}
