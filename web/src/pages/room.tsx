import { Suspense, useState } from 'react';
import { useParams } from 'react-router-dom';
import { ArrowRight, Share2 } from 'lucide-react';
import { toast } from 'sonner';

import amaLogo from '../assets/logo.svg';
import { Button } from '../components/button';
import { Input } from '../components/input';
import { Form } from '../components/form';
import { Questions } from '../components/question/list';
import { createQuestionRequest } from '../http/room/createQuestion';

export function Room() {
  const { id: roomId } = useParams();

  const [isErrored, setIsErrored] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  function handleShareRoom() {
    const url = window.location.href.toString();

    if (navigator.share !== undefined && navigator.canShare()) {
      navigator.share({ url });
      return;
    }

    navigator.clipboard.writeText(url);

    toast.info('Room URL copied to clipboard');
  }

  async function handleCreateQuestion(data: FormData) {
    try {
      if (!roomId) {
        toast.error('Room ID is required');
        return;
      }

      const text = data.get('text') as string;

      if (!text) {
        setIsErrored(true);
      }

      setIsErrored(false);
      setIsLoading(true);

      await createQuestionRequest({
        params: {
          id: roomId,
        },
        body: {
          text,
        },
        functions: {
          stopLoading: () => setIsLoading(false),
        },
      });
    } catch {
      setIsLoading(false);
      toast.error('An error occurred while creating the question');
    }
  }

  return (
    <div className="mx-auto max-w-[640px] flex flex-col gap-6 py-10 px-4">
      <div className="flex items-center gap-3 px-3">
        <img className="h-5" src={amaLogo} alt="AMA - Logo" />

        <h4 className="text-sm text-zinc-500 truncate">
          Room ID: <b className="text-zinc-300">{roomId}</b>
        </h4>

        <Button
          className="ml-auto bg-zinc-800 text-zinc-300 hover:bg-zinc-900"
          type="submit"
          onClick={handleShareRoom}
        >
          Share
          <Share2 className="size-4" />
        </Button>
      </div>

      <hr className="h-px w-full bg-zinc-900" />

      <Form hasError={isErrored} action={handleCreateQuestion}>
        <Input name="text" placeholder="What is your question?" />
        <Button
          className="bg-orange-400 text-orange-950 hover:bg-orange-500"
          type="submit"
          isLoading={isLoading}
        >
          Create question
          <ArrowRight className="size-4" />
        </Button>
      </Form>

      <Suspense fallback={<p>Loading questions...</p>}>
        <Questions />
      </Suspense>
    </div>
  );
}
