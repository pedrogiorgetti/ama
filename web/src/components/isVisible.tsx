import { PropsWithChildren } from 'react';

export function IsVisible({
  when,
  children,
}: PropsWithChildren<{ when: boolean }>) {
  if (!when) {
    return null;
  }

  return children;
}
