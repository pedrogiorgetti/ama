import { ArrowRight } from 'lucide-react';
import { toast } from 'sonner';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import amaLogo from '../assets/logo.svg';
import { pages } from '../constants/pages';
import { Button } from '../components/button';
import { Input } from '../components/input';
import { Form } from '../components/form';
import { createRoomRequest } from '../http/room/create';

export function CreateRoom() {
  const navigate = useNavigate();

  const [isErrored, setIsErrored] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  async function handleCreateRoom(data: FormData) {
    try {
      const name = data.get('name') as string;

      if (!name) {
        setIsErrored(true);
      }

      setIsErrored(false);
      setIsLoading(true);

      const room = await createRoomRequest({
        body: {
          name,
        },
        functions: {
          stopLoading: () => setIsLoading(false),
        },
      });

      navigate(pages.room.details(room.id));
    } catch {
      setIsLoading(false);
      toast.error(`
        An error occurred while creating the room
      `);
    }
  }

  return (
    <main className="h-screen flex items-center justify-center px-4">
      <div className="max-w-[450px] flex flex-col gap-6">
        <img className="h-10" src={amaLogo} alt="AMA - Logo" />

        <p className="leading-relaxed text-zinc-300 text-center">
          Create a public AMA (Ask me anything) room and prioritise the most
          important questions for the community.
        </p>

        <Form hasError={isErrored} action={handleCreateRoom}>
          <Input name="name" placeholder="Room's name" />
          {isErrored && (
            <small className="text-red-500">Room's name is required!</small>
          )}

          <Button
            isLoading={isLoading}
            className="bg-orange-400 text-orange-950 hover:bg-orange-500"
            type="submit"
          >
            Create room
            <ArrowRight className="size-4" />
          </Button>
        </Form>
      </div>
    </main>
  );
}
