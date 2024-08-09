import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { CreateRoom } from './pages/create-room';
import { pages } from './constants/pages';
import { Room } from './pages/room';
import { Toaster } from 'sonner';

const router = createBrowserRouter([
  {
    path: pages.room.create,
    element: <CreateRoom />,
  },
  {
    path: '/room/:id',
    element: <Room />,
  },
]);

export function App() {
  return (
    <>
      <RouterProvider router={router} />
      <Toaster invert richColors />
    </>
  );
}
