import classNames from 'classnames';
import {
  DetailedHTMLProps,
  FormHTMLAttributes,
  PropsWithChildren,
} from 'react';

export function FormComponent({
  children,
  className,
  action,
  ...rest
}: PropsWithChildren<
  DetailedHTMLProps<FormHTMLAttributes<HTMLFormElement>, HTMLFormElement>
>) {
  return (
    <form
      className={classNames(
        'flex items-center gap-2 bg-zinc-900 p-2 rounded-xl border border-zinc-800 ring-orange-400 ring-offset-2 ring-offset-zinc-950 focus-within:ring-1',
        className,
      )}
      action={action}
      {...rest}
    >
      {children}
    </form>
  );
}
