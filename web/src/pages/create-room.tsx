import { ArrowRight } from 'lucide-react';
import amaLogo from '../assets/logo.svg';
import { useNavigate } from 'react-router-dom';
import { pages } from '../constants/pages';
import { ButtonComponent } from '../components/button';
import { InputComponent } from '../components/input';
import { FormComponent } from '../components/form';

export function CreateRoom() {
  const navigate = useNavigate();

  function handleCreateRoom(data: FormData) {
    const theme = data.get('theme') as string;

    console.log('Criar sala com o tema:', theme);
    navigate(pages.room.details(123123));
  }

  return (
    <main className="h-screen flex items-center justify-center px-4">
      <div className="max-w-[450px] flex flex-col gap-6">
        <img className="h-10" src={amaLogo} alt="AMA - Logo" />

        <p className="leading-relaxed text-zinc-300 text-center">
          Create a public AMA (Ask me anything) room and prioritise the most
          important questions for the community.
        </p>

        <FormComponent action={handleCreateRoom}>
          <InputComponent name="theme" placeholder="Room's name" />
          <ButtonComponent
            className="bg-orange-400 text-orange-950 hover:bg-orange-500"
            type="submit"
          >
            Create room
            <ArrowRight className="size-4" />
          </ButtonComponent>
        </FormComponent>
      </div>
    </main>
  );
}
