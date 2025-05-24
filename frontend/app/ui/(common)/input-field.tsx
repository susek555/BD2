import { InputHTMLAttributes } from "react";

interface InputFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  id: string;
  name: string;
  type?: string;
  label: string;
  placeholder: string;
  icon: React.ComponentType<React.SVGProps<SVGSVGElement>>;
  defaultValue: string;
  errors: string[];
}

export default function InputField({
  id,
  name,
  type = 'text',
  label,
  placeholder,
  icon: Icon,
  defaultValue,
  errors,
  ...rest
}: InputFieldProps) {
  return (
    <div className='mb-4'>
      <label
        className='mb-3 block text-xs font-medium text-gray-900'
        htmlFor={id}
      >
        {label}
      </label>
      <div className='relative'>
        <input
          className='peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500'
          id={id}
          type={type}
          name={name}
          placeholder={placeholder}
          defaultValue={defaultValue}
          {...rest}
        />
        {Icon && (
          <Icon className='pointer-events-none absolute top-1/2 left-3 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900' />
        )}
      </div>
      <div aria-live='polite' aria-atomic='true'>
        {errors &&
          errors.map((error) => (
            <p className='mt-2 text-sm text-red-500' key={error}>
              {error}
            </p>
          ))}
      </div>
    </div>
  );
}
