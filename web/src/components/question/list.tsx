import { useParams } from 'react-router-dom';

import { useSuspenseQuery } from '@tanstack/react-query';

import { getRoomQuestionsRequest } from '../../http/room/getQuestions';

import { Question } from './item';
import { useQuestionsWebSocket } from '../../hooks/useQuestionsWebSocket';

export function Questions() {
  const { id: roomId } = useParams();

  if (!roomId) {
    throw new Error('Questions component must be used within room page');
  }

  const {
    data: { list, total },
  } = useSuspenseQuery({
    queryKey: ['questions', roomId],
    queryFn: async () =>
      await getRoomQuestionsRequest({
        params: { id: roomId },
      }),
  });

  useQuestionsWebSocket({ roomId });

  const sortedQuestions = list.sort(
    (a, b) => b.reactionCount - a.reactionCount,
  );

  return (
    <div>
      <span>
        Questions
        <b>({total})</b>
      </span>
      <ol className="list-decimal list-outside px-3 space-y-8">
        {sortedQuestions.map(question => (
          <Question key={question.id} question={question} />
        ))}
      </ol>
    </div>
  );
}
