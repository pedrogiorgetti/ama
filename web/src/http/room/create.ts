import { Room } from '../../interfaces/room';

interface CreateRoomRequest {
  body: {
    name: string;
  };
  functions: {
    stopLoading: () => void;
  };
}

interface CreateRoomResponse {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export async function createRoomRequest({
  body,
  functions,
}: CreateRoomRequest): Promise<Room> {
  const response = await fetch(`${import.meta.env.VITE_APP_API_URL}/rooms`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(body),
  });

  const data: CreateRoomResponse = await response.json();

  functions.stopLoading();

  return {
    id: data.id,
    name: data.name,
    createdAt: data.created_at,
    updatedAt: data.updated_at,
  };
}
