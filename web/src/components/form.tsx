import classNames from 'classnames';
import {
  DetailedHTMLProps,
  FormHTMLAttributes,
  PropsWithChildren,
} from 'react';

interface FormProps
  extends DetailedHTMLProps<
    FormHTMLAttributes<HTMLFormElement>,
    HTMLFormElement
  > {
  hasError?: boolean;
}

export function Form({
  children,
  className,
  hasError = false,
  action,
  ...rest
}: PropsWithChildren<FormProps>) {
  return (
    <form
      className={classNames(
        'flex items-center gap-2 bg-zinc-900 p-2 rounded-xl border border-zinc-800 ring-orange-400 ring-offset-2 ring-offset-zinc-950 focus-within:ring-1',
        className,
        {
          'border-red-500': hasError,
        },
      )}
      action={action}
      {...rest}
    >
      {children}
    </form>
  );
}
