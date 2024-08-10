import { useState } from 'react';
import { Button } from '../button';
import classNames from 'classnames';
import { ArrowUp } from 'lucide-react';
import { useParams } from 'react-router-dom';
import { Question as QuestionInterface } from '../../interfaces/question';
import { toast } from 'sonner';
import { addQuestionReactionRequest } from '../../http/room/addQuestionReaction';
import { removeQuestionReactionRequest } from '../../http/room/removeQuestionReaction';

const reactionRequestByValue = {
  remove: removeQuestionReactionRequest,
  add: addQuestionReactionRequest,
};

interface QuestionProps {
  question: QuestionInterface;
}

export function Question({
  question: { reactionCount, text, isAnswered, id },
}: QuestionProps) {
  const { id: roomId } = useParams();

  if (!roomId) {
    throw new Error('Questions component must be used within room page');
  }

  const [isReacted, setIsReacted] = useState<boolean>(false);

  async function handleToggleReaction() {
    try {
      if (!roomId) {
        toast.error('Room ID is required');
        return;
      }

      const reactionRequestType = isReacted ? 'remove' : 'add';

      const questionReactionRequest =
        reactionRequestByValue[reactionRequestType];

      await questionReactionRequest({
        params: {
          roomId,
          id,
        },
        functions: {
          onSuccess: () => setIsReacted(true),
        },
      });
    } catch {
      const errorMessageComplement = isReacted
        ? 'removing the reaction'
        : 'reacting to the question';
      toast.error(`An error occurred while ${errorMessageComplement}`);
    }
  }

  return (
    <li
      className={classNames('ml-4 leading-relaxed text-zinc-100', {
        'opacity-40 pointer-events-none': isAnswered,
      })}
    >
      {text}
      <Button
        className={classNames('mt-3 ', {
          'text-orange-400 hover:text-orange-500': isReacted,
          'text-zinc-400 hover:opacity-60': !isReacted,
        })}
        onClick={handleToggleReaction}
      >
        <ArrowUp className="size-4" />
        Like question ({reactionCount})
      </Button>
    </li>
  );
}
