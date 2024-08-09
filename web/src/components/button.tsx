import {
  ButtonHTMLAttributes,
  DetailedHTMLProps,
  PropsWithChildren,
} from 'react';
import classNames from 'classnames';

export function ButtonComponent({
  className,
  type = 'button',
  children,
  ...rest
}: PropsWithChildren<
  DetailedHTMLProps<ButtonHTMLAttributes<HTMLButtonElement>, HTMLButtonElement>
>) {
  return (
    <button
      className={classNames(
        'px-3 py-1.5 gap-1.5 flex items-center rounded-lg font-medium text-sm transition-colors',
        className,
      )}
      type={type}
      {...rest}
    >
      {children}
    </button>
  );
}
