import { useParams } from 'react-router-dom';
import amaLogo from '../assets/logo.svg';
import { ButtonComponent } from '../components/button';
import { ArrowRight, Share2 } from 'lucide-react';
import { InputComponent } from '../components/input';
import { FormComponent } from '../components/form';

import { toast } from 'sonner';
import { QuestionComponent } from '../components/question';

export function Room() {
  const params = useParams();

  function handleShareRoom() {
    const url = window.location.href.toString();

    if (navigator.share !== undefined && navigator.canShare()) {
      navigator.share({ url });
      return;
    }

    navigator.clipboard.writeText(url);

    toast.info('Room URL copied to clipboard');
  }

  function handleCreateQuestion(data: FormData) {
    const text = data.get('text') as string;

    console.log('Criar sala com o tema:', text);
  }

  return (
    <div className="mx-auto max-w-[640px] flex flex-col gap-6 py-10 px-4">
      <div className="flex items-center gap-3 px-3">
        <img className="h-5" src={amaLogo} alt="AMA - Logo" />

        <h4 className="text-sm text-zinc-500 truncate">
          Room ID: <b className="text-zinc-300">{params.id}</b>
        </h4>

        <ButtonComponent
          className="ml-auto bg-zinc-800 text-zinc-300 hover:bg-zinc-900"
          type="submit"
          onClick={handleShareRoom}
        >
          Share
          <Share2 className="size-4" />
        </ButtonComponent>
      </div>

      <hr className="h-px w-full bg-zinc-900" />

      <FormComponent action={handleCreateQuestion}>
        <InputComponent name="text" placeholder="What is your question?" />
        <ButtonComponent
          className="bg-orange-400 text-orange-950 hover:bg-orange-500"
          type="submit"
        >
          Create question
          <ArrowRight className="size-4" />
        </ButtonComponent>
      </FormComponent>

      <ol className="list-decimal list-outside px-3 space-y-8">
        <QuestionComponent reactionCount={24} text="Question one" />
      </ol>
    </div>
  );
}
