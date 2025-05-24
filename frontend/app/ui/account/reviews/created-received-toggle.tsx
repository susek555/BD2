'use client';

import clsx from 'clsx';
import { useRouter, useSearchParams } from 'next/navigation';

const orderOptions = [
  { label: 'By you', value: 'for' },
  { label: 'From you', value: 'by' },
];

export function ReviewToggle() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const currentVariant = searchParams.get('variant') || 'for';

  const handleToggle = (variant: string) => {
    const params = new URLSearchParams(searchParams);
    params.set('variant', variant);
    router.push(`?${params.toString()}`);
  };

  return (
    <div className='inline-flex rounded-md shadow-sm'>
      {orderOptions.map((option) => (
        <button
          key={option.value}
          onClick={() => handleToggle(option.value)}
          className={clsx(
            'border px-4 py-2 text-sm font-medium',
            option.value === currentVariant
              ? 'z-10 border-blue-500 bg-blue-50 text-blue-700'
              : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50',
            orderOptions[0].value === option.value && 'rounded-l-md',
            orderOptions[orderOptions.length - 1].value === option.value &&
              'rounded-r-md',
          )}
        >
          {option.label}
        </button>
      ))}
    </div>
  );
}
