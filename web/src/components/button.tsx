import {
  ButtonHTMLAttributes,
  DetailedHTMLProps,
  PropsWithChildren,
} from 'react';
import classNames from 'classnames';
import { IsVisible } from './isVisible';
import { ClipLoader } from 'react-spinners';

interface ButtonProps
  extends DetailedHTMLProps<
    ButtonHTMLAttributes<HTMLButtonElement>,
    HTMLButtonElement
  > {
  isLoading?: boolean;
}

export function Button({
  className,
  type = 'button',
  children,
  disabled = false,
  isLoading = false,
  ...rest
}: PropsWithChildren<ButtonProps>) {
  return (
    <button
      className={classNames(
        'px-3 py-1.5 gap-1.5 flex items-center rounded-lg font-medium text-sm transition-colors',
        className,
      )}
      type={type}
      disabled={disabled || isLoading}
      {...rest}
    >
      <IsVisible when={!isLoading}>{children}</IsVisible>
      <IsVisible when={isLoading}>
        <ClipLoader size={12} />
      </IsVisible>
    </button>
  );
}
