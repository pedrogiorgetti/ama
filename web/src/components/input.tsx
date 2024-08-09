import classNames from 'classnames';
import { DetailedHTMLProps, InputHTMLAttributes } from 'react';

export function InputComponent({
  name,
  type = 'text',
  placeholder,
  className,
  ...rest
}: DetailedHTMLProps<InputHTMLAttributes<HTMLInputElement>, HTMLInputElement>) {
  return (
    <input
      type={type}
      name={name}
      placeholder={placeholder}
      autoComplete="off"
      className={classNames(
        'bg-transparent flex-1 mx-2 outline-none text-zinc-100 placeholder:text-zinc-500',
        className,
      )}
      {...rest}
    />
  );
}
