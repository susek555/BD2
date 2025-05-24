import clsx from 'clsx';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
}

export function BasePriceButton({
  children,
  className,
  onClick,
  ...rest
}: ButtonProps & { onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void }) {
  return (
    <button
      {...rest}
      onClick={onClick}
      className={clsx(
        'flex md:w-80 h-12 items-center rounded bg-blue-500 px-2 text-sm font-medium text-white transition-colors hover:bg-blue-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-500 active:bg-blue-600 aria-disabled:cursor-not-allowed aria-disabled:opacity-50',
        className,
      )}
    >
      {children}
    </button>
  );
}