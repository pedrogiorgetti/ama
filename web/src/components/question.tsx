import { useState } from 'react';
import { ButtonComponent } from './button';
import classNames from 'classnames';
import { ArrowUp } from 'lucide-react';

interface QuestionComponentProps {
  text: string;
  reactionCount: number;
  isAnswered?: boolean;
}

export function QuestionComponent({
  reactionCount,
  text,
  isAnswered,
}: QuestionComponentProps) {
  const [isReacted, setIsReacted] = useState<boolean>(false);

  return (
    <li
      className={classNames('ml-4 leading-relaxed text-zinc-100', {
        'opacity-40 pointer-events-none': isAnswered,
      })}
    >
      {text}
      <ButtonComponent
        className={classNames('mt-3 ', {
          'text-orange-400 hover:text-orange-500': isReacted,
          'text-zinc-400 hover:opacity-60': !isReacted,
        })}
        onClick={() => setIsReacted(!isReacted)}
      >
        <ArrowUp className="size-4" />
        Like question ({reactionCount})
      </ButtonComponent>
    </li>
  );
}
