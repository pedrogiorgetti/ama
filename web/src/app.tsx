import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { CreateRoom } from './pages/create-room';
import { pages } from './constants/pages';
import { Room } from './pages/room';
import { Toaster } from 'sonner';
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from './lib/react-query';

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
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      <Toaster invert richColors />
    </QueryClientProvider>
  );
}
